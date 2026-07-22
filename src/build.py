import re, sys
import importlib.util
import latex2mathml.converter

from .hooks import hook, hook_call, reset_hooks
from .hash import handle_hash_sync, clear_all_hashes
from .api import expose, api

from jinja2 import Environment, FileSystemLoader
from markdown import markdown
from markdown.extensions.attr_list import AttrListTreeprocessor
from os import chdir, getcwd, makedirs, path, walk, readlink

from .config import find_project_from_path, load_config
from .errors import fatal, html_fatal
from .fileops import is_dotfile, read_file, write_file
from .log import GRAY, BLUE, die, info, warn, progress
from . import log

api.log = log

AttrListTreeprocessor.BASE_RE   = r'\[\:?[ ]*([^\]\n ][^\n]*)[ ]*\]'
AttrListTreeprocessor.BLOCK_RE  = re.compile(r'\n[ ]*{}[ ]*$'.format(AttrListTreeprocessor.BASE_RE))
AttrListTreeprocessor.HEADER_RE = re.compile(r'[ ]+{}[ ]*$'.format(AttrListTreeprocessor.BASE_RE))
AttrListTreeprocessor.INLINE_RE = re.compile(r'^{}'.format(AttrListTreeprocessor.BASE_RE))

replace_filters = [
    ("<p>{%" ,  "{%" ),
    ("%}</p>",  "%}" ),
]

def load_plugins(config):
    reset_hooks()
    module_dir = path.abspath(config.tree.plugins)
    sys.path.insert(0, module_dir)
    try:
        main_path = path.join(module_dir, "main.py")
        if not path.isfile(main_path):
            die("can not find `main.py` on plugins directory")
        for name in list(sys.modules):
            mod_file = getattr(sys.modules[name], "__file__", "") or ""
            if mod_file.startswith(module_dir):
                del sys.modules[name]
        spec = importlib.util.spec_from_file_location("plugins", main_path)
        if spec is None or spec.loader is None:
            raise RuntimeError("Failed to create plugin module spec.")
        module = importlib.util.module_from_spec(spec)

        module.api = api
        module.hook = hook
        spec.loader.exec_module(module)
        hook_call("on_plugins_before_export", config)

        exports = {
            name: value
            for name, value in vars(module).items()
            if not name.startswith("_")
        }
        hook_call("on_plugins_loaded", exports)
        return exports

    except Exception as e:
        raise RuntimeError(f"Failed loading plugin {config.tree.plugins}: {e}") from e
    finally:
        sys.path.pop(0)


@expose("html_filter")
def html_filter(html_content: str):
    html = []
    for line in html_content.splitlines():
        for fltr in replace_filters:
            line = line.replace(*fltr)
        html.append(line)
    result = "\n".join(html)
    result = re.sub(r'(}}\s*)\[[^\]]*\]', r'\1', result)
    return result

def render_math(md_content: str) -> str:
    def convert(m):
        return latex2mathml.converter.convert(m.group(1))

    md_content = re.sub(r'\$\$(.+?)\$\$', convert, md_content, flags=re.DOTALL)
    md_content = re.sub(r'\$(.+?)\$', convert, md_content)

    return md_content

@expose("jinja_handler")
def jinja_handler(config, html_content, plugins=None):
    try:
        tree = config.tree if config else None
        if tree:
            templates = tree.templates
            templates = hook_call("on_jinja_template_dir", templates) or templates
            env = (Environment(loader=FileSystemLoader(templates)))

        else:
            env = (Environment())

        template = env.from_string(html_content)
        result = template.render(**plugins) if plugins is not None else template.render()
        result = hook_call("on_jinja_template_renderer", result) or result
        return result

    except Exception as e:
        hook_call("on_jinja_error", e)
        return html_fatal(e, f"Template error")

@expose("save_html")
def save_html(html_content: str, html_dest: str):
    makedirs(path.dirname(html_dest), exist_ok=True)
    write_file(html_dest, html_content)
    hook_call("on_page_written", html_dest, html_content)

def process_highlighting(config):
    highlight = "noclasses"
    if config:
        highlight = config.extras.highlight
    highlight = hook_call("on_highlight_config", highlight) or highlight
    return {
        "use_pygments": True,
        "linenums": True,
        "linenums_style": "table",
        "noclasses": highlight != "noclasses",
        **(
            {}
            if highlight == "noclasses"
            else {"pygments_style": highlight}
        ),
    }

@expose("md_to_html")
def md_to_html(config, md_content: str):
    math_rendered = render_math(md_content)
    math_rendered = hook_call("on_math_renderer", math_rendered) or math_rendered
    raw_html_content = markdown(
        math_rendered,
        extensions=[
            "extra",
            "md_in_html",
            "pymdownx.betterem",
            "pymdownx.critic",
            "pymdownx.details",
            "pymdownx.highlight",
            "pymdownx.inlinehilite",
            "pymdownx.keys",
            "pymdownx.mark",
            "pymdownx.superfences",
            "pymdownx.tabbed",
            "pymdownx.tilde",
        ], extension_configs={
            "pymdownx.highlight": process_highlighting(config)
        }
    )
    raw_html_content = hook_call("on_md_to_html", raw_html_content) or raw_html_content
    filtered_html = html_filter(raw_html_content)
    filtered_html = hook_call("on_html_filter", filtered_html) or filtered_html
    return filtered_html


@expose("compile_page")
def compile_page(md_content: str, html_dest: str | None = None, config=None, plugins=None):
    html_raw = md_to_html(config, md_content)
    html_content = jinja_handler(config, html_raw, plugins)
    html_content = hook_call("on_page_rendered", html_content) or html_content
    if html_dest is None:
        return html_content
    else:
        save_html(html_content, html_dest)
        return True

def read_md_content(md_file):
    md_content = read_file(md_file)
    md_content = hook_call("on_page_read", md_file, md_content) or md_content
    return md_content

@expose("compile_file")
def compile_file(md_file, html_dest, config=None, plugins=None, force: bool = False):
    hook_call("on_compile_start", md_file)
    if config and not force:
        hash_result = handle_hash_sync(config, md_file)
        if hash_result is None and path.exists(html_dest):
            hook_call("on_page_skip", md_file)
            info(f"Skipping {GRAY(md_file)}...")
            return None

    md_content = read_md_content(md_file)
    result = compile_page(md_content, html_dest, config, plugins)
    info(f"Building {GRAY(md_file)}...")
    hook_call("on_page_built", md_file, html_dest)
    return result

def walk_and_build(config, plugins, dest, force_all=False):
    md_path = config.tree.markdown
    hook_call("on_walk_start", config)
    files = []
    for parent, _, filenames in walk(md_path):
        for filename in filenames:
            files.append(path.join(parent, filename))
    for md_file in progress(files, "Building"):
        hook_call("on_walk_file", md_file)
        md_relpath = path.relpath(md_file, md_path)
        if is_dotfile(md_relpath): continue
        html_dest = path.join(dest, md_relpath)
        if html_dest.endswith(".md"):
            html_dest = html_dest.removesuffix(".md") + ".html"
        source_file = readlink(md_file) if path.islink(md_file) else md_file
        compile_file(source_file, html_dest, config, plugins, force=force_all)
    hook_call("on_walk_end", config)

def build_all(config, plugins, dest):
    clear_all_hashes(config)
    walk_and_build(config, plugins, dest, force_all=True)

def validate_build(config, plugins):
    errors = []
    md_path = config.tree.markdown
    for parent, _, filenames in walk(md_path):
        for filename in filenames:
            md_file = path.join(parent, filename)
            if is_dotfile(path.relpath(md_file, md_path)):
                continue
            source = readlink(md_file) if path.islink(md_file) else md_file
            try:
                compile_page(read_md_content(source), config=config, plugins=plugins)
            except Exception as e:
                errors.append(f"{path.relpath(md_file, md_path)}: {e}")
    if errors:
        for e in errors:
            warn(e)
        raise Exception(f"Validation failed with {len(errors)} error(s)")
    info("Validation passed.")

def build_if_file(project_path, file):
    if project_path:
        raise ValueError("Please specify either a project path or a file path, not both.")
    elif len(file) != 2:
        raise ValueError("A file path must include exactly a source and a destination.")
    else:
        compile_file(file[0], file[1])

def build(mode, project_path, file, validate=False, force=False):
    try:
        if file:
            return build_if_file(project_path, file)

        project_path = project_path or getcwd()
        find_project_from_path(project_path)
        chdir(project_path)
        config = load_config()
        plugins = load_plugins(config)
        dest = config.tree.draft_dest if mode == "draft" else config.tree.release_dest

        if validate:
            validate_build(config, plugins)
        if mode == "release":
            import shutil
            if path.exists(dest):
                shutil.rmtree(dest)

        hook_call("on_build_start", config)
        walk_and_build(config, plugins, dest, force_all=force)
        hook_call("on_build_end", config)
    except Exception as e:
        fatal(e, f"Build failed: {e}")

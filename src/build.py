import re, sys
import importlib.util
import latex2mathml.converter

from .hooks import hook, hook_call, hooks
from .hash import handle_hash_sync

from jinja2 import Environment, FileSystemLoader
from markdown import markdown
from markdown.extensions.attr_list import AttrListTreeprocessor
from os import chdir, getcwd, makedirs, path, walk, readlink

from .config import find_project_from_path, load_config
from .errors import fatal, html_fatal
from .fileops import is_dotfile, read_file, write_file
from .log import GRAY, die, info, warn
from . import log

# patch [ ] instead of { }
AttrListTreeprocessor.BASE_RE   = r'\[\:?[ ]*([^\]\n ][^\n]*)[ ]*\]'
AttrListTreeprocessor.BLOCK_RE  = re.compile(r'\n[ ]*{}[ ]*$'.format(AttrListTreeprocessor.BASE_RE))
AttrListTreeprocessor.HEADER_RE = re.compile(r'[ ]+{}[ ]*$'.format(AttrListTreeprocessor.BASE_RE))
AttrListTreeprocessor.INLINE_RE = re.compile(r'^{}'.format(AttrListTreeprocessor.BASE_RE))

replace_filters = [
    ("<p>{%" ,  "{%" ),
    ("%}</p>",  "%}" ),
]

def load_plugins(config):
    for k in hooks:
        hooks[k] = None
    module_dir = path.abspath(config.tree.plugins)
    sys.path.insert(0, module_dir)
    try:
        if not path.exists(path.join(module_dir, "main.py") ):
            die("can not find `main.py` on plugins directory")
        for name in list(sys.modules):
            mod_file = getattr(sys.modules[name], "__file__", "") or ""
            if mod_file.startswith(module_dir):
                del sys.modules[name]
        spec = importlib.util.spec_from_file_location("plugins", path.join(module_dir, "main.py") )
        if spec is None or spec.loader is None:
            raise RuntimeError("Failed to create plugin module spec.")
        module = importlib.util.module_from_spec(spec)
        module.CONFIG = config
        module.compile_page = compile_page
        module.hook = hook
        module.log = log
        spec.loader.exec_module(module)

        exports =  {
            name: value
            for name, value in vars(module).items()
            if not name.startswith("_")
        }
        hook_call("on_plugins_loaded",exports)
        return exports

    except Exception as e:
        raise RuntimeError(f"Failed loading plugin {config.tree.plugins}: {e}") from e
    finally:
        sys.path.pop(0)


def html_filter(html_content:str):
    """ filter jinja placeholders from html """
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

def jinja_handler(config, html_content, plugins=None):
    try:
        tree = config.tree if config else None
        if tree:
            templates = tree.templates 
            templates = hook_call("on_jinja_template_dir", templates) or templates
            env = ( Environment(loader=FileSystemLoader(templates)) )

        else:
            env = ( Environment() )

        template = env.from_string(html_content)
        result =  template.render(**plugins) if plugins is not None else template.render()
        result = hook_call("on_jinja_template_renderer", result) or result
        return result

    except Exception as e:
        hook_call("on_jinja_error", e)
        return html_fatal(e, f"Template error")

def escape_code_blocks(html_content: str) -> str:
    """ Neutralize '{' inside rendered <pre>/<code> blocks so Jinja never
    reads {% %} or {{ }} out of code samples, even if a sample's text
    literally contains Jinja syntax (e.g. docs showing `{% raw %}`).
    &#123; renders back to '{' in the browser, so output is unaffected. """
    def neutralize(m):
        block = m.group(0)
        block = block.replace("${", "__JINJA_OPEN__")
        block = block.replace("{", "&#123;")
        block = block.replace("__JINJA_OPEN__", "{")
        return block

    html_content = re.sub(r'<pre\b[\s\S]*?</pre>', neutralize, html_content)
    html_content = re.sub(r'<code\b[\s\S]*?</code>', neutralize, html_content)
    return html_content

def save_html(html_content:str, html_dest:str):
    makedirs(path.dirname(html_dest), exist_ok=True )
    write_file( html_dest, html_content)
    hook_call("on_page_written", html_dest, html_content)

def process_highlighting(config):
    highlight = "noclasses"
    if config:
        highlight = config.extras.highlight
    highlight = hook_call("on_highlight_config", highlight) or highlight
    return {
        "use_pygments": True,
        "noclasses": highlight != "noclasses",
        **(
            {}
            if highlight == "noclasses"
            else {"pygments_style": highlight}
        ),
    }

def md_to_html(config, md_content:str):
    """Convert a Markdown file to HTML, applying filters and Jinja2 processing, and save to dest."""
    math_rendered = render_math(md_content) 
    math_rendered = hook_call("on_math_renderer", math_rendered) or math_rendered
    raw_html_content = markdown(
        math_rendered,
        extensions = [
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
        ],extension_configs={
            "pymdownx.highlight": process_highlighting(config)
        }
    )
    raw_html_content = hook_call("on_md_to_html", raw_html_content) or raw_html_content
    escaped_html = escape_code_blocks(raw_html_content)
    escaped_html = hook_call("on_escape_code", escaped_html) or escaped_html
    filtered_html = html_filter(escaped_html)
    filtered_html = hook_call("on_html_filter", filtered_html) or filtered_html
    return filtered_html


def compile_page(md_content:str, html_dest:str | None=None, config=None, plugins=None):
    html_raw = md_to_html(config, md_content)
    html_content = jinja_handler(config, html_raw, plugins)
    html_content = hook_call("on_page_rendered", html_content) or html_content
    if html_dest is None:
        return html_content
    else:
        save_html(html_content, html_dest)
        return True

def compile_file(md_file, html_dest, config=None, plugins=None, force:bool = False):
    md_content = read_file(md_file)
    md_content = hook_call("on_page_read", md_file, md_content) or md_content
    if config and not force:
        hash = handle_hash_sync(config, md_file)
        if  hash is None and path.exists(html_dest):
            info(f"Skipping {GRAY(md_file)}...")
            return None 

    info(f"Building {GRAY(md_file)}...")
    return compile_page(md_content, html_dest, config, plugins)

def walk_and_build(config, plugins):
    md_path = config.tree.markdown
    hook_call("on_walk_start", config)
    for parent, _, files  in walk(md_path):
        for filename in files:
            md_file = path.join(parent, filename)
            hook_call("on_walk_file", md_file)
            md_relpath = path.relpath(md_file,md_path)
            if is_dotfile(md_relpath): continue
            html_dest = path.join(config.tree.dest, md_relpath)
            if html_dest.endswith(".md"):
                html_dest  = html_dest.removesuffix(".md") + ".html"
            source_file = readlink(md_file) if path.islink(md_file) else md_file
            compile_file(source_file, html_dest, config, plugins)
    hook_call("on_walk_end", config)

def build_if_file(project_path, file):
    if project_path:
        raise ValueError("Please specify either a project path or a file path, not both.")
    elif len(file) != 2:
        raise ValueError("A file path must include exactly a source and a destination.")
    else:
        compile_file(file[0], file[1])

def build(building_type:str,project_path:str, file:list[str]):
    try:
        if file:
            build_if_file(project_path, file)
        else:
            project_path = project_path if project_path else getcwd()
            find_project_from_path(project_path)
            chdir(project_path)
            config = load_config()
            plugins = load_plugins(config)
            hook_call("on_build_start", config)
            walk_and_build(config, plugins)
            hook_call("on_build_end", config)


    except Exception as e:
        fatal(e, f"Build failed: {e}")

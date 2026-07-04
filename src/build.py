import re, sys
import importlib.util
import latex2mathml.converter
from . import settings

from jinja2 import Environment, FileSystemLoader
from markdown import markdown
from markdown.extensions.attr_list import AttrListTreeprocessor
from os import chdir, getcwd, makedirs, path, walk
from os.path import abspath, dirname

from .config import find_project_from_path, load_config
from .errors import fatal, html_fatal
from .fileops import read_file, write_file
from .log import GRAY, info

# patch [ ] instead of { }
AttrListTreeprocessor.BASE_RE   = r'\[\:?[ ]*([^\]\n ][^\n]*)[ ]*\]'
AttrListTreeprocessor.BLOCK_RE  = re.compile(r'\n[ ]*{}[ ]*$'.format(AttrListTreeprocessor.BASE_RE))
AttrListTreeprocessor.HEADER_RE = re.compile(r'[ ]+{}[ ]*$'.format(AttrListTreeprocessor.BASE_RE))
AttrListTreeprocessor.INLINE_RE = re.compile(r'^{}'.format(AttrListTreeprocessor.BASE_RE))

replace_filters = [
    ("<p>{%" ,  "{%" ),
    ("%}</p>",  "%}" ),
]

def load_plugins(path):
    module_dir = dirname(abspath(path))
    sys.path.insert(0, module_dir)
    try:
        for name in list(sys.modules):
            mod_file = getattr(sys.modules[name], "__file__", "") or ""
            if mod_file.startswith(module_dir):
                del sys.modules[name]
        spec = importlib.util.spec_from_file_location("plugins", path)
        module = importlib.util.module_from_spec(spec)
        module.CONFIG = settings.CONFIG
        spec.loader.exec_module(module)

        return {
            name: value
            for name, value in vars(module).items()
            if not name.startswith("_")
        }

    except Exception as e:
        raise RuntimeError(f"Failed loading plugin {path}: {e}") from e
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
    def block(m):
        return latex2mathml.converter.convert(m.group(1))
    def inline(m):
        return latex2mathml.converter.convert(m.group(1))
    
    md_content = re.sub(r'\$\$(.+?)\$\$', block, md_content, flags=re.DOTALL)
    md_content = re.sub(r'\$(.+?)\$', inline, md_content)
    return md_content

def jinja_handler(file, html_content):
    try:
        tree = settings.CONFIG.tree if settings.CONFIG else None
        env = (
            Environment(loader=FileSystemLoader(tree.templates))
            if tree is not None
            else Environment()
        )
        template = env.from_string(html_content)
        return (
            template.render(**load_plugins(tree.plugins))
            if tree is not None
            else template.render()
        )

    except Exception as e:
        return html_fatal(e, f"Template error in {file}")

def escape_code_blocks(html_content: str) -> str:
    """ Neutralize '{' inside rendered <pre>/<code> blocks so Jinja never
    reads {% %} or {{ }} out of code samples, even if a sample's text
    literally contains Jinja syntax (e.g. docs showing `{% raw %}`).
    &#123; renders back to '{' in the browser, so output is unaffected. """
    def neutralize(m):
        return m.group(0).replace("{", "&#123;")

    html_content = re.sub(r'<pre\b[\s\S]*?</pre>', neutralize, html_content)
    html_content = re.sub(r'<code\b[\s\S]*?</code>', neutralize, html_content)
    return html_content

def save_html(html_content:str, html_dest:str):
    makedirs( dirname( abspath(html_dest)), exist_ok=True )
    write_file( html_dest, html_content)

def compile_md_to_html(md_file:str, html_dest:str ):
    """Convert a Markdown file to HTML, applying filters and Jinja2 processing, and save to dest."""

    info(f"Building {GRAY(md_file)}...")
    md_content = read_file(md_file)
    math_rendered = render_math(md_content)
    highlight = "noclasses"
    if settings.CONFIG:
        highlight = settings.CONFIG.extras.highlight
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
            "pymdownx.highlight": {
                "use_pygments": True,
                "noclasses": highlight != "noclasses",
                **(
                    {}
                    if highlight == "noclasses"
                    else {"pygments_style": highlight}
                ),
            }
        }
    )
    escaped_html = escape_code_blocks(raw_html_content)
    filtered_html = html_filter(escaped_html)
    html_content = jinja_handler(md_file, filtered_html)
    save_html(html_content, html_dest)
    info(f"Done: {GRAY(html_dest)}...")

def build(building_type:str,project_path:str, file:list[str]):
    try:
        if file:
            if project_path:
                raise ValueError("Please specify either a project path or a file path, not both.")
 
            elif len(file) != 2:
                raise ValueError("A file path must include exactly a source and a destination.")

            else:
                compile_md_to_html(file[0], file[1].removesuffix(".md") + ".html")
                return

        else:
            project_path = project_path if project_path else getcwd()
            find_project_from_path(project_path)
            chdir(project_path)
            settings.CONFIG = load_config()
            config = settings.CONFIG
            md_path = config.tree.markdown

            for parent, _, files  in walk(md_path):
                for filename in files:
                    md_file = path.join(parent, filename)
                    md_relpath = path.relpath(md_file,md_path)
                    html_dest  = path.join(config.tree.dest, md_relpath).removesuffix(".md") + ".html"
                    compile_md_to_html(md_file, html_dest)

    except Exception as e:
        fatal(e, f"Build failed: {e}")

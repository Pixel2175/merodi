import importlib.util
import re

from jinja2 import Environment, FileSystemLoader
from markdown import markdown
from markdown.extensions.attr_list import AttrListTreeprocessor
from os import getcwd, makedirs, path, walk
from os.path import abspath, dirname

from .config import find_project_from_path, load_tree_config
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
    global context
    spec = importlib.util.spec_from_file_location("plugins", path)
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)

    return {
        name: value
        for name, value in vars(module).items()
        if not name.startswith("_")
    }


def html_filter(html_content:str):
    """ filter jinja placeholders from html """
    html = []
    for line in html_content.splitlines():
        for fltr in replace_filters:
            line = line.replace(*fltr)
        html.append(line)
    return("\n".join(html))

def jinja_handler(file, html_content, config=None):
    try:
        env = Environment( loader=FileSystemLoader(config.templates)) if config else Environment()
        template = env.from_string(html_content)

        if config:
            return template.render(**load_plugins(config.plugins))
        return template.render()
    except Exception as e:
        return html_fatal(e, f"Template error in {file}")

def save_html(html_content:str, html_dest:str):
    makedirs( dirname( abspath(html_dest)), exist_ok=True )
    write_file( html_dest, html_content)

def compile_md_to_html(md_file:str, html_dest:str, config =None):
    """Convert a Markdown file to HTML, applying filters and Jinja2 processing, and save to dest."""
    info(f"Building {GRAY(md_file)}...")
    md_content = read_file(md_file)
    raw_html_content = markdown(
        md_content,
        extensions=[
            "extra",
            "toc",
            "codehilite",
        ]
    )
    filtered_html = html_filter(raw_html_content)
    html_content = jinja_handler(md_file, filtered_html, config)
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
            config  = load_tree_config(project_path)
            md_path = config.markdown

            for parent, _, files  in walk(md_path):
                for filename in files:
                    md_file = path.join(parent, filename)
                    md_relpath = path.relpath(md_file,md_path)
                    html_dest  = path.join(config.dest, md_relpath).removesuffix(".md") + ".html"
                    compile_md_to_html(md_file, html_dest, config)

    except Exception as e:
        fatal(e, f"Build failed: {e}")

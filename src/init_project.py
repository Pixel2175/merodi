from os import path, getcwd, makedirs

from .modules import Config, Project, Tree, Webview, Extras
from .log import *
from .errors import fatal
from .templates import *
from .fileops import *

def init_config_struct(project_name:str) -> Config:
    """Build a Config with default values for a newly initialized project."""
    return Config(
        project=Project(
            name        = project_name,
            version     = "0.1.0",
            description = "Add your description here",
        ),
        tree = Tree(
            markdown  = "src/md",
            static    = "src/static",
            templates = "src/templates",
            dest      = "src/dest",
            plugins   = "src/plugins.py",
        ),
        webview = Webview(
            host        = "localhost",
            port        = 8866,
            dev_tools   = False,
            html_path   = "/",
            static_path = "/static"
        ),
        extras = Extras(
            highlight = "monokai"
        )
    )

def render_config(config:Config) -> str:
    return f"""[project]
name        = "{config.project.name}"
version     = "{config.project.version}"
description = "{config.project.description}"

[tree]
markdown  = "{config.tree.markdown}"
static    = "{config.tree.static}"
templates = "{config.tree.templates}"
dest      = "{config.tree.dest}"
plugins   = "{config.tree.plugins}"

[webview]
host        = "{config.webview.host}"
port        = {config.webview.port}
dev_tools   = {str(config.webview.dev_tools).lower()}
html_path   = "{config.webview.html_path}"
static_path = "{config.webview.static_path}"

[extras]
highlight = "{config.extras.highlight}"
"""

def check_not_already_initialized(file_path:str) -> None:
    """raise FileExistsError when config file exists"""
    if path.exists(file_path):
        raise Exception("Project already initialized: `config.toml` already exists")

def write_default_content(tree:Tree) -> None:
    """ Write default content """
    md_file = path.join(tree.markdown, "index.md")
    write_file(md_file,MARKDOWN_CONTENT )

    html_file = path.join(tree.templates, "layout.html")
    write_file(html_file,HTML_CONTENT )

    css_file = path.join(tree.static, "style.css")
    write_file(css_file,CSS_CONTENT )

    write_file(tree.plugins,PLUGINS_CONTENT )

def write_config_file(project_path, config_content):
    makedirs(project_path, exist_ok=True)
    config_path = path.join(project_path, "config.toml")
    check_not_already_initialized(config_path)

    # Writing config.toml
    write_file(config_path, config_content)


def init(project_path):
    try:
        project_path = project_path if project_path else getcwd()
        config = init_config_struct(path.basename(project_path) ) 
        default_config_content = render_config(config)
        config.tree = resolve_tree_paths(project_path, config.tree)
        write_config_file(project_path, default_config_content)

        create_tree_dirs(config.tree)
        write_default_content(config.tree)
        info("Project initialized.")
    except Exception as e:
        fatal(e, f"Initialization failed: {e}")


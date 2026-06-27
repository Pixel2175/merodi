from os import makedirs, path
from os.path import exists
from tomllib import loads

from .fileops import read_file
from .log import GRAY
from .modules import Config, Extras, Project, Tree, Webview

def init_config_struct(project_name:str) -> Config:
    """Build a Config with default values for a newly initialized project."""
    return Config(
        project=Project(
            name = project_name,
            version = "0.1.0",
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
            host = "localhost",
            port = 8866,
            html_path   = "/",
            static_path = "/static"
        ),
        extras= Extras(
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
host  = "{config.webview.host}"
port  = {config.webview.port}
html_path   = "{config.webview.html_path}"
static_path = "{config.webview.static_path}"
"""

def find_project_from_path(project_path: str):
    if not path.exists(project_path):
        raise Exception("project does not exist")
    if not path.exists(path.join(project_path, "config.toml")):
        raise Exception(f"cannot find `config.toml` in: {GRAY(project_path)}")

def load_tree_config(project_path) -> Tree:
    config_raw_content  = read_file(path.join(project_path,"config.toml"))
    config = loads(config_raw_content) 
    markdown  = path.join(project_path, config["tree"]["markdown"])
    static    = path.join(project_path, config["tree"]["static"])
    templates = path.join(project_path, config["tree"]["templates"])
    dest      = path.join(project_path, config["tree"]["dest"])
    plugins   = path.join(project_path, config["tree"]["plugins"])

    if not exists(markdown):
        raise FileNotFoundError(f"markdown directory does not exist: {markdown}")

    if not exists(static):
        raise FileNotFoundError(f"static directory does not exist: {static}")

    if not exists(templates):
        raise FileNotFoundError(f"templates directory does not exist: {templates}")

    if not exists(plugins):
        raise FileNotFoundError(f"plugins file does not exist: {plugins}")

    if not exists(dest):
        makedirs(dest)

    return Tree(
        markdown  = markdown,
        static    = static,
        templates = templates,
        dest      = dest,
        plugins   = plugins,
    )


def load_webview_config(project_path:str) -> Webview:
    config_raw_content  = read_file(path.join(project_path,"config.toml"))
    config = loads(config_raw_content) 
    return Webview(
        host        = config["webview"]["host"],
        port        = config["webview"]["port"],
        html_path   = config["webview"]["html_path"],
        static_path = config["webview"]["static_path"],
    )

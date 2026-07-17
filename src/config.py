from os import makedirs, path
from os.path import exists
from tomllib import loads
from .fileops import read_file
from .log import GRAY
from .modules import Cache, Config, Extras, Project, Tree, Webview
from .hooks import hook_call


def find_project_from_path(project_path: str):
    if not path.exists(project_path):
        raise Exception("project does not exist")
    if not path.exists(path.join(project_path, "config.toml")):
        raise Exception(f"cannot find `config.toml` in: {GRAY(project_path)}")

def load_project_config(config) -> Project:
    return Project(
        name        = config["project"].get("name", "project"),
        version     = config["project"].get("version", "0.1.0"),
        description = config["project"].get("description", "Add your description here"),
    ) 

def load_tree_config(config) -> Tree:
    markdown  = config["tree"].get("markdown",  "src/md")
    static    = config["tree"].get("static",    "src/static")
    templates = config["tree"].get("templates", "src/templates")
    dest      = config["tree"].get("dest",      "src/dest")
    plugins   = config["tree"].get("plugins",   "src/plugins.py")

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

def load_webview_config(config) -> Webview:
    return Webview(
        host        = config["webview"].get("host",        "localhost"),
        port        = config["webview"].get("port",        8866),
        dev_tools   = config["webview"].get("dev_tools",   False),
        html_path   = config["webview"].get("html_path",   "/"),
        static_path = config["webview"].get("static_path", "/static"),
    )

def load_extras_config(config) -> Extras:
    return Extras(
        highlight = config["extras"].get("highlight", "monokai"),
    )

def load_cache_config(config) -> Cache:
    return Cache(
        hash = config["cache"].get("hash", ".hash.cache"),
    )

def load_config() -> Config:
    config_raw_content = read_file("config.toml")
    config_loaded = loads(config_raw_content)
    config = Config(
        project = load_project_config(config_loaded),
        tree    = load_tree_config(config_loaded),
        webview = load_webview_config(config_loaded),
        extras  = load_extras_config(config_loaded),
        cache   = load_cache_config(config_loaded),
    )
    config = hook_call("on_config_load", config) or config
    return config


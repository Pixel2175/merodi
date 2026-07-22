from os import makedirs, path
from os.path import exists
from tomllib import loads
from .fileops import read_file
from .log import GRAY
from .modules import Cache, Config, Extras, Project, Tree, Webview
from .hooks import hook_call
from .api import api


def find_project_from_path(project_path):
    if not path.exists(project_path):
        raise Exception("project does not exist")
    if not path.exists(path.join(project_path, "config.toml")):
        raise Exception(f"cannot find `config.toml` in: {GRAY(project_path)}")

def load_project_config(c) -> Project:
    return Project(name=c["project"]["name"], version=c["project"]["version"], description=c["project"]["description"])

def load_tree_config(c) -> Tree:
    t = c["tree"]
    markdown, static, templates = t["markdown"], t["static"], t["templates"]
    draft_dest, release_dest, plugins = t["draft_dest"], t["release_dest"], t["plugins"]
    for name, p in [("markdown", markdown), ("static", static), ("templates", templates), ("plugins", plugins)]:
        if not exists(p):
            raise FileNotFoundError(f"{name} directory does not exist: {p}")
    return Tree(markdown=markdown, static=static, templates=templates, draft_dest=draft_dest, release_dest=release_dest, plugins=plugins)

def load_webview_config(c) -> Webview:
    w = c["webview"]
    return Webview(host=w["host"], port=w["port"], dev_tools=w["dev_tools"], html_path=w["html_path"], static_path=w["static_path"])

def load_extras_config(c) -> Extras:
    return Extras(highlight=c["extras"]["highlight"])

def load_cache_config(c) -> Cache:
    return Cache(hash=c["cache"]["hash"])

def load_config() -> Config:
    raw = read_file("config.toml")
    c = loads(raw)
    config = Config(
        project=load_project_config(c),
        tree=load_tree_config(c),
        webview=load_webview_config(c),
        extras=load_extras_config(c),
        cache=load_cache_config(c),
    )
    config = hook_call("on_config_load", config) or config
    api.config = config
    return config

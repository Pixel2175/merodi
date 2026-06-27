from dataclasses import dataclass

@dataclass()
class Project:
    name        : str
    version     : str
    description : str

@dataclass()
class Tree:
    markdown  : str
    static    : str
    templates : str
    dest      : str
    plugins   : str

@dataclass()
class Webview:
    host        : str
    port        : int
    dev_tools   : str
    html_path   : str
    static_path : str

@dataclass
class Extras:
    highlight :str 

@dataclass()
class Config:
    project:Project
    tree:Tree
    webview:Webview
    extras:Extras

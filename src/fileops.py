from os import makedirs, path, sep
from .modules import Tree

def is_dotfile(file_path):
    return any(part.startswith(".") for part in path.normpath(file_path).split(sep))

def write_file(file_path: str, file_content: str) -> None:
    with open(file_path, "w") as f:
        f.write(file_content)

def read_file(file_path:str) -> str:
    """Read and return file contents"""
    with open(file_path, "r") as f :
        return f.read()

def resolve_tree_paths(project_dir: str, tree: Tree) -> Tree:
    """Return a new Tree with all paths joined onto project_dir."""
    return Tree(
        markdown=path.join(project_dir, tree.markdown),
        static=path.join(project_dir, tree.static),
        templates=path.join(project_dir, tree.templates),
        dest=path.join(project_dir, tree.dest),
        plugins=path.join(project_dir, tree.plugins),
    )

def create_tree_dirs (tree:Tree) -> None:
    """ Creates parent dirs before start writing """
    makedirs(tree.markdown , exist_ok=True)
    makedirs(tree.static   , exist_ok=True)
    makedirs(tree.templates, exist_ok=True)
    makedirs(tree.dest     , exist_ok=True)
    makedirs(tree.plugins  , exist_ok=True)

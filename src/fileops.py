from os import path
from .modules import Tree

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


from os import makedirs, path, sep
from .modules import Tree

def is_dotfile(file_path):
    return any(part.startswith(".") for part in path.normpath(file_path).split(sep))

def write_file(file_path, file_content):
    with open(file_path, "w") as f:
        f.write(file_content)

def read_file(file_path):
    with open(file_path, "r") as f:
        return f.read()

def resolve_tree_paths(project_dir, tree):
    return Tree(
        markdown=path.join(project_dir, tree.markdown),
        static=path.join(project_dir, tree.static),
        templates=path.join(project_dir, tree.templates),
        draft_dest=path.join(project_dir, tree.draft_dest),
        release_dest=path.join(project_dir, tree.release_dest),
        plugins=path.join(project_dir, tree.plugins),
    )

def create_tree_dirs(tree):
    for d in [tree.markdown, tree.static, tree.templates, tree.draft_dest, tree.release_dest, tree.plugins]:
        makedirs(d, exist_ok=True)

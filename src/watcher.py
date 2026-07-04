from os import  chdir, getcwd, path
from typing import Callable
from watchdog.observers import Observer
from watchdog.events import FileSystemEventHandler
import time

from src import settings
from .config import find_project_from_path, load_config
from .build import compile_md_to_html
from .log import  warn
from .errors import fatal

current_md_file = None

def reload(changed_path, tree_config, extras_config):
    global current_md_file
    try:
        if changed_path.endswith(".md"):
            current_md_file = changed_path
        if current_md_file is None:
            return None

        md_relpath = path.relpath(current_md_file, tree_config.markdown)
        html = path.splitext(md_relpath)[0] + ".html"
        dest = path.join(tree_config.dest, html)
        compile_md_to_html(current_md_file, dest)
        return html
    except Exception as e:
        warn(str(e))

def watch_files(reload_func:Callable | None=None):
    last_reload = {}
    config = settings.CONFIG
    
    class ReloadHandler(FileSystemEventHandler):
        def on_modified(self, event):
            reload_path = event.src_path
            if event.is_directory :
                    return
            now = time.time()
            last = last_reload.get(reload_path, 0)
            if now - last < 0.3:
                return
            last_reload[reload_path] = now
            file = reload( reload_path, config.tree , config.extras)
            if reload_func and file: reload_func(file)
    observer = Observer()
    handler = ReloadHandler()
    for file in [config.tree.markdown, config.tree.templates, path.dirname(config.tree.plugins), config.tree.static]:
        observer.schedule(
            handler,
            path=file,
            recursive=True
        )
    return observer

def run_watcher(project_path):
    try:
        project_path = project_path if project_path else getcwd()
        find_project_from_path(project_path)
        chdir(project_path)
        settings.CONFIG = load_config()
        observers = watch_files()
        observers.start()
        try:
            while True:
                time.sleep(1)
        except KeyboardInterrupt:
            observers.stop()
        observers.join()

    except Exception as e:
        fatal(e, f"Build failed: {e}")

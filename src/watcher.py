from os import  getcwd, path
from typing import Callable
from watchdog.observers import Observer
from watchdog.events import FileSystemEventHandler
import time
from .config import find_project_from_path, load_extras_config, load_tree_config 
from .build import compile_md_to_html
from .log import  warn
from .errors import fatal

def reload( file, tree_config, extras_config):
    global current_url
    try:
        md_relpath = path.relpath(file, tree_config.markdown)
        dest = path.join(tree_config.dest, md_relpath)
        if dest.endswith(".md"):
            dest = dest.removesuffix(".md") + ".html"
            url = path.relpath(dest, tree_config.dest)
            compile_md_to_html(file, dest, tree_config, extras_config)
            return url
    except Exception as e:
        warn(str(e))

def watch_files(tree_config, extra_config, reload_func:Callable | None=None):
    last_reload = {}
    
    class ReloadHandler(FileSystemEventHandler):
        def on_modified(self, event):
            reload_path = event.src_path
            if event.is_directory :
                if path.exists(path.join(event.src_path, "index.html")):
                    reload_path = path.join(event.src_path, "index.html")
                else: 
                    return
            now = time.time()
            last = last_reload.get(reload_path, 0)
            if now - last < 0.3:
                return
            last_reload[reload_path] = now
            file = reload( reload_path, tree_config, extra_config)
            if reload_func: reload_func(file)



    observer = Observer()
    handler = ReloadHandler()
    for file in [tree_config.markdown, tree_config.templates, tree_config.plugins, tree_config.static]:
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
        extra_config = load_extras_config(project_path)
        tree_config  = load_tree_config(project_path)
        observers = watch_files(tree_config, extra_config)
        observers.start()
        try:
            while True:
                time.sleep(1)
        except KeyboardInterrupt:
            observers.stop()
        observers.join()

    except Exception as e:
        fatal(e, f"Build failed: {e}")

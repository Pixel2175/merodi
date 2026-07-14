from os import  chdir, getcwd, path
from typing import Callable
from watchdog.observers import Observer
from watchdog.events import FileSystemEventHandler
import time

from .config import find_project_from_path, load_config
from .build import compile_file, load_plugins
from .log import warn
from .errors import fatal

current_md_file = None

def reload(changed_path, config, plugins, force=None):
    global current_md_file
    try:
        if changed_path.endswith(".md"):
            current_md_file = changed_path
        if current_md_file is None:
            return None

        md_relpath = path.relpath(current_md_file, config.tree.markdown)
        html = path.splitext(md_relpath)[0] + ".html"
        dest = path.join(config.tree.dest, html)
        compile_status = compile_file(current_md_file, dest, config, plugins, force=force)
        if compile_status is None:
            return None
        return html
    except Exception as e:
        warn(str(e))

def watch_files(config, reload_func:Callable | None=None):
    last_reload = {}
    plugins = load_plugins(config)
    class ReloadHandler(FileSystemEventHandler):
        def on_modified(self, event):
            nonlocal plugins
            reload_path = event.src_path
            force = False
            if reload_path.startswith(config.tree.templates):
                force = True
            if reload_path.startswith(config.tree.plugins):
                plugins = load_plugins(config)
            if event.is_directory :
                    return

            now = time.time()
            last = last_reload.get(reload_path, 0)
            if now - last < 0.3:
                return
            last_reload[reload_path] = now
            file = reload( reload_path, config, plugins, force=force)
            if reload_func and file: reload_func(file)
    observer = Observer()
    handler = ReloadHandler()
    for file in [config.tree.markdown, config.tree.templates, config.tree.plugins, config.tree.static]:
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
        config = load_config()
        observers = watch_files(config)
        observers.start()
        try:
            while True:
                time.sleep(1)
        except KeyboardInterrupt:
            observers.stop()
        observers.join()

    except Exception as e:
        fatal(e, f"Build failed: {e}")

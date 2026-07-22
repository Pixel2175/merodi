from os import chdir, getcwd, path
from watchdog.observers import Observer
from watchdog.events import FileSystemEventHandler
import time

from src.hooks import hook_call
from .config import find_project_from_path, load_config
from .build import compile_file, load_plugins, build_all
from .log import warn
from .errors import fatal

current_page = None
_plugins = None

def page_to_url(filesystem_path, config):
    rel = path.relpath(filesystem_path, config.tree.markdown)
    return path.splitext(rel)[0] + ".html"

def is_template_or_plugin(file_path, config):
    return (file_path.startswith(config.tree.templates) or
            file_path.startswith(config.tree.plugins))

def single_file_reload(changed_path, config):
    global current_page
    try:
        if changed_path.endswith(".md"):
            current_page = changed_path
        if current_page is None:
            return None
        html = page_to_url(current_page, config)
        dest = path.join(config.tree.draft_dest, html)
        if compile_file(current_page, dest, config, _plugins) is None:
            return None
        return html
    except Exception as e:
        hook_call("on_reload_error", e)
        warn(str(e))

def full_reload(config, plugins):
    global current_page
    try:
        build_all(config, plugins, config.tree.draft_dest)
        return page_to_url(current_page, config) if current_page else None
    except Exception as e:
        hook_call("on_reload_error", e)
        warn(str(e))

def watch_files(config, reload_func=None):
    global _plugins
    last_reload = {}
    _plugins = load_plugins(config)

    class ReloadHandler(FileSystemEventHandler):
        def on_modified(self, event):
            global _plugins
            try:
                p = event.src_path
                if event.is_directory:
                    return
                if p.startswith(config.tree.static):
                    if reload_func:
                        reload_func("")
                    return
                if is_template_or_plugin(p, config):
                    if p.startswith(config.tree.plugins):
                        _plugins = load_plugins(config)
                        hook_call("on_plugin_changed", p)
                    file = full_reload(config, _plugins)
                    if reload_func and file:
                        reload_func(file)
                    return
                now = time.time()
                if now - last_reload.get(p, 0) < 0.3:
                    return
                last_reload[p] = now
                hook_call("on_file_changed", p, config)
                file = single_file_reload(p, config)
                if reload_func and file:
                    reload_func(file)
            except Exception as e:
                hook_call("on_reload_error", e)
                warn(str(e))

    observer = Observer()
    handler = ReloadHandler()
    for dir_path in [config.tree.markdown, config.tree.templates, config.tree.plugins, config.tree.static]:
        observer.schedule(handler, path=dir_path, recursive=True)
    hook_call("on_watch_start", config)
    return observer

def run_watcher(project_path):
    try:
        project_path = project_path or getcwd()
        find_project_from_path(project_path)
        chdir(project_path)
        config = load_config()
        observer = watch_files(config)
        observer.start()
        try:
            while True:
                time.sleep(1)
        except KeyboardInterrupt:
            observer.stop()
            hook_call("on_watch_stop", config)
        observer.join()
    except Exception as e:
        fatal(e, f"Watch failed: {e}")

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

        md_relpath = path.relpath(current_page, config.tree.markdown)
        html = path.splitext(md_relpath)[0] + ".html"
        dest = path.join(config.tree.dest, html)
        compile_status = compile_file(current_page, dest, config, _plugins)
        if compile_status is None:
            return None
        return html
    except Exception as e:
        hook_call("on_reload_error", e)
        warn(str(e))

def full_reload(config, plugins):
    global current_page
    try:
        build_all(config, plugins)
        return current_page
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
                reload_path = event.src_path
                if event.is_directory:
                    return
                if reload_path.startswith(config.tree.static):
                    if reload_func:
                        reload_func("")
                    return

                if is_template_or_plugin(reload_path, config):
                    if reload_path.startswith(config.tree.plugins):
                        _plugins = load_plugins(config)
                        hook_call("on_plugin_changed", reload_path)
                    file = full_reload(config, _plugins)
                    if reload_func and file:
                        reload_func(file)
                    return

                now = time.time()
                last = last_reload.get(reload_path, 0)
                if now - last < 0.3:
                    return
                last_reload[reload_path] = now
                hook_call("on_file_changed", reload_path, config)
                file = single_file_reload(reload_path, config)
                if reload_func and file:
                    reload_func(file)
            except Exception as e:
                hook_call("on_reload_error", e)
                warn(str(e))

    observer = Observer()
    handler = ReloadHandler()
    for file in [config.tree.markdown, config.tree.templates, config.tree.plugins, config.tree.static]:
        observer.schedule(
            handler,
            path=file,
            recursive=True
        )
    hook_call("on_watch_start", config)
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
            hook_call("on_watch_stop", config)
        observers.join()

    except Exception as e:
        fatal(e, f"Build failed: {e}")

from os import getcwd, path
from http.server import HTTPServer, SimpleHTTPRequestHandler
from threading import Thread
from watchdog.observers import Observer
from watchdog.events import FileSystemEventHandler
import webview
import time

from .config import find_project_from_path, load_tree_config, load_webview_config
from .errors import fatal, html_fatal
from .build import compile_md_to_html
from .log import info

current_url = ""

def http_server(host, port, routes):
    class Handler(SimpleHTTPRequestHandler):
        def translate_path(self, url_path: str) -> str:
            url = routes["url_path"]
            fs = routes["fs_path"]

            if url_path.startswith(url["static"]):
                rel = url_path.removeprefix(url["static"]).lstrip("/")
                resolved = path.join(fs["static"], rel)

            elif url_path.startswith(url["html"]):
                rel = url_path.removeprefix(url["html"]).lstrip("/")
                resolved = path.join(fs["html"], rel)

            else:
                rel = url_path.lstrip("/")
                resolved = path.join(fs["html"], rel)

            if not path.exists(resolved) and path.exists(resolved + ".html"):
                resolved += ".html"

            if path.isdir(resolved):
                index = path.join(resolved, "index.html")
                if path.exists(index):
                    resolved = index

            return resolved

        def send_error(self, code, message=None, explain=None):
            html = html_fatal(Exception(f"{code}: {message}"), explain or self.path)

            self.send_response(code)
            self.send_header("Content-Type", "text/html")
            self.send_header("Content-Length", len(html))
            self.end_headers()

            self.wfile.write(html.encode())

        def log_message(self, format, *args):
            method, path, _ = args[0].split(" ", 2)
            info(f"{method} {path} {args[1]}")
                    
    server = HTTPServer((host, port), Handler)
    return server

def reload(window, file, tree_config, webview_config):
    global current_url
    try:
        md_relpath = path.relpath(file, tree_config.markdown)
        dest = path.join(tree_config.dest, md_relpath)
        if dest.endswith(".md"):
            dest = dest.removesuffix(".md") + ".html"
            current_url = path.relpath(dest, tree_config.dest)
            compile_md_to_html(file, dest, tree_config)
        info(f"Reloading: {current_url}")
        if current_url:
            window.load_url(f"http://{webview_config.host}:{webview_config.port}/{current_url}")
    except Exception as e:
        window.load_url(f"data:text/html,{html_fatal(e, f'Failed to reload {file}')}")

def watch_files(window, tree_config, webview_config):
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
            reload(window, reload_path, tree_config, webview_config)


    observer = Observer()
    handler = ReloadHandler()
    for file in [tree_config.markdown, tree_config.templates, tree_config.plugins, tree_config.static]:
        observer.schedule(
            handler,
            path=file,
            recursive=True
        )
    observer.start()

def run(project_path):
    try:
        project_path = project_path if project_path else getcwd()
        find_project_from_path(project_path)
        webview_config = load_webview_config(project_path)
        host = webview_config.host
        port = webview_config.port
        tree_config    = load_tree_config(project_path)
        routes = {
            "url_path" : {
                "html":   webview_config.html_path,
                "static": webview_config.static_path
            },
            "fs_path"  : {
                "html":   tree_config.dest,
                "static": tree_config.static
            },
        }

        server = http_server(host, port, routes)
        Thread(target=server.serve_forever, daemon=True).start()

        window = webview.create_window("Merodi", url=f"http://{host}:{port}/")
        watch_files(window, tree_config, webview_config )
        webview.start()


    except KeyboardInterrupt:
        info("Exiting...")
        exit(0)
    except Exception as e:
        fatal(e, f"Webview failed: {e}")

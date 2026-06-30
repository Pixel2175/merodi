from os import getcwd, path
from http.server import HTTPServer, SimpleHTTPRequestHandler
from threading import Thread
from watchdog.observers import Observer
from watchdog.events import FileSystemEventHandler
import webview
import time

from .config import find_project_from_path, load_extras_config, load_tree_config, load_webview_config
from .errors import fatal, html_fatal
from .build import compile_md_to_html
from .watcher import watch_files 
from .log import die, info, warn


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
            try:
                method, path, _ = args[0].split(" ", 2)
                info(f"{method} {path} {args[1]}")
            except (ValueError, IndexError):
                info(format % args)
                    
    server = HTTPServer((host, port), Handler)
    return server

def reload_webview(window, url, webview_config):
    try:
        window.load_url(f"http://{webview_config.host}:{webview_config.port}/{url}")
    except Exception as e:
        window.load_url(f"data:text/html,{html_fatal(e, f'Failed to reload {url}')}")
        warn(str(e))

def run(project_path):
    try:
        project_path = project_path if project_path else getcwd()
        find_project_from_path(project_path)
        webview_config = load_webview_config(project_path)
        host = webview_config.host
        port = webview_config.port
        tree_config    = load_tree_config(project_path)
        extras_config = load_extras_config(project_path)
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
        server_thread = Thread(target=server.serve_forever, daemon=True)
        server_thread.start()

        window = webview.create_window("Merodi", url=f"http://{host}:{port}/")
        observer = watch_files(
                tree_config, extras_config,
                lambda path: reload_webview(window, path, webview_config),
                )
        observer.start()
        webview.start(debug=webview_config.dev_tools in ["true",1])


    except KeyboardInterrupt:
        info("Exiting...")
        exit(0)
    except Exception as e:
        fatal(e, f"Webview failed: {e}")

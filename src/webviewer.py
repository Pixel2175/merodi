from os import getcwd, path
from http.server import HTTPServer, SimpleHTTPRequestHandler
from threading import Thread
import webview

from .config import find_project_from_path, load_extras_config, load_tree_config, load_webview_config
from .errors import fatal, html_fatal
from .watcher import watch_files 
from .log import info, warn
import ifaddr
import socket

_ip_cache: list[str] = []

import ifaddr

def resolve_ips(host: str, port: int) -> list[str]:
    global _ip_cache
    if host != "0.0.0.0":
        _ip_cache = [host]
    elif not _ip_cache:
        _ip_cache = [
            ip_str
            for adapter in ifaddr.get_adapters()
            for ip in adapter.ips
            if ip.is_IPv4
            and (ip_str := str(ip.ip))
            and not ip_str.startswith("127.")
            and not ip_str.startswith("169.254.")
        ] or ["127.0.0.1"]

    return _ip_cache

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

def reload_webview(window, url, host, port):
    try:
        window.load_url(f"http://{host}:{port}/{url}")
    except Exception as e:
        window.load_url(f"data:text/html,{html_fatal(e, f'Failed to reload {url}')}")
        warn(str(e))

def run(project_path):
    try:
        global _ip_cache
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
        _ip_cache = resolve_ips(host, port)
        info(f"HTTP SERVER STARTS ON:")
        for ip in _ip_cache:
            info(f"  IP: http://{ip}:{port}")
        print()

        window = webview.create_window("Merodi", url=f"http://{_ip_cache[0]}:{port}/")
        observer = watch_files(
                tree_config, extras_config,
                lambda path: reload_webview(window, path, _ip_cache[0], port),
                )
        observer.start()
        webview.start(debug=webview_config.dev_tools in ["true",1])


    except KeyboardInterrupt:
        info("Exiting...")
        exit(0)
    except Exception as e:
        fatal(e, f"Webview failed: {e}")

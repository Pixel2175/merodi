import sys
sys.dont_write_bytecode = True

import atexit
from argparse import ArgumentParser
from os import environ
from .log import *
from . import settings
from .hooks import hook_call

def main():
    common = ArgumentParser(add_help=False)
    common.add_argument("--verbose", action="store_true")
    common.add_argument("--no-color", action="store_true")

    parser = ArgumentParser()
    parser.add_argument("--version", action="store_true")
    sub_parser = parser.add_subparsers(dest="command")

    init_parser = sub_parser.add_parser("init", parents=[common])
    init_parser.add_argument("path", nargs="?")

    webview_parser = sub_parser.add_parser("webview", parents=[common])
    webview_parser.add_argument("path", nargs="?")

    webview_parser = sub_parser.add_parser("watch", parents=[common])
    webview_parser.add_argument("path", nargs="?")

    build_parser = sub_parser.add_parser("build", parents=[common])
    build_mode = build_parser.add_mutually_exclusive_group()
    build_mode.add_argument("--draft", action="store_true", default=True)
    build_mode.add_argument("--release", action="store_true")
    build_parser.add_argument("path", nargs="?")
    build_parser.add_argument("--file", nargs=2, metavar=("SRC", "DEST"))

    args = parser.parse_args()

    if args.version:
        from importlib.metadata import version
        info(version("merodi"), title="VERSION")
        return

    if args.command is None:
        parser.print_help()
        return

    settings.VERBOSE  =  args.verbose  or ( environ.get("VERBOSE")  in ["true", "1"])
    settings.NO_COLOR =  args.no_color or ( environ.get("NO_COLOR") in ["true", "1"])

    hook_call("on_start")

    if args.command == "init":
        from .init_project import init
        init(project_path = args.path)

    elif args.command == "build":
        from .api import api
        mode = "release" if args.release else "draft"
        api.mode = mode
        from .build import build
        build(mode = mode, project_path = args.path, file=args.file)

    elif args.command == "webview":
        from .webviewer import run
        run(args.path)

    elif args.command == "watch":
        from .watcher import run_watcher
        run_watcher(args.path)

    atexit.register(lambda: hook_call("on_end"))

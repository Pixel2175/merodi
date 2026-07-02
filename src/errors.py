from os.path import basename
from html import escape
import traceback

from .log import *
from . import settings

def get_user_frame(exc: Exception):
    tb = traceback.extract_tb(exc.__traceback__)
    for frame in reversed(tb):
        if "site-packages" not in frame.filename and "<frozen" not in frame.filename:
            return frame
    if not tb:
        return None
    return tb[-1]

def get_error_line(exc: Exception):
    if hasattr(exc, 'source') and exc.source:
        lines = exc.source.splitlines()
        lineno = getattr(exc, 'lineno', 1) or 1
        if lineno <= len(lines):
            return lines[lineno - 1].strip()
    frame = get_user_frame(exc)
    if not frame:
        return None
    return frame.line.strip() if frame.line else None

def fatal(exc: Exception, message: str):
    frame = get_user_frame(exc)
    if settings.VERBOSE and frame:
        print(f"{BLUE(frame.name + '()')} {GRAY('—')} {basename(frame.filename)}:{YELLOW(str(frame.lineno))}")
        print(f"{GRAY('│')}  {frame.line}\n")
    die(message)

def html_fatal(exc: Exception, message: str) -> str:
    frame = get_user_frame(exc)
    source_line = get_error_line(exc)
    warn(escape(message))
    warn(f"{type(exc).__name__}: {exc.message if hasattr(exc, 'message') else exc}")
    line = escape(source_line) if source_line else "<i style='color:#666'>line not available</i>"
    detail = escape(str(exc.message)) if hasattr(exc, 'message') else escape(str(exc))

    if frame is not None:
        location_block = f"""<div style="background: #2a2a2a; border-radius: 6px; overflow: hidden;">
                <div style="background: #333; padding: 0.5rem 1rem; display: flex; justify-content: space-between;">
                    <span style="color: #569cd6;">{escape(frame.name)}()</span>
                    <span style="color: #888;">{escape(basename(frame.filename))}:{frame.lineno}</span>
                </div>
                <pre style="margin: 0; padding: 1rem; color: #ddd; white-space: pre-wrap; border-left: 3px solid #f44; font-size: 1.1rem;">{line}</pre>
            </div>"""
    else:
        location_block = """<div style="background: #2a2a2a; border-radius: 6px; overflow: hidden;">
                <div style="background: #333; padding: 0.5rem 1rem;">
                    <span style="color: #888;">no source location available</span>
                </div>
            </div>"""

    return f"""<html>
    <head><meta charset="UTF-8"></head>
    <body style="font-family: monospace; margin: 0; background: #1e1e1e; color: #ccc; min-height: 100vh; font-size: 18px;">
        <div style="background: #f44; padding: 0.75rem 2rem;">
            <span style="color: #fff; font-size: 1.2rem; font-weight: bold;">! {type(exc).__name__}</span>
        </div>
        <div style="padding: 2rem;">
            <p style="color: #f90; font-size: 1.3rem; margin: 0 0 0.25rem;">{detail}</p>
            <p style="color: #888; font-size: 1.2rem; margin: 0 0 2rem;">{escape(message)}</p>
            {location_block}
        </div>
    </body>
    </html>"""

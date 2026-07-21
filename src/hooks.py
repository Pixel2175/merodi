from .log import GRAY, warn

DEFAULT_PRIORITY = 100

hooks = {
    "on_build_end"              : [],
    "on_build_start"            : [],
    "on_compile_start"          : [],
    "on_config_load"            : [],
    "on_end"                    : [],
    "on_file_changed"           : [],
    "on_hash_check"             : [],
    "on_hash_written"           : [],
    "on_highlight_config"       : [],
    "on_html_filter"            : [],
    "on_jinja_error"            : [],
    "on_jinja_template_dir"     : [],
    "on_jinja_template_renderer": [],
    "on_math_renderer"          : [],
    "on_md_to_html"             : [],
    "on_page_built"             : [],
    "on_page_read"              : [],
    "on_page_rendered"          : [],
    "on_page_skip"              : [],
    "on_page_written"           : [],
    "on_plugin_changed"         : [],
    "on_plugins_before_export"  : [],
    "on_plugins_loaded"         : [],
    "on_reload_error"           : [],
    "on_start"                  : [],
    "on_walk_end"               : [],
    "on_walk_file"              : [],
    "on_walk_start"             : [],
    "on_watch_start"            : [],
    "on_watch_stop"             : [],
}

def hook(hook_name: str, priority: int = DEFAULT_PRIORITY):
    if hook_name not in hooks:
        warn(f"Skipping hook `{GRAY(hook_name)}`: hook does not exist")
        return lambda fn: fn

    def hook_fn(fn):
        hooks[hook_name].append((priority, fn))
        hooks[hook_name].sort(key=lambda x: x[0])
        return fn
    return hook_fn

def hook_call(hook_name, *args, **kwargs):
    subscribers = hooks.get(hook_name, [])
    if not subscribers:
        return None
    result = None
    for _priority, fn in subscribers:
        try:
            result = fn(*args, **kwargs)
        except Exception as e:
            warn(f"Hook '{hook_name}' failed in {fn.__name__}: {e}")
    return result

def reset_hooks():
    for key in hooks:
        hooks[key] = []

def remove_hook(hook_name, fn):
    if hook_name in hooks:
        hooks[hook_name] = [(p, f) for p, f in hooks[hook_name] if f is not fn]

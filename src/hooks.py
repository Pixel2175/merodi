from .log import GRAY, warn

hooks = {
    "on_after_compile_file"  : None,
    "on_after_html_filter"   : None,      
    "on_after_jinja"         : None,
    "on_after_math"          : None,
    "on_after_escape_code"   : None,
    "on_after_md_to_html"    : None,
    "on_after_markdown"      : None,
    "on_after_template"      : None,  
    "on_before_escape_code"  : None,
    "on_before_compile_file" : None,
    "on_before_html_filter"  : None,      
    "on_before_markdown"     : None,  
    "on_before_math"         : None,  
    "on_before_md_to_html"   : None,
    "on_before_save"         : None,
    "on_before_template"     : None,  
    "on_build_end"           : None,
    "on_build_start"         : None,  
    "on_config_load"         : None,  
    "on_file_changed"        : None,  
    "on_hash_check"          : None,  
    "on_hash_written"        : None,
    "on_highlight_config"    : None,  
    "on_page_read"           : None,  
    "on_page_skip"           : None,      
    "on_page_written"        : None,      
    "on_plugins_loaded"      : None,      
    "on_reload_error"        : None,  
    "on_template_error"      : None,      
    "on_walk_end"            : None,      
    "on_walk_file"           : None,  
    "on_walk_start"          : None,  
    "on_watch_start"         : None,  
    "on_watch_stop"          : None,  
}                            
                             
def hook(hook_name:str):
    if hook_name in hooks:
        def hook_fn(fn):
            if not hooks.get(hook_name) == None:
                warn(f"Overwriting existing hook '{hook_name}'.")
            hooks[hook_name]=fn
            return fn
        return hook_fn

    else:
        warn(f"Skipping hook `{GRAY(hook_name)}`: hook is not exists")
        return lambda fn:fn

def hook_call(hook_name, *args, **kwargs):
    if hook_name not in  hooks:
        warn(hook_name)
    func = hooks.get(hook_name)
    if callable(func):
        return func(*args, **kwargs)

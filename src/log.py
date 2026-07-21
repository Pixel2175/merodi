from . import settings

def info(message: str, title="INFO", **kwarg):
    print(f"{BLUE(f'[{title}]:')} {message}", **kwarg, flush=True)

def warn(message: str):
    print(f"{YELLOW('[WARN]:')} {message}")

def die(err: str, status_code: int = 1, exit_from_code=True):
    print(f"{RED('[ERROR]:')} {err}")
    if exit_from_code: exit(status_code)

def progress(iterable, label="Progress"):
    total = len(iterable)
    for i, item in enumerate(iterable, 1):
        print(f"{BLUE(f'[{label}]')}: [{i}/{total}]", end="\r", flush=True)
        yield item

def YELLOW(text):
    if settings.NO_COLOR: return text
    return f"\033[33m{text}\033[0m"

def BLUE(text):
    if settings.NO_COLOR: return text
    return f"\033[34m{text}\033[0m"

def GRAY(text):
    if settings.NO_COLOR: return text
    return f"\033[90m{text}\033[0m"

def RED(text):
    if settings.NO_COLOR: return text
    return f"\033[31m{text}\033[0m"

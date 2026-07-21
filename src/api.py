"""Bridge between merodi core and user plugins"""

class GlobalStore:
    def __init__(self):
        self.data = {}

class API:
    def __init__(self):
        self.globals = GlobalStore()
        self.log = None
        self.config = None

api = API()

def expose(name):
    def decorator(fn):
        setattr(api, name, fn)
        return fn
    return decorator

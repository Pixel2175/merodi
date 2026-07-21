"""Bridge between merodi core and user plugins"""

class GlobalStore:
    def __init__(self):
        self._data = {}

    def get(self, key, default=None):
        return self._data.get(key, default)

    def set(self, key, value):
        self._data[key] = value

    def has(self, key):
        return key in self._data

    def delete(self, key):
        self._data.pop(key, None)

    def clear(self):
        self._data.clear()

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

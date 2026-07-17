import hashlib
import json
from os import path

from src.fileops import read_file, write_file
from src.hooks import hook_call
from src.log import warn

HASH_FILE_CONTENT: dict | None = None


def read_hash_file(config):
    global HASH_FILE_CONTENT
    if HASH_FILE_CONTENT is None:
        with open(config.cache.hash, "r") as f:
            HASH_FILE_CONTENT = json.load(f)

    return


def write_hash_file(config):
    global HASH_FILE_CONTENT
    with open(config.cache.hash, "w") as f:
        json.dump(HASH_FILE_CONTENT, f)

def md5_handler(md_path):
    md_content = read_file(md_path)
    return hashlib.md5(md_content.encode()).hexdigest()


def hash_is_modified(md_path, new_hash):
    global HASH_FILE_CONTENT
    old_hash = HASH_FILE_CONTENT.get(md_path)
    if old_hash == new_hash:
        return None
    return new_hash


def append_md5(key, value):
    global HASH_FILE_CONTENT
    HASH_FILE_CONTENT[key] = value


def handle_hash_sync(config, md_path):
    global HASH_FILE_CONTENT
    try:
        hash_dir = config.cache.hash
        if config is None:
            return
        if not (path.exists(hash_dir) and path.isfile(hash_dir)):
            write_file(hash_dir, "{}")

        read_hash_file(config)

        if HASH_FILE_CONTENT is None:
            HASH_FILE_CONTENT = {}

        current_hash = md5_handler(md_path)
        md_rel_path = path.relpath(md_path, config.tree.markdown)
        hook_call("on_hash_check", md_rel_path, current_hash)
        new_hash = hash_is_modified(md_rel_path, current_hash)

        if new_hash is None:
            return None

        append_md5(md_rel_path, new_hash)
        write_hash_file(config)
        hook_call("on_hash_written", md_rel_path, new_hash)
        return True

    except Exception as e:
        warn(f"Hashing failed: {e}")

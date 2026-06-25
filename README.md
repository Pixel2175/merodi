# Merodi

A markdown-based static site generator built with Python.

## Build

- with pip
```
pip(x) install .
```

## Usage

### Init a new project

```
merodi init <path>  # default is the current directory
```

Creates:

```
my-site/
  config.toml
  src/
    md/
      index.md
    templates/
      layout.html
    static/
      style.css
    plugins.py
    dest/
```

### Build

Build the entire project:

```
merodi build <path>  # default is the current directory
```

Build a single file:

```
merodi build --file input.md output.html

### Webview (live preview)

```
merodi webview <path>  # default is the current directory
```

Opens a GUI window with live reload on file changes.

## Configuration

Project settings are defined in `config.toml`:

```toml
[project]
name        = "my-site"
version     = "0.1.0"
description = "Add your description here"

[tree]
markdown  = "src/md"
static    = "src/static"
templates = "src/templates"
dest      = "src/dest"
plugins   = "src/plugins.py"

[webview]
host        = "localhost"
port        = 8866
html_path   = "/"
static_path = "/static"
```

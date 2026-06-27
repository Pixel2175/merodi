# Merodi

A markdown-based static site generator built with Python.
Write pages in markdown, style them with Jinja2 templates, and get a ready-to-publish website.

## Install

```bash
pip install .
# or using `pipx`
pipx install .
```

## Quick start

```bash
merodi init my-site
cd my-site
merodi build
```

Open `src/dest/index.html` or run `merodi webview` for a live preview.

## Commands

### `merodi init [path]`

Scaffolds a new project (defaults to the current directory):

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

### `merodi build [path]`

Builds all markdown files from `src/md/` into HTML files in `src/dest/`.

| Flag | Description |
|---|---|
| `--file INPUT OUTPUT` | Build a single markdown file |
| `--release` / `--debug` | Build mode (currently reserved for future use) |

### `merodi webview [path]`

Opens a native GUI window with live preview. Changes to markdown, templates, static files, or plugins automatically rebuild and reload the page.

## Markdown features

Merodi supports a wide set of markdown extensions:

- **Standard extras** — tables, footnotes, definition lists, fenced code blocks, abbreviations, attr_list, and more
- **Math** — LaTeX math rendered as MathML (`$...$` for inline, `$$...$$` for blocks)
- **Syntax highlighting** — Pygments-based, style configurable in `config.toml`
- **Inline highlights** — `==highlighted text==`
- **Strikethrough** — `~~strikethrough~~`
- **Better emphasis** — smart handling of `*` and `_`
- **Magic links** — URLs auto-link without wrapping syntax
- **Keyboard keys** — `{++Ctrl+Alt+Del++}` renders styled keyboard keys
- **Details/summary** — collapsible `<details>` blocks
- **Tabbed content** — tabbed code blocks and sections
- **Critic markup** — track suggested edits with `{--delete--}` and `{++add++}`
- **Attribute lists** — use `[.class]` or `[:#id]` instead of `{.class}` / `{:#id}`

## Templates

Pages are rendered with **Jinja2**. Template files live in `src/templates/`. Your layout can use Jinja2 blocks and other features.

### Plugin functions in templates

Functions defined in `src/plugins.py` are available as globals in your templates:

```jinja2
{{ fetch("https://api.example.com/data", type="json") }}
{{ read("src/data/content.txt") }}
```

The default `plugins.py` provides `fetch(url, type)` for HTTP requests and `read(file)` for local files. Add your own functions — any public name (no leading underscore) is automatically available.

## Configuration

Project settings live in `config.toml`:

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

[extras]
highlight = "monokai"
```

### `[extras]` options

| Key | Default | Description |
|---|---|---|
| `highlight` | `"monokai"` | Pygments style name. Set to `"noclasses"` to use CSS classes instead of inline styles. |

## Global flags

| Flag | Description |
|---|---|
| `--verbose` | Show detailed error information (also via `VERBOSE=true`) |
| `--no-color` | Disable colored output (also via `NO_COLOR=true`) |

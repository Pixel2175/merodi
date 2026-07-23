# Merodi

A markdown-based static site generator built with Python. Write pages in markdown, style them with Jinja2 templates, and get a ready-to-publish website.

## Install

```bash
pip install .
# or using pipx
pipx install .
```

## Quick start

```bash
merodi init my-site
cd my-site
merodi build
```

Open `src/dest/index.html`, or run `merodi webview` for a live preview.

## Basics

- Pages are written in Markdown under `src/md/`
- Templates use Jinja2 and live in `src/templates/`
- Static assets go in `src/static/`
- Project settings live in `config.toml`
- Custom functions can be added in `src/plugins.py` and used directly in templates

## Status

Merodi is still under active development (pre-1.0.0). Commands, config, and template behavior may change. Full documentation will be published once the project reaches **1.0.0**.

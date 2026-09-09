# Merodi

A Markdown-based static site generator built with Go. Write pages in Markdown, style them with Jinja2 templates, extend it with Lua plugins, and build a ready-to-publish website.

## Installation

```bash
make           # Build Merodi
make install   # Install Merodi
make uninstall # Remove the installed binary
```

## Quick Start

```bash
merodi init my-site
cd my-site
merodi build
```

## Basics

* Project settings are stored in `config.toml`.
* Pages are written in Markdown and stored in the directory specified in `config.toml`.
* See the `examples/` directory for more information.

## Status

Merodi is still under active development and is currently pre-1.0.0. Commands, configuration, and template behavior may change as the project evolves.

Full documentation will be published once Merodi reaches a stable release.

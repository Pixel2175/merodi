-- merodi.compile lets you convert Markdown to HTML without writing files.

merodi.compile.convert("# Hello\nSome **markdown** content")
-- Runs Jinja and Markdown, then returns the HTML.

merodi.compile.mdtohtml("# Hello\nSome **markdown** content")
-- Converts Markdown to HTML without running Jinja.

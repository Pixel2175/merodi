-- Hooks let your plugin run code at different stages of the build.
-- Hook functions don't take arguments or return anything. Use merodi.build.* to

merodi.hook("before_read", function() end)      -- runs before the Markdown file is read
merodi.hook("read_md", function() end)          -- runs after the raw Markdown has been read
merodi.hook("html_content", function() end)     -- runs after the Markdown is converted to HTML
merodi.hook("apply_html_jinja", function() end) -- runs after Jinja templates are applied
merodi.hook("before_write", function() end)     -- runs before the final HTML is written
merodi.hook("after_write", function() end)      -- runs after the file has been written

-- merodi.abort(reason) stops the build for the current file only.
-- It does not stop the rest of the project from being built.
-- Use it only inside a hook. Calling it outside a hook has no effect.
-- The reason is shown as a warning in the build output.

merodi.abort("reason string")

merodi.build.mode()
-- Returns a string: "release" or "draft" depending on how the build was run

merodi.build.content.get()
-- Returns the current content (raw markdown, HTML, or final output,
-- depending on which hook stage you're in).

merodi.build.content.set("new content")
-- Replaces the current content with the string you provide. Takes effect for
-- whatever the next pipeline stage does with it.

merodi.build.md_path.get()
-- Returns the source markdown file's path

merodi.build.html_path.get()
-- Returns the destination HTML file's path

merodi.build.html_path.set("/some/path.html")
-- Overrides the destination path the final HTML will be written to.

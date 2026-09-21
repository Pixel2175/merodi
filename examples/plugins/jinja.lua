-- Variables available in Jinja templates as {{ key }}.

merodi.jinja.set("key", "value")
-- Adds a variable that can be used in templates as {{ key }}.
-- Values can be strings, numbers, booleans, or tables.

merodi.jinja.get("key")
-- Returns the value stored for "key", or nil if it hasn't been set.

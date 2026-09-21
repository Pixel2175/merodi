By default, Merodi converts raw Markdown to raw HTML.

From:

```markdown
# Hello
```

To:

```html
<h1>Hello</h1>
```

But if you use:

```jinja
{% document title="My Website" lang="en" %}

# Hello
```

It will create a complete HTML document:

```html
<!DOCTYPE html>
<html lang="en">
<head>
<title>My Website</title>
</head>
<body>
<h1>Hello</h1>
</body>
</html>
```

### More specifications

You can also add styles, scripts, and metadata:

```jinja
{% document lang="en" title="Pi66" %}
{% style "/static/tailwind.css" %}
{% style "/static/style.css" %}
{% script "/static/js/script.js" %}
{% meta name="description" content="My Personal website" %}

# Hello
```

This will create:

```html
<!DOCTYPE html>
<html lang="en">
<head>
<title>Pi66</title>
<meta name="description" content="My Personal website">
<link rel="stylesheet" href="/static/tailwind.css">
<link rel="stylesheet" href="/static/style.css">
</head>
<body>
<h1>Hello</h1>

<script src="/static/js/script.js"></script>
</body>
</html>
```


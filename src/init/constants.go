package init

const (
	MarkdownContent = `
{% set page_title = "Merodi" %}
{% set metadata_desc = "A markdown-based static site generator." %}

{% block body %}
# Hello, Alice
{% endblock %}`

	HTMLContent = `<!DOCTYPE html>
<html lang="en">
	<head>
		<meta name="description" content="{{ metadata_desc }}">
		<title>{{ page_title }}</title>
	</head>
	<body>
		{% block body %}{% endblock %}
	</body>
</html>`

	PluginsContent = ``
)

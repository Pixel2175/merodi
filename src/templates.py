MARKDOWN_CONTENT = """{% block content %}
# Hello, Alice
{% endblock %}"""
CSS_CONTENT = """body {
	background-color: #1d2021;
	color: #d5c4a1;
}"""
HTML_CONTENT = """<!DOCTYPE html>
<html lang="en">
	<head>
		<title>{% block title %}{% endblock %}</title>
		<link rel="stylesheet" href="{{ style_css }}">
	</head>
	<body>
		{% block content %}{% endblock %}
	</body>
</html>"""
PLUGINS_CONTENT = """from urllib.request import urlopen
from json import loads

def fetch(url, type:str="text"):
    with  urlopen(url) as request:
        if request.status == 200:
            value = request.read().decode("utf-8")
            if type=="json":
                return  loads(value)
            if value.endswith("\\n"):
                return value.removesuffix("\\n")
            return value
        else:
            raise RuntimeError(f"Failed to fetch '{url}': {request.status}")

def read(file:str):
    content = open(file, "r").read()
    if content.endswith("\\n"):
        return content.removesuffix("\\n")
    return content
"""

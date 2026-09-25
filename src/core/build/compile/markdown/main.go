package markdown

import (
	"bytes"
	"merodi/src/core/build/compile/jinja"
	. "merodi/src/core/state"
	"os"

	L "merodi/src/core/lua"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type Markdown struct {
	Extensions    []goldmark.Extender
	ParsersOption []parser.Option
	Build         L.Build
	Page          *jinja.Page
}

func (self *Markdown) Init(state *State) {
	self.Extensions = state.Lua.Extensions
	self.ParsersOption = state.Lua.ParsersOption
}

func (self *Markdown) ReadMD(mdfile string) (string, error) {
	content, err := os.ReadFile(mdfile)
	return string(content), err
}

func (self *Markdown) MdToDoc(mdcontent string) (string, error) {
	var raw_doc bytes.Buffer
	mdBytes := []byte(mdcontent)

	gm := goldmark.New(
		goldmark.WithParserOptions(self.ParsersOption...),
		goldmark.WithExtensions(self.Extensions...),
		goldmark.WithRendererOptions(html.WithUnsafe()),
		goldmark.WithParserOptions(
			parser.WithBlockParsers(util.Prioritized(&jinja.JinjaEscaper{}, 50)),
		),
	)

	doc := gm.Parser().Parse(text.NewReader(mdBytes))
	if err := gm.Renderer().Render(&raw_doc, mdBytes, doc); err != nil {
		return "", err
	}
	body := raw_doc.String()

	if self.Page == nil || !self.Page.Set {
		return body, nil
	}

	return self.wrapHTML(body), nil
}

func (self *Markdown) wrapHTML(body string) string {
	var b bytes.Buffer

	b.WriteString("<!DOCTYPE html>\n")
	b.WriteString(`<html lang="` + self.Page.Lang + `">` + "\n<head>\n")
	b.WriteString("<title>" + self.Page.Title + "</title>\n")

	for _, m := range self.Page.Metas {
		b.WriteString(`<meta name="` + m.Name + `" content="` + m.Content + `">` + "\n")
	}
	for _, s := range self.Page.Styles {
		b.WriteString(`<link rel="stylesheet" href="` + s + `">` + "\n")
	}

	b.WriteString("</head>\n<body>\n")
	b.WriteString(body)
	b.WriteString("\n")

	for _, s := range self.Page.Scripts {
		b.WriteString(`<script src="` + s + `"></script>` + "\n")
	}

	b.WriteString("</body>\n</html>")
	return b.String()
}

package compile

import (
	"bytes"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"merodi/src/config"
	"merodi/src/core/state"
	"merodi/src/utils"
	"os"
	"path/filepath"
	"strings"
)

type Compile struct {
	State *state.State
	Jinja JinjaEngine
}

var (
	handle = utils.Handle
	check  = utils.Check
)

func (self *Compile) ReadMD(mdfile string) error {
	content, err := os.ReadFile(mdfile)
	self.State.Lua.Build.Content = string(content)
	return err
}

func (self *Compile) MdToHtml(md string) error {
	var raw_html bytes.Buffer
	mdBytes := []byte(md)

	gm := goldmark.New(
		goldmark.WithParserOptions(self.State.Lua.ParsersOption...),
		goldmark.WithExtensions(self.State.Lua.Extensions...),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)

	doc := gm.Parser().Parse(text.NewReader(mdBytes))
	doc = self.Jinja.CleanJinja(mdBytes, doc)
	err := gm.Renderer().Render(&raw_html, mdBytes, doc)
	self.State.Lua.Build.Content = raw_html.String()

	return err
}

func (self *Compile) resolveHTMLPath(mdPath string) (string, error) {
	absMarkdown, err := filepath.Abs(self.State.Config.Tree.Markdown)
	if err != nil {
		return "", err
	}
	absMdPath, err := filepath.Abs(mdPath)
	if err != nil {
		return "", err
	}
	mdRelPath, err := filepath.Rel(absMarkdown, absMdPath)
	if err != nil {
		return "", err
	}
	destDir := self.State.Config.Tree.DraftDest
	if self.State.Lua.Build.Mode == config.Release {
		destDir = self.State.Config.Tree.ReleaseDest
	}

	return filepath.Abs(
		filepath.Join(
			destDir,
			strings.TrimSuffix(mdRelPath, ".md")+".html",
		),
	)
}

func (self *Compile) HtmlToFile(htmlPath string) error {
	if err := os.MkdirAll(filepath.Dir(htmlPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(htmlPath, []byte(self.State.Lua.Build.Content), 0644)
}

func (self *Compile) ResolvePaths(mdfile *string, htmlfile *string) (err error) {
	*mdfile = filepath.Join(self.State.ProjectDir, *mdfile)
	*htmlfile, err = self.resolveHTMLPath(*mdfile)
	return err
}

func (self *Compile) Run(file string) (err error) {
	defer handle(&err)
	var mdfile = file
	var htmlfile string
	check(self.ResolvePaths(&mdfile, &htmlfile))
	self.State.Lua.Build.HTMLPath = htmlfile
	self.State.Lua.Build.MDPath = mdfile
	check(self.State.Lua.RunHook("before_read"))
	check(self.ReadMD(mdfile))
	check(self.State.Lua.RunHook("read_md"))
	check(self.MdToHtml(self.State.Lua.Build.Content))
	check(self.State.Lua.RunHook("html_content"))
	check(self.InitJinjaEngine())
	check(self.Jinja.JinjaHandler(&self.State.Lua.Build.Content))
	check(self.State.Lua.RunHook("apply_html_jinja"))
	check(self.State.Lua.RunHook("before_write"))
	check(self.HtmlToFile(htmlfile))
	check(self.State.Lua.RunHook("after_write"))
	self.State.Lua.Build.Content = ""
	self.State.Lua.Build.HTMLPath = ""
	self.State.Lua.Build.MDPath = ""
	return nil
}

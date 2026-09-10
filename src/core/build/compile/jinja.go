package compile

import (
	"strings"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/config"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"
	"github.com/yuin/goldmark/ast"
)

type JinjaEngine struct {
	Loader loaders.Loader
	Config *config.Config
	Data   map[string]any
}

func (self *Compile) InitJinjaEngine() error {
	loader, err := loaders.NewFileSystemLoader(self.State.Config.Tree.Templates)
	if err != nil {
		return err
	}
	self.Jinja.Loader = loader
	self.Jinja.Config = gonja.DefaultConfig
	self.Jinja.Data = self.State.Lua.Jinja
	return nil
}
func (self *JinjaEngine) CleanJinja(source []byte, doc ast.Node) ast.Node {
	child := doc.FirstChild()
	for child != nil {
		next := child.NextSibling()
		if paragraph, ok := child.(*ast.Paragraph); ok {
			lines := paragraph.Lines()
			hasJinja := false
			for i := 0; i < lines.Len(); i++ {
				segment := lines.At(i)
				if strings.HasPrefix(strings.ReplaceAll(string(segment.Value(source)), " ", ""), "{%") {
					hasJinja = true
					break
				}
			}

			if hasJinja {
				for i := 0; i < lines.Len(); i++ {
					segment := lines.At(i)
					raw := ast.NewRawHTML()
					raw.Segments.Append(segment)
					doc.InsertBefore(doc, child, raw)
				}
				doc.RemoveChild(doc, child)
			}
		}
		child = next
	}
	return doc
}

func (self *JinjaEngine) JinjaHandler(html_content *string) error {
	shiftedLoader, err := loaders.NewShiftedLoader("root", strings.NewReader(*html_content), self.Loader)
	if err != nil {
		return err
	}
	tpl, err := exec.NewTemplate("root", self.Config, shiftedLoader, gonja.DefaultEnvironment)
	if err != nil {
		return err
	}
	*html_content, err = tpl.ExecuteToString(exec.NewContext(self.Data))
	return err
}

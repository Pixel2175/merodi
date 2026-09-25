package jinja

import (
	"errors"
	"fmt"
	"merodi/src/utils"
	"strings"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/config"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"
	"github.com/nikolalohinski/gonja/v2/parser"
)

type JinjaEngine struct {
	Loader loaders.Loader
	Config *config.Config
	Data   map[string]any
	Page   *Page
}

func (self *JinjaEngine) Init(templates string) error {
	loader, err := loaders.NewFileSystemLoader(templates)

	if err != nil {
		return err
	}
	self.Loader = loader
	self.Config = gonja.DefaultConfig
	return nil
}
func (self *JinjaEngine) JinjaHandler(html_content *string) error {
	self.Page = &Page{}
	check := utils.Check

	check(self.registerDocument())

	shiftedLoader, err := loaders.NewShiftedLoader(
		"root",
		strings.NewReader(*html_content),
		self.Loader,
	)
	check(err)

	tpl, err := exec.NewTemplate(
		"root",
		self.Config,
		shiftedLoader,
		gonja.DefaultEnvironment,
	)
	if err != nil {
		var syntaxErr *parser.SyntaxError
		if errors.As(err, &syntaxErr) {
			return fmt.Errorf(
				"Jinja syntax error:\n- Line=%d, Column=%d\n- %s",
				syntaxErr.Line,
				syntaxErr.Column,
				syntaxErr.Message,
			)
		}

		return err
	}

	ctx := exec.NewContext(self.Data)
	ctx.Set("document", self.Page)

	*html_content, err = tpl.ExecuteToString(ctx)
	return err
}

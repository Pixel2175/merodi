package jinja

import (
	"fmt"
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/nodes"
	"github.com/nikolalohinski/gonja/v2/parser"
	"github.com/nikolalohinski/gonja/v2/tokens"
)

type Page struct {
	Type, Title, Lang string
	Styles, Scripts   []string
	Metas             []Meta
	Set               bool
}

type Meta struct {
	Name    string
	Content string
}

type Tag struct {
	loc *tokens.Token
	run func(r *exec.Renderer)
}

func (t *Tag) Position() *tokens.Token { return t.loc }
func (t *Tag) String() string          { return "tag" }
func (t *Tag) Execute(r *exec.Renderer, _ *nodes.ControlStructureBlock) error {
	t.run(r)
	return nil
}

func parseArgs(args *parser.Parser, allowed map[string]bool) (map[string]nodes.Expression, error) {
	exprs := map[string]nodes.Expression{}
	for !args.End() {
		key := args.Match(tokens.Name)
		if key == nil {
			return nil, args.Error("expected argument name", nil)
		}
		if !allowed[key.Val] {
			return nil, args.Error(fmt.Sprintf("argument %q not exists", key.Val), key)
		}
		if args.Match(tokens.Assign) == nil {
			return nil, args.Error("expected '='", key)
		}
		e, err := args.ParseExpression()
		if err != nil {
			return nil, err
		}
		exprs[key.Val] = e
	}
	return exprs, nil
}

func (self *JinjaEngine) fields() map[string]*string {
	return map[string]*string{
		"type": &self.Page.Type, "title": &self.Page.Title, "lang": &self.Page.Lang,
	}
}

func (self *JinjaEngine) parseDocument(p *parser.Parser, args *parser.Parser) (nodes.ControlStructure, error) {
	dst := self.fields()
	allowed := map[string]bool{}
	for k := range dst {
		allowed[k] = true
	}
	exprs, err := parseArgs(args, allowed)
	if err != nil {
		return nil, err
	}
	return &Tag{p.Current(), func(r *exec.Renderer) {
		for k, e := range exprs {
			*dst[k] = r.Eval(e).String()
		}
		self.Page.Set = true
	}}, nil
}

func (self *JinjaEngine) parseMeta(p *parser.Parser, args *parser.Parser) (nodes.ControlStructure, error) {
	exprs, err := parseArgs(args, map[string]bool{"name": true, "content": true})
	if err != nil {
		return nil, err
	}
	return &Tag{p.Current(), func(r *exec.Renderer) {
		m := Meta{}
		if e, ok := exprs["name"]; ok {
			m.Name = r.Eval(e).String()
		}
		if e, ok := exprs["content"]; ok {
			m.Content = r.Eval(e).String()
		}
		self.Page.Metas = append(self.Page.Metas, m)
	}}, nil
}

func parseList(list func() *[]string) parser.ControlStructureParser {
	return func(p *parser.Parser, args *parser.Parser) (nodes.ControlStructure, error) {
		e, err := args.ParseExpression()
		if err != nil {
			return nil, err
		}
		return &Tag{p.Current(), func(r *exec.Renderer) {
			l := list()
			*l = append(*l, r.Eval(e).String())
		}}, nil
	}
}

func (self *JinjaEngine) registerDocument() error {
	cs := gonja.DefaultEnvironment.ControlStructures
	tags := map[string]parser.ControlStructureParser{
		"document": self.parseDocument,
		"meta":     self.parseMeta,
		"style":    parseList(func() *[]string { return &self.Page.Styles }),
		"script":   parseList(func() *[]string { return &self.Page.Scripts }),
	}
	for name, p := range tags {
		if cs.Exists(name) {
			if err := cs.Replace(name, p); err != nil {
				return err
			}
			continue
		}
		if err := cs.Register(name, p); err != nil {
			return err
		}
	}
	return nil
}

package jinja

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type JinjaEscaper struct{}

func (p *JinjaEscaper) Trigger() []byte { return []byte{'{'} }

func (p *JinjaEscaper) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, segment := reader.PeekLine()
	if !bytes.HasPrefix(bytes.TrimLeft(line, " \t"), []byte("{%")) {
		return nil, parser.NoChildren
	}
	node := ast.NewHTMLBlock(ast.HTMLBlockType7)
	reader.Advance(segment.Len() - 1)
	if bytes.Contains(line, []byte("%}")) {
		node.ClosureLine = segment
	} else {
		node.Lines().Append(segment)
	}
	return node, parser.NoChildren
}

func (p *JinjaEscaper) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	n := node.(*ast.HTMLBlock)
	if !n.ClosureLine.IsEmpty() {
		return parser.Close
	}
	line, segment := reader.PeekLine()
	reader.Advance(segment.Len() - 1)
	if bytes.Contains(line, []byte("%}")) {
		n.ClosureLine = segment
	} else {
		n.Lines().Append(segment)
	}
	return parser.Continue | parser.NoChildren
}

func (p *JinjaEscaper) Close(node ast.Node, reader text.Reader, pc parser.Context) {}
func (p *JinjaEscaper) CanInterruptParagraph() bool                                { return true }
func (p *JinjaEscaper) CanAcceptIndentedLine() bool                                { return false }

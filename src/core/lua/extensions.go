package lua

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	glua "github.com/yuin/gopher-lua"
)

var extensionsMap = map[string]goldmark.Extender{
	"Table":          extension.Table,
	"Strikethrough":  extension.Strikethrough,
	"Linkify":        extension.Linkify,
	"TaskList":       extension.TaskList,
	"GFM":            extension.GFM,
	"DefinitionList": extension.DefinitionList,
	"Footnote":       extension.Footnote,
	"Typographer":    extension.Typographer,
	"CJK":            extension.CJK,
}

func (self *Lua) registerExtensions() {
	extensionTable := self.Context.NewTable()
	meta := self.Context.NewTable()
	self.Context.SetField(meta, "__index", self.Context.NewFunction(
		func(L *glua.LState) int {
			name := L.CheckString(2)
			ext, ok := extensionsMap[name]
			if !ok {
				L.RaiseError("unknown extension: %s", name)
				return 0
			}
			self.Extensions = append(self.Extensions, ext)
			return 0
		},
	))
	extensionTable.Metatable = meta
	self.Merodi.RawSetString("extensions", extensionTable)
}

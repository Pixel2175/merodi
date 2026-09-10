package lua

import (
	"github.com/yuin/goldmark/renderer/html"
	glua "github.com/yuin/gopher-lua"
)

var disableMap = map[string]func(){
	"GlobalAttribute": func() {
		html.GlobalAttributeFilter = nil
	},
	"HeadingAttribute": func() {
		html.HeadingAttributeFilter = nil
	},
	"BlockquoteAttribute": func() {
		html.BlockquoteAttributeFilter = nil
	},
	"CodeAttribute": func() {
		html.CodeAttributeFilter = nil
	},
	"EmphasisAttribute": func() {
		html.EmphasisAttributeFilter = nil
	},
	"ImageAttribute": func() {
		html.ImageAttributeFilter = nil
	},
	"LinkAttribute": func() {
		html.LinkAttributeFilter = nil
	},
	"ListAttribute": func() {
		html.ListAttributeFilter = nil
	},
	"ListItemAttribute": func() {
		html.ListItemAttributeFilter = nil
	},
	"ParagraphAttribute": func() {
		html.ParagraphAttributeFilter = nil
	},
	"ThematicAttribute": func() {
		html.ThematicAttributeFilter = nil
	},
}

func (self *Lua) registerDisable() {
	disableTable := self.Context.NewTable()
	meta := self.Context.NewTable()

	self.Context.SetField(meta, "__index", self.Context.NewFunction(
		func(L *glua.LState) int {
			name := L.CheckString(2)

			disable, ok := disableMap[name]
			if !ok {
				L.RaiseError("unknown filter: %s", name)
				return 0
			}

			fn := self.Context.NewFunction(func(L *glua.LState) int {
				disable()
				return 0
			})

			L.Push(fn)
			return 1
		},
	))

	disableTable.Metatable = meta
	self.Merodi.RawSetString("disable", disableTable)
}

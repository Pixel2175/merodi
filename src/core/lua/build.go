package lua

import (
	glua "github.com/yuin/gopher-lua"
	"merodi/src/config"
)

type Build struct {
	MDPath   string
	DocPath string
	Content  string
	Mode     config.Mode
}

func (self *Lua) registerBuild() {
	dataTable := self.Context.NewTable()
	dataTable.RawSetString("mode", self.Context.NewFunction(func(L *glua.LState) int {
		L.Push(glua.LString(config.ModeToString(self.Build.Mode)))
		return 1
	}))

	dataTable.RawSetString("content", self.makeGetSetTable(
		func() string { return self.Build.Content },
		func(v string) { self.Build.Content = v },
	))
	dataTable.RawSetString("md_path", self.makeGetSetTable(
		func() string { return self.Build.MDPath },
		func(v string) { self.Build.MDPath = v },
	))
	dataTable.RawSetString("document_path", self.makeGetSetTable(
		func() string { return self.Build.DocPath },
		func(v string) { self.Build.DocPath = v },
	))

	self.Merodi.RawSetString("build", dataTable)
}

func (self *Lua) makeGetSetTable(get func() string, set func(string)) *glua.LTable {
	t := self.Context.NewTable()
	self.Context.SetFuncs(t, map[string]glua.LGFunction{
		"get": func(L *glua.LState) int {
			L.Push(glua.LString(get()))
			return 1
		},
		"set": func(L *glua.LState) int {
			set(L.ToString(1))
			return 0
		},
	})
	return t
}

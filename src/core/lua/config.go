package lua

import (
	lua "github.com/yuin/gopher-lua"
)

func (self *Lua) registerConfigTree() *lua.LTable {
	tree := self.Context.NewTable()
	tree.RawSetString("markdown", lua.LString(self.Config.Tree.Markdown))
	tree.RawSetString("templates", lua.LString(self.Config.Tree.Templates))
	tree.RawSetString("draft_dest", lua.LString(self.Config.Tree.DraftDest))
	tree.RawSetString("release_dest", lua.LString(self.Config.Tree.ReleaseDest))
	return tree
}

func (self *Lua) registerConfigProject() *lua.LTable {
	projectTable := self.Context.NewTable()
	projectTable.RawSetString("name", lua.LString(self.Config.Project.Name))
	projectTable.RawSetString("version", lua.LString(self.Config.Project.Version))
	projectTable.RawSetString("description", lua.LString(self.Config.Project.Description))
	return projectTable
}

func (self *Lua) registerConfig() {
	configTable := self.Context.NewTable()

	configTable.RawSetString("tree", self.registerConfigTree())
	configTable.RawSetString("project", self.registerConfigProject())

	self.Merodi.RawSetString("config", configTable)
}

package lua

import (
	lua "github.com/yuin/gopher-lua"
)

func (self *Lua) registerJinja() {
	jinja := self.Context.NewTable()
	self.Context.SetFuncs(jinja, map[string]lua.LGFunction{
		"set": self.luaJinjaSet,
		"get": self.luaJinjaGet,
	})
	self.Merodi.RawSetString("jinja", jinja)
}

func (self *Lua) luaJinjaSet(L *lua.LState) int {
	key := L.ToString(1)
	value := gluaToGo(L.Get(2))
	self.Jinja[key] = value
	return 0
}

func (self *Lua) luaJinjaGet(L *lua.LState) int {
	key := L.ToString(1)
	L.Push(goToGlua(self.Jinja[key]))
	return 1
}

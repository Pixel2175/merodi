package lua

import glua "github.com/yuin/gopher-lua"

func (self *Lua) registerAction() {
	self.Merodi.RawSetString("action", self.Context.NewFunction(func(L *glua.LState) int {
		L.Push(glua.LString(string(self.State.Action)))
		return 1
	}))
}

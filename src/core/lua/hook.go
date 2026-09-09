package lua

import (
	"merodi/src/utils"
	"merodi/src/utils/log"

	lua "github.com/yuin/gopher-lua"
)

func (self *Lua) registerHook(L *lua.LState) int {
	stage := L.ToString(1)
	fn := L.ToFunction(2)
	self.Hooks[stage] = fn
	return 0
}

func (self *Lua) luaAbort(L *lua.LState) int {
	log.Warn(L.ToString(1))
	self.Aborted = true
	panic(utils.ErrAborted)
}

func (self *Lua) registerHooks() {
	self.Context.SetFuncs(self.Merodi, map[string]lua.LGFunction{
		"hook":  self.registerHook,
		"abort": self.luaAbort,
	})
}

func (self *Lua) RunHook(stage string) (err error) {
	fn, ok := self.Hooks[stage]
	if !ok {
		return nil
	}

	self.Aborted = false
	err = self.Context.CallByParam(lua.P{
		Fn:      fn,
		NRet:    0,
		Protect: true,
	})
	if self.Aborted {
		return utils.ErrAborted
	}

	return err
}

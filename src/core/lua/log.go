package lua

import (
	"merodi/src/utils/log"

	lua "github.com/yuin/gopher-lua"
)

func (self *Lua) luaLogDie(L *lua.LState) int {
	msg := L.CheckString(1)
	log.Die(msg)
	return 0
}

func (self *Lua) luaLogWarn(L *lua.LState) int {
	msg := L.CheckString(1)
	log.Warn(msg)
	return 0
}

func (self *Lua) luaLogInfo(L *lua.LState) int {
	msg := L.CheckString(1)
	log.Info("%s", msg)
	return 0
}

func (self *Lua) registerLog() {
	logTable := self.Context.NewTable()
	self.Context.SetFuncs(logTable, map[string]lua.LGFunction{
		"info": self.luaLogInfo,
		"warn": self.luaLogWarn,
		"die":  self.luaLogDie,
	})

	self.Merodi.RawSetString("log", logTable)
}

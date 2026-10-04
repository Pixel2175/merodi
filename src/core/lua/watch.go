package lua

import (
	glua "github.com/yuin/gopher-lua"
)

func (self *Lua) registerWatch() {
	watchTable := self.Context.NewTable()
	watchTable.RawSetString("add", self.Context.NewFunction(func(l *glua.LState) int {
		self.State.Watch.Add = append(self.State.Watch.Add, l.ToString(1))
		return 0
	}))
	self.Merodi.RawSetString("watch", watchTable)
}

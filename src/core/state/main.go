package state

import (
	. "merodi/src/config"
	"merodi/src/core/lua"
)

type State struct {
	ProjectDir string
	Config     Data
	Lua        *lua.Lua
}

package state

import (
	"merodi/src/config"
	. "merodi/src/config"
	"merodi/src/core/lua"
)

type State struct {
	ProjectDir string
	Config     Data
	BuildMode  config.Mode
	Lua        *lua.Lua
}

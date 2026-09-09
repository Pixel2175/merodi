package core

import (
	"merodi/src/core/build"
	"merodi/src/core/lua"
	"merodi/src/core/serve"
	"merodi/src/core/state"
	"merodi/src/utils"
)

var (
	handle = utils.Handle
	check  = utils.Check
)

type Core state.State
type Action interface {
	Run(state *state.State, args *[]string) error
}

var Actions = map[string]Action{
	"build": &build.Build{},
	"serve": &serve.Serve{},
}

func Run(act string, args *[]string) (err error) {
	defer handle(&err)
	self := Core{Lua: &lua.Lua{}}
	action, ok := Actions[act]
	if !ok {
		utils.PrintHelp()
		return
	}
	check(self.GoToProjectDir(args))
	check(self.LoadConfig())

	check(self.Lua.InitLua(&self.Config))
	defer self.Lua.Context.Close()

	check(action.Run((*state.State)(&self), args))
	return nil
}

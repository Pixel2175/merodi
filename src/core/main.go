package core

import (
	"merodi/src/core/build"
	"merodi/src/core/lua"
	"merodi/src/core/state"
	"merodi/src/core/watch"
	"merodi/src/utils"
	. "merodi/src/utils/errs"
)

type Core state.State
type Action interface {
	Run(state *state.State, args *[]string) error
}

var Actions = map[string]Action{
	"build": &build.Build{},
	"watch": &watch.Watcher{},
}

func Run(act string, args *[]string) (err error) {
	defer Handle(&err)
	self := Core{Lua: &lua.Lua{}}
	action, ok := Actions[act]
	if !ok {
		utils.PrintHelp()
		return
	}
	self.BuildMode = build.SetBuildMode(args)
	CheckE(self.GoToProjectDir(args))
	CheckE(self.LoadConfig())

	CheckE(self.Lua.InitLua(&self.Config))
	defer self.Lua.Context.Close()

	CheckE(action.Run((*state.State)(&self), args))
	return nil
}

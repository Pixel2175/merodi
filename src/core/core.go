package core

import (
	"merodi/src/core/build"
	"merodi/src/core/lua"
	. "merodi/src/core/state"
	"merodi/src/core/watch"
	"merodi/src/utils"
	. "merodi/src/utils/errs"
)

type Action interface {
	Run(state *State, args *[]string) error
}

var Actions = map[string]Action{
	"build": &build.Build{},
	"watch": &watch.Watcher{},
}

func Run(act string, args *[]string) (err error) {
	defer Handle(&err)
	state := InitGlobals()

	action, ok := Actions[act]
	if !ok {
		utils.PrintHelp()
		return
	}

	state.Build.Mode = build.SetBuildMode(args)
	CheckE(state.GoToProjectDir(args))
	CheckE(state.LoadConfig())
	CheckE(lua.Init(&state))
	CheckE(action.Run(&state, args))

	return nil
}

package core

import (
	"merodi/src/core/build"
	"merodi/src/core/clean"
	"merodi/src/core/hook"
	"merodi/src/core/lua"
	. "merodi/src/core/state"
	"merodi/src/core/watch"
	"merodi/src/parser"
	. "merodi/src/utils/errs"
)

type Action interface {
	Run(state *State, opt *parser.Option) error
}

var Actions = map[parser.Action]Action{
	parser.Build: &build.Build{},
	parser.Watch: &watch.Watcher{},
	parser.Clean: &clean.Clean{},
	parser.Hook:  &hook.Hook{},
}

func Run(opt *parser.Option) (err error) {
	defer Handle(&err)
	state := InitGlobals()

	action := Actions[opt.Action]

	state.Build.Mode = opt.Mode
	CheckE(state.GoToProjectDir(opt.ProjectDir))
	CheckE(state.LoadConfig())
	if opt.Action != parser.Clean {
		CheckE(lua.Init(&state))
		defer state.Lua.Context.Close()
	}
	CheckE(action.Run(&state, opt))

	return nil
}

package build

import (
	"io/fs"
	"merodi/src/config"
	"merodi/src/core/build/compile"
	. "merodi/src/core/state"
	"merodi/src/utils"
	"path/filepath"
)

type Build struct {
	state   *State
	mode    config.Mode
	compile *compile.Compile
}

var (
	handle = utils.Handle
	check  = utils.Check
)

func (build *Build) SetBuildMode(args *[]string) {
	if utils.PopString(args, "--release") != "" {
		build.mode = config.Release
		return
	}

	build.mode = config.Draft
}

func (build *Build) WalkAndBuild() (err error) {
	return filepath.WalkDir(build.state.Config.Tree.Markdown, func(
		path string,
		entry fs.DirEntry,
		err error,
	) error {
		if err != nil {
			return err
		}

		if !entry.Type().IsRegular() && entry.Type()&fs.ModeSymlink == 0 {
			return nil
		}

		err = build.compile.Run(path)
		if err == utils.ErrAborted {
			return nil
		}

		return err
	})
}

func (build *Build) Run(state *State, args *[]string) (err error) {
	defer handle(&err)

	build.SetBuildMode(args)

	build.compile = &compile.Compile{}
	build.state = state
	build.compile.State = build.state
	build.state.Lua.Build.Mode = build.mode

	check(build.WalkAndBuild())

	return nil
}

package build

import (
	"io/fs"
	"merodi/src/core/build/compile"
	"merodi/src/core/state"
	. "merodi/src/utils/errs"
)

type Build struct {
	state   *state.State
	compile *compile.Compile
}

func (self *Build) InitCompile() {
	self.compile = &compile.Compile{State: self.state}
	CheckE(self.compile.Init())
}

func (self *Build) BuildFile(path string) (err error) {
	defer Handle(&err)
	html := CheckV(self.compile.Run(path))
	_, docPath, rerr := self.compile.ResolveDocPaths(path)
	CheckE(rerr)
	CheckE(self.writeDoc(docPath, html))
	return
}

func (self *Build) Visit(path string, entry fs.DirEntry, err error) (rerr error) {
	defer Handle(&rerr, func(err *error) {
		if *err == ErrAborted {
			*err = self.state.RunHook("on_abort")
		}
	})

	CheckE(err)
	if !self.isSource(entry) {
		return
	}
	CheckE(self.BuildFile(path))
	return
}

func (self *Build) Run(state *state.State, args *[]string) (err error) {
	defer Handle(&err)

	self.state = state
	self.InitCompile()

	CheckE(self.walkAndBuild())
	return
}

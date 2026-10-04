package build

import (
	"io/fs"
	"merodi/src/config"
	"merodi/src/core/build/compile"
	. "merodi/src/core/state"
	. "merodi/src/utils/errs"
)

type Build struct {
	state   *State
	mode    config.Mode
	compile *compile.Compile
}

func (self *Build) InitCompile() {
	self.state.Lua.Build.Mode = self.mode
	self.compile = &compile.Compile{State: self.state}
	CheckE(self.compile.Init(self.mode))
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
	defer Handle(&rerr)
	CheckE(err)
	if !self.isSource(entry) {
		return
	}
	CheckE(self.BuildFile(path))
	return
}

func (self *Build) Run(state *State, args *[]string) (err error) {
	defer Handle(&err)

	self.state = state
	self.mode = state.BuildMode
	self.InitCompile()

	CheckE(self.walkAndBuild())
	return
}

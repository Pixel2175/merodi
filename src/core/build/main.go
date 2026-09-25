package build

import (
	"io/fs"
	"merodi/src/config"
	"merodi/src/core/build/compile"
	. "merodi/src/core/state"
	"merodi/src/utils"
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

func (self *Build) InitCompile() {
	self.state.Lua.Build.Mode = self.mode
	self.compile = &compile.Compile{State: self.state}
	check(self.compile.Init(self.mode))
}

func (self *Build) BuildFile(path string) error {
	html, err := self.compile.Run(path)
	check(err)

	_, docPath, err := self.compile.ResolveDocPaths(path)
	check(err)

	return self.writeDoc(docPath, html)
}

func (self *Build) Visit(path string, entry fs.DirEntry, err error) error {
	check(err)
	if !self.isSource(entry) {
		return nil
	}
	check(self.BuildFile(path))
	return nil
}

func (self *Build) Run(state *State, args *[]string) (err error) {
	defer handle(&err)

	self.state = state
	self.setBuildMode(args)

	self.InitCompile()

	return self.walkAndBuild()
}

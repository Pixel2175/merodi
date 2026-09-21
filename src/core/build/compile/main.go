package compile

import (
	"merodi/src/config"
	. "merodi/src/core/build/compile/jinja"
	. "merodi/src/core/build/compile/markdown"
	L "merodi/src/core/lua"
	"merodi/src/core/state"
	"merodi/src/utils"
)

type Compile struct {
	State *state.State
	Lua   *L.Lua
	Jinja JinjaEngine
	MD    Markdown
}

var (
	handle = utils.Handle
	check  = utils.Check
)

func (self *Compile) Init(mode config.Mode) (err error) {
	defer handle(&err)

	self.Lua = self.State.Lua

	self.Jinja = JinjaEngine{}
	check(self.Jinja.Init(self.State.Config.Tree.Templates))

	self.MD = Markdown{}
	self.MD.Init(self.State)

	return nil
}

func (self *Compile) Run(mdfile string) (doc string, err error) {
	defer handle(&err)
	defer self.reset()

	self.resolvePaths(mdfile)
	self.readMarkdown()
	self.applyJinja()
	self.convertToDoc()
	self.hook("before_write")

	return self.Lua.Build.Content, nil
}


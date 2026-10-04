package compile

import (
	. "merodi/src/core/build/compile/jinja"
	. "merodi/src/core/build/compile/markdown"
	"merodi/src/core/state"
	. "merodi/src/utils/errs"
)

type Compile struct {
	State  *state.State
	Jinja  JinjaEngine
	MD     Markdown
	Build  state.Build
}

func (self *Compile) Init() (err error) {
	defer Handle(&err)

	self.Jinja = JinjaEngine{Title: self.State.Config.Project.Description}
	CheckE(self.Jinja.Init(self.State.Config.Tree.Templates))

	self.MD = Markdown{}
	self.MD.Init(self.State)

	return nil
}

func (self *Compile) Run(mdfile string) (doc string, err error) {
	defer Handle(&err)
	defer self.reset()

	self.resolvePaths(mdfile)
	self.readMarkdown()
	self.applyJinja()
	self.convertToDoc()
	self.hook("before_write")

	return self.State.Build.Content, nil
}

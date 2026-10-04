package compile

import (
	"merodi/src/config"
	. "merodi/src/utils/errs"
)

func (self *Compile) hook(name string) {
	CheckE(self.State.RunHook(name))
}

func (self *Compile) resolvePaths(mdfile string) {
	build := self.State.Build
	mdPath, docPath, err := self.ResolveDocPaths(mdfile)
	CheckE(err)
	build.MDPath, build.DocPath = mdPath, docPath
}

func (self *Compile) readMarkdown() {
	build := self.State.Build
	self.hook("before_read")
	build.Content = CheckV(self.MD.ReadMD(build.MDPath))
	self.hook("read_md")
}

func (self *Compile) convertToDoc() {
	build := self.State.Build
	self.MD.Page = *self.Jinja.Page
	build.Content = CheckV(self.MD.MdToDoc(build.Content))
	self.hook("doc_content")
}

func (self *Compile) applyJinja() {
	build := self.State.Build
	self.Jinja.Data = self.State.Jinja
	CheckE(self.Jinja.JinjaHandler(&build.Content))
	self.hook("apply_doc_jinja")
}

func (self *Compile) reset() {
	build := self.State.Build
	build.Content = ""
	build.DocPath = ""
	build.MDPath = ""
	build.Mode = config.None
}

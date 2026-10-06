package compile

import (
	. "merodi/src/utils/errs"
)

func (self *Compile) hook(name string) {
	CheckE(self.State.RunHook(name))
}

func (self *Compile) resolvePaths(mdfile string) {
	mdPath, docPath, err := self.ResolveDocPaths(mdfile)
	CheckE(err)
	self.Build.MDPath, self.Build.DocPath = mdPath, docPath
}

func (self *Compile) readMarkdown() {
	self.hook("before_read")
	self.Build.Content = CheckV(self.MD.ReadMD(self.Build.MDPath))
	self.hook("read_md")
}

func (self *Compile) convertToDoc() {
	self.MD.Page = *self.Jinja.Page
	self.Build.Content = CheckV(self.MD.MdToDoc(self.Build.Content))
	self.hook("doc_content")
}

func (self *Compile) applyJinja() {
	self.Jinja.Data = self.State.Jinja
	CheckE(self.Jinja.JinjaHandler(&self.Build.Content))
	self.hook("apply_doc_jinja")
}

func (self *Compile) reset() {
	self.Build.Content = ""
	self.Build.DocPath = ""
	self.Build.MDPath = ""
}

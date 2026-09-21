package compile

import "merodi/src/config"

func (self *Compile) hook(name string) {
	check(self.Lua.RunHook(name))
}

func (self *Compile) resolvePaths(mdfile string) {
	var err error
	build := self.Lua.Build
	build.MDPath, build.DocPath, err = self.ResolveDocPaths(mdfile)
	check(err)
}

func (self *Compile) readMarkdown() {
	var err error
	build := self.Lua.Build
	self.hook("before_read")
	build.Content, err = self.MD.ReadMD(build.MDPath)
	check(err)
	self.hook("read_md")
}

func (self *Compile) convertToDoc() {
	var err error
	build := self.Lua.Build
	self.MD.Page = self.Jinja.Page
	build.Content, err = self.MD.MdToDoc(build.Content)
	check(err)
	self.hook("doc_content")
}

func (self *Compile) applyJinja() {
	build := self.Lua.Build
	self.Jinja.Data = self.State.Lua.Jinja
	check(self.Jinja.JinjaHandler(&build.Content))
	self.hook("apply_doc_jinja")
}

func (self *Compile) reset() {
	build := self.Lua.Build
	build.Content = ""
	build.DocPath = ""
	build.MDPath = ""
	build.Mode = config.None
}

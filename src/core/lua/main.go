package lua

import (
	"merodi/src/config"
	"path/filepath"

	"github.com/yuin/goldmark"
	glua "github.com/yuin/gopher-lua"
)

type Lua struct {
	Context *glua.LState
	Config  *config.Data

	Hooks      map[string]*glua.LFunction
	Jinja      map[string]any
	Build      *Build
	Http       *Http
	Extensions []goldmark.Extender

	Merodi  *glua.LTable
	Aborted bool
}

func (self *Lua) InitLua(cfg *config.Data) error {
	self.Config = cfg
	self.Hooks = make(map[string]*glua.LFunction)
	self.Jinja = make(map[string]any)
	self.Context = glua.NewState()
	self.Http = &Http{}
	self.Build = &Build{}
	self.Aborted = false

	self.registerMerodi()
	return self.Context.DoFile(filepath.Join(cfg.Tree.Plugins, "main.lua"))
}

func (self *Lua) registerMerodi() {
	self.Merodi = self.Context.NewTable()
	self.registerHooks()
	self.registerConfig()
	self.registerJinja()
	self.registerBuild()
	self.registerHttp()
	self.registerExtensions()
	self.registerLog()
	self.Context.SetGlobal("merodi", self.Merodi)
}

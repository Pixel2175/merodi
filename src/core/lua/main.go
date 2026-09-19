package lua

import (
	"merodi/src/config"
	"path/filepath"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	glua "github.com/yuin/gopher-lua"
)

type Lua struct {
	Context *glua.LState
	Config  *config.Data

	Hooks map[string]*glua.LFunction
	Jinja map[string]any
	Build *Build
	Watch *Watcher

	Extensions    []goldmark.Extender
	ParsersOption []parser.Option

	Merodi  *glua.LTable
	Aborted bool
}

func (self *Lua) InitLua(cfg *config.Data) error {
	self.Config = cfg
	self.Hooks = make(map[string]*glua.LFunction)
	self.Jinja = make(map[string]any)
	self.Context = glua.NewState()
	self.Watch = &Watcher{}
	self.Build = &Build{}
	self.Aborted = false

	self.registerMerodi()

	path := filepath.Join(cfg.Tree.Plugins, "main.lua")
	if err := self.loadFile(path); err != nil {
		return err
	}
	return nil
}

func (self *Lua) loadFile(path string) error {
	fn, err := self.Context.LoadFile(path)
	if err != nil {
		return describeLuaError(path, err)
	}

	self.Context.Push(fn)
	if err := self.Context.PCall(0, glua.MultRet, nil); err != nil {
		return describeLuaError(path, err)
	}

	return nil
}
func (self *Lua) registerMerodi() {
	self.Merodi = self.Context.NewTable()
	self.registerHooks()
	self.registerConfig()
	self.registerJinja()
	self.registerBuild()
	self.registerWatch()
	self.registerEnable()
	self.registerDisable()
	self.registerLog()
	self.Context.SetGlobal("merodi", self.Merodi)
}

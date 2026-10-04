package lua

import (
	"merodi/src/config"
	. "merodi/src/utils/errs"
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

func (self *Lua) InitLua(cfg *config.Data) (err error) {
	defer Handle(&err)
	self.Config = cfg
	self.Hooks = make(map[string]*glua.LFunction)
	self.Jinja = make(map[string]any)
	self.Context = glua.NewState()
	self.Watch = &Watcher{}
	self.Build = &Build{}
	self.Aborted = false

	self.registerMerodi()

	path := filepath.Join(cfg.Tree.Plugins, "main.lua")
	CheckE(self.loadFile(path))
	return
}

func (self *Lua) loadFile(path string) (err error) {
	defer Handle(&err)
	fn, lerr := self.Context.LoadFile(path)
	if lerr != nil {
		CheckE(describeLuaError(path, lerr))
	}

	self.Context.Push(fn)
	if perr := self.Context.PCall(0, glua.MultRet, nil); perr != nil {
		CheckE(describeLuaError(path, perr))
	}
	return
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

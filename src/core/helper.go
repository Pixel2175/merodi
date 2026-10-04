package core

import (
	"merodi/src/config"
	"merodi/src/project"
	. "merodi/src/utils/errs"
	"os"
)

func (self *Core) GoToProjectDir(args *[]string) (err error) {
	defer Handle(&err)
	dir := CheckV(project.SetProjectDir(args))
	self.ProjectDir = dir
	CheckE(os.Chdir(dir))
	return
}

func (self *Core) LoadConfig() (err error) {
	defer Handle(&err)
	var cfg config.Config
	cfg.ProjectDir = self.ProjectDir
	CheckE(cfg.Read())
	self.Config = cfg.Data
	return
}

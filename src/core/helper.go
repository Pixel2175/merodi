package core

import (
	"merodi/src/config"
	"merodi/src/project"
	"os"
)

func (self *Core) GoToProjectDir(args *[]string) error {
	dir, err := project.SetProjectDir(args)

	if err != nil {
		return err
	}
	self.ProjectDir = dir
	return os.Chdir(dir)
}

func (self *Core) LoadConfig() error {
	var cfg config.Config
	cfg.ProjectDir = self.ProjectDir
	if err := cfg.Read(); err != nil {
		return err
	}
	self.Config = cfg.Data
	return nil
}

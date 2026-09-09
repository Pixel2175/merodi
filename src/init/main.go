package init

import (
	"fmt"
	"merodi/src/config"
	"merodi/src/project"
	"merodi/src/utils"
	"merodi/src/utils/log"
)

type Init struct {
	ProjectDir string
	cfg        config.Tree
}

var (
	handle = utils.Handle
	check  = utils.Check
)

func (init *Init) GetProjectDir(args *[]string) error {
	var err error
	init.ProjectDir, err = project.SetProjectDir(args)
	if project.IsProjectExists(init.ProjectDir) {
		err = fmt.Errorf("project already initialized: `config.toml` already exists")
	}
	return err
}

func (init *Init) InitConfig() error {
	cfg := config.Config{}
	cfg.ProjectDir = init.ProjectDir
	cfg.Init()
	init.cfg = cfg.Data.Tree
	return cfg.Write()
}

func (init Init) Run(args *[]string) (err error) {
	defer handle(&err)
	check(init.GetProjectDir(args))
	check(init.InitConfig())
	check(init.Bootstrap())
	log.Info("Project initialized.")
	return nil
}

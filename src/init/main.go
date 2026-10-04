package init

import (
	"errors"
	"merodi/src/config"
	"merodi/src/project"
	. "merodi/src/utils/errs"
	"merodi/src/utils/log"
)

type Init struct {
	ProjectDir string
	cfg        config.Tree
}

var ErrProjectExists = errors.New("project already initialized: `config.toml` already exists")

func (init *Init) GetProjectDir(args *[]string) (err error) {
	defer Handle(&err)
	init.ProjectDir = CheckV(project.SetProjectDir(args))
	if project.IsProjectExists(init.ProjectDir) {
		CheckE(ErrProjectExists)
	}
	return
}

func (init *Init) InitConfig() error {
	cfg := config.Config{}
	cfg.ProjectDir = init.ProjectDir
	cfg.Init()
	init.cfg = cfg.Data.Tree
	return cfg.Write()
}

func (init Init) Run(args *[]string) (err error) {
	defer Handle(&err)
	CheckE(init.GetProjectDir(args))
	CheckE(init.InitConfig())
	CheckE(init.Bootstrap())
	log.Info("Project initialized.")
	return nil
}

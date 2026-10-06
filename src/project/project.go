package project

import (
	"merodi/src/utils"
	. "merodi/src/utils/errs"
	"merodi/src/utils/path"
	"os"
	"path/filepath"
)

func SetProjectDir(args *[]string) (string, error) {
	projectPath := path.Getwd()
	if len(*args) > 0 {
		projectPath = utils.Pop(args, 0)
	}
	projectPath = CheckV(filepath.Abs(projectPath))

	for projectPath != "/" {
		if IsProjectExists(projectPath) {
			return projectPath, nil
		}

		projectPath = filepath.Dir(projectPath)
	}

	return "", os.ErrNotExist
}

func IsProjectExists(path string) bool {
	configPath := filepath.Join(path, "config.toml")

	_, err := os.Stat(configPath)
	return err == nil
}

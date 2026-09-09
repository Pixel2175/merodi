package project

import (
	"merodi/src/utils"
	"merodi/src/utils/path"
	"os"
	"path/filepath"
)

func SetProjectDir(args *[]string) (string, error) {
	projectPath := path.Getwd()
	if len(*args) > 0 {
		projectPath = utils.Pop(args, 0)
	}
	return filepath.Abs(projectPath)
}

func IsProjectExists(path string) bool {
	configPath := filepath.Join(path, "config.toml")

	_, err := os.Stat(configPath)
	return err == nil
}

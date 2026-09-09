package path

import (
	"os"
	"path/filepath"
)

func Getwd() string {
	dir, err := os.Getwd()
	if err != nil {
		return filepath.Clean(".")
	}

	return dir
}

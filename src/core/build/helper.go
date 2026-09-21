package build

import (
	"io/fs"
	"merodi/src/config"
	"merodi/src/utils"
	"os"
	"path/filepath"
)

func (build *Build) isSource(entry fs.DirEntry) bool {
	return entry.Type().IsRegular() || entry.Type()&fs.ModeSymlink != 0
}

func (build *Build) writeDoc(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

func (build *Build) setBuildMode(args *[]string) {
	build.mode = config.Draft
	if utils.PopString(args, "--release") != "" {
		build.mode = config.Release
	}
}

func (build *Build) walkAndBuild() error {
	return filepath.WalkDir(build.state.Config.Tree.Markdown, build.Visit)
}

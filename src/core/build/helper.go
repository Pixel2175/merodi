package build

import (
	"io/fs"
	"merodi/src/config"
	"merodi/src/utils"
	. "merodi/src/utils/errs"
	"os"
	"path/filepath"
)

func (build *Build) isSource(entry fs.DirEntry) bool {
	return entry.Type().IsRegular() || entry.Type()&fs.ModeSymlink != 0
}

func (build *Build) writeDoc(path, content string) (err error) {
	defer Handle(&err)
	CheckE(os.MkdirAll(filepath.Dir(path), 0o755))
	CheckE(os.WriteFile(path, []byte(content), 0o644))
	return
}

func (build *Build) walkAndBuild() (err error) {
	defer Handle(&err)
	CheckE(filepath.WalkDir(build.state.Config.Tree.Markdown, build.Visit))
	return
}

func SetBuildMode(args *[]string) config.Mode {
	mode := config.Draft
	if utils.PopString(args, "--release") != "" {
		mode = config.Release
	}
	return mode
}

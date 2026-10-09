package build

import (
	"io/fs"
	. "merodi/src/utils/errs"
	"path/filepath"
)

func (build *Build) isSource(entry fs.DirEntry) bool {
	return entry.Type().IsRegular() || entry.Type()&fs.ModeSymlink != 0
}

func (build *Build) walkAndBuild() (err error) {
	defer Handle(&err)
	CheckE(filepath.WalkDir(build.state.Config.Tree.Markdown, build.Visit))
	return
}

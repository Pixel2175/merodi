package compile

import (
	. "merodi/src/utils/errs"
	"os"
	"path/filepath"
	"strings"

	"merodi/src/config"
)

func (self *Compile) destRoot() string {
	if self.State.BuildMode == config.Release {
		return self.State.Config.Tree.ReleaseDest
	}
	return self.State.Config.Tree.DraftDest
}

func (self *Compile) docDestPath(mdPath string) (dest string, err error) {
	defer Handle(&err)
	mdRoot := CheckV(filepath.Abs(self.State.Config.Tree.Markdown))
	absMd := CheckV(filepath.Abs(mdPath))
	rel := CheckV(filepath.Rel(mdRoot, absMd))

	name := strings.TrimSuffix(rel, ".md") + "." + "html"
	dest = CheckV(filepath.Abs(filepath.Join(self.destRoot(), name)))
	return
}

func (self *Compile) WriteDoc(docContent string, docPath string) (err error) {
	defer Handle(&err)
	CheckE(os.MkdirAll(filepath.Dir(docPath), 0o755))
	CheckE(os.WriteFile(docPath, []byte(docContent), 0o644))
	return
}

func (self *Compile) ResolveDocPaths(mdFile string) (mdPath, docPath string, err error) {
	defer Handle(&err)
	mdPath = filepath.Join(self.State.ProjectDir, mdFile)
	docPath = CheckV(self.docDestPath(mdPath))
	return
}

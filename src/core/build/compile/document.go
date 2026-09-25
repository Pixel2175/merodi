package compile

import (
	"os"
	"path/filepath"
	"strings"

	"merodi/src/config"
)

func (self *Compile) destRoot() string {
	if self.Lua.Build.Mode == config.Release {
		return self.State.Config.Tree.ReleaseDest
	}
	return self.State.Config.Tree.DraftDest
}

func (self *Compile) docDestPath(mdPath string) (string, error) {
	mdRoot, err := filepath.Abs(self.State.Config.Tree.Markdown)
	if err != nil {
		return "", err
	}
	absMd, err := filepath.Abs(mdPath)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(mdRoot, absMd)
	if err != nil {
		return "", err
	}
	name := strings.TrimSuffix(rel, ".md") + "." + "html"
	return filepath.Abs(filepath.Join(self.destRoot(), name))
}

func (self *Compile) WriteDoc(docContent string, docPath string) error {
	if err := os.MkdirAll(filepath.Dir(docPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(docPath, []byte(docContent), 0o644)
}

func (self *Compile) ResolveDocPaths(mdFile string) (mdPath, docPath string, err error) {
	mdPath = filepath.Join(self.State.ProjectDir, mdFile)
	docPath, err = self.docDestPath(mdPath)
	return mdPath, docPath, err
}

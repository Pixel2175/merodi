package init

import (
	"merodi/src/config"
	. "merodi/src/utils/errs"
	"os"
	"path/filepath"
)

func (init *Init) CreateDirectories(tree config.Tree) (err error) {
	defer Handle(&err)
	for _, dir := range []string{
		tree.Markdown,
		tree.Templates,
		tree.Plugins,
	} {
		path := filepath.Join(init.ProjectDir, dir)
		CheckE(os.MkdirAll(path, 0755))
	}
	return
}

func (init *Init) CreateFiles(tree config.Tree) (err error) {
	defer Handle(&err)
	files := map[string]string{
		filepath.Join(init.ProjectDir, tree.Markdown, "index.md"):     MarkdownContent,
		filepath.Join(init.ProjectDir, tree.Templates, "layout.html"): HTMLContent,
		filepath.Join(init.ProjectDir, tree.Plugins, "main.lua"):      PluginsContent,
	}

	for path, content := range files {
		CheckE(os.WriteFile(path, []byte(content), 0644))
	}
	return
}

func (self *Init) Bootstrap() (err error) {
	defer Handle(&err)
	CheckE(self.CreateDirectories(self.cfg))
	CheckE(self.CreateFiles(self.cfg))
	return
}

package init

import (
	"merodi/src/config"
	"os"
	"path/filepath"
)

func (init *Init) CreateDirectories(tree config.Tree) error {
	for _, dir := range []string{
		tree.Markdown,
		tree.Templates,
		tree.Plugins,
	} {
		path := filepath.Join(init.ProjectDir, dir)

		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}

	return nil
}

func (init *Init) CreateFiles(tree config.Tree) error {
	files := map[string]string{
		filepath.Join(init.ProjectDir, tree.Markdown, "index.md"):     MarkdownContent,
		filepath.Join(init.ProjectDir, tree.Templates, "layout.html"): HTMLContent,
		filepath.Join(init.ProjectDir, tree.Plugins, "main.lua"):      PluginsContent,
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

func (self *Init) Bootstrap() error {
	err := self.CreateDirectories(self.cfg)
	if err != nil {
		return err
	}

	err = self.CreateFiles(self.cfg)
	return err
}

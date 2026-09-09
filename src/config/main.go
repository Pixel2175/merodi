package config

import (
	"merodi/src/utils"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Mode int

const ConfigFileName = "config.toml"
const (
	None Mode = iota
	Release
	Draft
)

type Project struct {
	Name        string `toml:"name"`
	Version     string `toml:"version"`
	Description string `toml:"description"`
}

type Tree struct {
	Markdown    string `toml:"markdown"`
	Templates   string `toml:"templates"`
	Plugins     string `toml:"plugins"`
	DraftDest   string `toml:"draft_dest"`
	ReleaseDest string `toml:"release_dest"`
}

type Http struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
}

type Data struct {
	Project Project `toml:"project"`
	Tree    Tree    `toml:"tree"`
	Http    Http    `toml:"http"`
}

type Config struct {
	Data       Data
	ProjectDir string
}

func (self *Config) DefaultData() {
	self.Data = Data{
		Project: Project{
			Name:        "potato",
			Version:     "0.1.0",
			Description: "Add your description here",
		},
		Tree: Tree{
			Markdown:    "src/md",
			Templates:   "src/templates",
			Plugins:     "src/plugins",
			DraftDest:   "build/draft",
			ReleaseDest: "build/release",
		},
		Http: Http{
			Host: "localhost",
			Port: "8000",
		},
	}
}

func ModeToString(mode Mode) string {
	switch mode {
	case Release:
		return "release"
	case Draft:
		return "Draft"
	default:
		return "None"
	}
}

func (self *Config) ConfigPath() string {
	return filepath.Join(self.ProjectDir, ConfigFileName)
}

func (self *Config) Init() error {

	if err := os.MkdirAll(self.ProjectDir, 0755); err != nil {
		return err
	}
	f, err := os.Create(self.ConfigPath())
	if err != nil {
		return err
	}
	defer f.Close()
	self.DefaultData()
	return toml.NewEncoder(f).Encode(self.Data)
}

func (self *Config) Read() error {
	_, err := toml.DecodeFile(self.ConfigPath(), &self.Data)
	return err
}

func (self *Config) Write() (err error) {
	defer utils.Handle(&err)
	f, err := os.Create(self.ConfigPath())
	if err != nil {
		return err
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(self.Data)
}

package main

import (
	"merodi/src/core"
	. "merodi/src/init"
	"merodi/src/parser"
	. "merodi/src/utils/errs"
	"merodi/src/utils/log"
	"os"
)

func main() {
	var err error
	defer Handle(&err, func(err *error) {
		log.Die((*err).Error())
	})

	opt := parser.Parse(os.Args[1:])

	switch opt.Action {
	case parser.Help:
		parser.PrintHelp()
	case parser.Init:
		CheckE((&Init{ProjectDir: opt.ProjectDir}).Run())
	case parser.Version:
		log.Info(log.Title("Version"), "Merodi %s", parser.VERSION)
	default:
		CheckE(core.Run(&opt))
	}
}

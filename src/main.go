package main

import (
	"merodi/src/core"
	. "merodi/src/init"
	"merodi/src/utils"
	"merodi/src/utils/log"
	"os"
)

const VERSION = "0.5.1"

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		utils.PrintHelp()
		return
	}
	pos := utils.Pop(&args, 0)

	var err error
	switch pos {
	case "init", "new":
		init := Init{}
		err = init.Run(&args)

	case "version":
		log.Info(log.Title("Version"), "Merodi %s", VERSION)

	default:
		err = core.Run(pos, &args)
	}

	if err != nil {
		log.Die(err.Error())
	}
}

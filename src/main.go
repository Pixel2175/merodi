package main

import (
	"merodi/src/core"
	. "merodi/src/init"
	"merodi/src/utils"
	. "merodi/src/utils/errs"
	"merodi/src/utils/log"
	"os"
)

const VERSION = "0.5.2"

func main() {
	args := os.Args[1:]
	var err error 
	defer Handle(&err, func(err *error) {
		log.Die((*err).Error())
	})

	if len(args) == 0 {
		utils.PrintHelp()
		return
	}
	cmd := utils.Pop(&args, 0)

	switch cmd {
	case "init", "new":
		CheckE((&Init{}).Run(&args))
	case "version":
		log.Info(log.Title("Version"), "Merodi %s", VERSION)
	default:
		CheckE(core.Run(cmd, &args))
	}
	return
}

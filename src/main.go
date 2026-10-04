package main

import (
	"merodi/src/core"
	. "merodi/src/init"
	"merodi/src/utils"
	. "merodi/src/utils/errs"
	"merodi/src/utils/log"
	"os"
)

const VERSION = "0.5.1"

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Die(err.Error())
	}
}

func run(args []string) (err error) {
	defer Handle(&err)

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

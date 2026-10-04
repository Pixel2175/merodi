package utils

import (
	"merodi/src/utils/log"
)

func PrintHelp() {
	log.Info(log.Title("Usage"), ": merodi {init,watch,build,version} <path> ...")
}

func Pop(args *[]string, idx int) string {
	value := (*args)[idx]
	*args = append((*args)[:idx], (*args)[idx+1:]...)
	return value
}

func PopString(args *[]string, value string) string {
	for i, arg := range *args {
		if arg == value {
			*args = append((*args)[:i], (*args)[i+1:]...)
			return value
		}
	}

	return ""
}

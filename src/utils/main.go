package utils

import (
	"errors"
	"merodi/src/utils/log"
)

var ErrAborted = errors.New("build aborted")

func Handle(err *error) {
	if r := recover(); r != nil {
		if e, ok := r.(error); ok {
			*err = e
		} else {
			panic(r)
		}
	}
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func PrintHelp() {
	log.Info(log.Title("Usage"), ": merodi {init,watch,build,serve,version} <path> ...")
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

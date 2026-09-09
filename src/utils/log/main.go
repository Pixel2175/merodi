package log

import (
	"fmt"
	"io"
	"os"
)

var file io.Writer = os.Stdout

type Title string

func Info(args ...any) {
	title := "INFO"

	if t, ok := args[0].(Title); ok {
		title = string(t)
		args = args[1:]
	}

	msg := fmt.Sprintf(args[0].(string), args[1:]...)
	fmt.Fprintf(file, "[%s]: %s\n", blue(title), msg)
}

func Warn(args ...any) {
	title := "WARN"

	if t, ok := args[0].(Title); ok {
		title = string(t)
		args = args[1:]
	}

	msg := fmt.Sprintf(args[0].(string), args[1:]...)
	fmt.Fprintf(file, "[%s]: %s\n", yellow(title), msg)
}

func Die(args ...any) {
	title := "ERROR"

	if t, ok := args[0].(Title); ok {
		title = string(t)
		args = args[1:]
	}

	msg := fmt.Sprintf(args[0].(string), args[1:]...)
	fmt.Fprintf(file, "[%s]: %s\n", red(title), msg)

	os.Exit(1)
}

func blue(msg string) string {
	return fmt.Sprintf("\033[34m%s\033[0m", msg)
}

func yellow(msg string) string {
	return fmt.Sprintf("\033[33m%s\033[0m", msg)
}

func red(msg string) string {
	return fmt.Sprintf("\033[31m%s\033[0m", msg)
}

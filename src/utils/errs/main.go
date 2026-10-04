package errs

import (
	"errors"
)

var ErrAborted = errors.New("build aborted")

func Handle(err *error, fn ...func(error)) {
	if r := recover(); r != nil {
		e, ok := r.(error)
		if !ok {
			panic(r)
		}
		*err = e
		if len(fn) > 0 && fn[0] != nil {
			fn[0](e)
		}
	}
}

func CheckE(err error) {
	if err != nil {
		panic(err)
	}
}

func CheckV[T any](v T, err error) T {
	CheckE(err)
	return v
}

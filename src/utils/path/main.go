package path

import (
	. "merodi/src/utils/errs"
	"os"
	"path/filepath"
)

func Getwd() (dir string) {
	var err error
	defer Handle(&err, func(*error) { dir = filepath.Clean(".") })
	dir = CheckV(os.Getwd())
	return
}

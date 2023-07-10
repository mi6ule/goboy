package util

import (
	"runtime"
	"strconv"
)

func GetInvokedPath(path string) string {
	if path != "" {
		return path
	}
	_, file, line, _ := runtime.Caller(2)
	return file + ":" + strconv.Itoa(line)
}

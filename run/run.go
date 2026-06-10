// Package run
package run

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/kovey/debug-go/debug"
)

func RunCo(f func()) {
	defer func() {
		Panic(recover())
	}()
	f()
}

func Panic(err any) bool {
	if err == nil {
		return false
	}

	logs := []string{fmt.Sprintf("panic error[%s]", err)}
	for i := 3; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		logs = append(logs, fmt.Sprintf("%s(%d)", file, line))
	}

	debug.Erro(strings.Join(logs, "\r\n"))
	return true
}

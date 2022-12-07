package run

import (
	"runtime"

	"github.com/kovey/debug-go/debug"
)

func RunCo(f func()) {
	defer func() {
		Panic(recover())
	}()
	f()
}

func Panic(err interface{}) bool {
	if err == nil {
		return false
	}

	debug.Erro("panic error[%s]", err)

	for i := 3; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		debug.Erro("%s(%d)", file, line)
	}

	return true
}

package debug

import (
	"fmt"
	"runtime"
	"time"

	"github.com/kovey/debug-go/color"
	"github.com/kovey/debug-go/util"
)

type DebugType string

type DebugValue int32

type DebugLevels map[DebugType]DebugValue

var level DebugValue = val_info

func (d DebugLevels) CanShow(t DebugType) bool {
	if l, ok := d[t]; ok {
		return l >= level
	}

	return false
}

func (d DebugLevels) Get(t DebugType) DebugValue {
	if val, ok := d[t]; ok {
		return val
	}

	return val_none
}

const (
	val_info DebugValue = 1
	val_dbug DebugValue = 2
	val_warn DebugValue = 3
	val_erro DebugValue = 4
	val_none DebugValue = 5
)

const (
	Debug_Info DebugType = "info"
	Debug_Dbug DebugType = "dbug"
	Debug_Warn DebugType = "warn"
	Debug_Erro DebugType = "erro"
	Debug_None DebugType = "none"
)

const (
	echoFormat         = "[%s][%s] %s\n"
	echoFormatFileLine = "[%s][%s] %s in file %s on line %d\n"
)

const (
	caller = 2
)

var maps = DebugLevels{
	Debug_Info: val_info,
	Debug_Dbug: val_dbug,
	Debug_Warn: val_warn,
	Debug_Erro: val_erro,
}

func SetLevel(t DebugType) {
	level = maps.Get(t)
}

func echo(format string, t DebugType, args ...interface{}) {
	if !maps.CanShow(t) {
		return
	}

	switch t {
	case Debug_Warn:
		if fileLineSwitch == File_Line_Show {
			_, file, line, _ := runtime.Caller(caller)
			str := fmt.Sprintf(echoFormatFileLine, time.Now().Format(util.Golang_Birthday_Time), t, fmt.Sprintf(format, args...), file, line)
			fmt.Print(color.Yellow(str))
			return
		}
		str := fmt.Sprintf(echoFormat, time.Now().Format(util.Golang_Birthday_Time), t, fmt.Sprintf(format, args...))
		fmt.Print(color.Yellow(str))
	case Debug_Erro:
		if fileLineSwitch == File_Line_Show {
			_, file, line, _ := runtime.Caller(caller)
			str := fmt.Sprintf(echoFormatFileLine, time.Now().Format(util.Golang_Birthday_Time), t, fmt.Sprintf(format, args...), file, line)
			fmt.Print(color.Red(str))
			return
		}
		str := fmt.Sprintf(echoFormat, time.Now().Format(util.Golang_Birthday_Time), t, fmt.Sprintf(format, args...))
		fmt.Print(color.Red(str))
	case Debug_Dbug:
		if fileLineSwitch == File_Line_Show {
			_, file, line, _ := runtime.Caller(caller)
			str := fmt.Sprintf(echoFormatFileLine, time.Now().Format(util.Golang_Birthday_Time), t, fmt.Sprintf(format, args...), file, line)
			fmt.Print(color.Magenta(str))
			return
		}
		str := fmt.Sprintf(echoFormat, time.Now().Format(util.Golang_Birthday_Time), t, fmt.Sprintf(format, args...))
		fmt.Print(color.Magenta(str))
	default:
		if fileLineSwitch == File_Line_Show {
			_, file, line, _ := runtime.Caller(caller)
			fmt.Printf(echoFormatFileLine, time.Now().Format(util.Golang_Birthday_Time), t, fmt.Sprintf(format, args...), file, line)
			return
		}
		fmt.Printf(echoFormat, time.Now().Format(util.Golang_Birthday_Time), t, fmt.Sprintf(format, args...))
	}
}

func Info(format string, args ...interface{}) {
	echo(format, Debug_Info, args...)
}

func Erro(format string, args ...interface{}) {
	echo(format, Debug_Erro, args...)
}

func Warn(format string, args ...interface{}) {
	echo(format, Debug_Warn, args...)
}

func Dbug(format string, args ...interface{}) {
	echo(format, Debug_Dbug, args...)
}

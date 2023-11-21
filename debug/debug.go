package debug

import (
	"fmt"
	"runtime"

	"github.com/kovey/debug-go/color"
	"github.com/kovey/debug-go/util/now"
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
	val_test DebugValue = 5
	val_none DebugValue = 6
)

const (
	Debug_Info DebugType = "info"
	Debug_Dbug DebugType = "dbug"
	Debug_Warn DebugType = "warn"
	Debug_Erro DebugType = "erro"
	Debug_Test DebugType = "test"
	Debug_None DebugType = "none"
)

const (
	echoFormat         = "[%s][%s] %s\n"
	echoFormatFileLine = "[%s][%s] %s(%d): %s\n"
)

const (
	caller = 3
)

var maps = DebugLevels{
	Debug_Info: val_info,
	Debug_Dbug: val_dbug,
	Debug_Warn: val_warn,
	Debug_Erro: val_erro,
	Debug_Test: val_test,
}

func SetLevel(t DebugType) {
	level = maps.Get(t)
}

func echo(format string, t DebugType, args ...any) {
	if !maps.CanShow(t) {
		return
	}

	switch t {
	case Debug_Warn:
		_echo(t, format, args, color.Yellow)
	case Debug_Erro:
		_echo(t, format, args, color.Red)
	case Debug_Dbug:
		_echo(t, format, args, color.Magenta)
	case Debug_Test:
		_echo(t, format, args, color.Green)
	default:
		_echo(t, format, args, _info)
	}
}

func _info(content string) string {
	return content
}

func _echo(t DebugType, format string, args []any, color func(string) string) {
	if fileLineSwitch == File_Line_Show {
		_, file, line, _ := runtime.Caller(caller)
		str := fmt.Sprintf(echoFormatFileLine, now.DateTime(), t, file, line, fmt.Sprintf(format, args...))
		fmt.Print(color(str))
		return
	}
	str := fmt.Sprintf(echoFormat, now.DateTime(), t, fmt.Sprintf(format, args...))
	fmt.Print(color(str))
}

func Info(format string, args ...any) {
	echo(format, Debug_Info, args...)
}

func Erro(format string, args ...any) {
	echo(format, Debug_Erro, args...)
}

func Warn(format string, args ...any) {
	echo(format, Debug_Warn, args...)
}

func Dbug(format string, args ...any) {
	echo(format, Debug_Dbug, args...)
}

func Test(format string, args ...any) {
	echo(format, Debug_Test, args...)
}

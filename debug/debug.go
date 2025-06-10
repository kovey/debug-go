package debug

import (
	"io"
	"os"
)

type DebugType string

type DebugValue int32

type DebugLevels map[DebugType]DebugValue

type FormatType byte

type _json struct {
	Level   DebugType `json:"level"`
	LogMsg  string    `json:"log_msg"`
	TraceId string    `json:"trace_id"`
	SpanId  string    `json:"span_id"`
	File    string    `json:"file"`
	Line    string    `json:"line"`
	Time    string    `json:"time"`
}

const (
	Format_Normal FormatType = 1
	Format_Json   FormatType = 2
)

var level DebugValue = val_info
var writer io.Writer = os.Stdout
var formatType FormatType = Format_Normal
var log = &Log{caller: 4}

func SetWriter(w io.Writer) {
	writer = w
}

func UseJsonFormat() {
	formatType = Format_Json
}

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
	echoFormat              = "[%s][%s] %s\r\n"
	echoFormatFileLine      = "[%s][%s] %s(%d): %s\r\n"
	echoFormatTrace         = "[%s][%s][%s][%s] %s\r\n"
	echoFormatFileLineTrace = "[%s][%s][%s][%s] %s(%d): %s\r\n"
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

func Info(format string, args ...any) {
	log.Info(format, args...)
}

func Erro(format string, args ...any) {
	log.Erro(format, args...)
}

func Warn(format string, args ...any) {
	log.Warn(format, args...)
}

func Dbug(format string, args ...any) {
	log.Dbug(format, args...)
}

func Test(format string, args ...any) {
	log.Test(format, args...)
}

func Json(data any) {
	log.Json(data)
}

func LogWith(traceId, spandId string) *Log {
	return &Log{traceId: traceId, spanId: spandId, caller: caller}
}

package debug

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/kovey/debug-go/color"
	"github.com/kovey/debug-go/util/now"
)

type Log struct {
	traceId string
	spanId  string
	caller  int
}

func (l *Log) echo(format string, level DebugType, args ...any) {
	if !maps.CanShow(level) {
		return
	}

	switch level {
	case Debug_Warn:
		l._echo(format, level, args, color.Yellow)
	case Debug_Erro:
		l._echo(format, level, args, color.Red)
	case Debug_Dbug:
		l._echo(format, level, args, color.Magenta)
	case Debug_Test:
		l._echo(format, level, args, color.Green)
	default:
		l._echo(format, level, args, l._info)
	}
}

func (l *Log) _info(content string) string {
	return content
}

func (l *Log) _echo(format string, level DebugType, args []any, color func(string) string) {
	if fileLineSwitch == File_Line_Show {
		_, file, line, _ := runtime.Caller(l.caller)
		switch formatType {
		case Format_Normal:
			if l.traceId == "" {
				str := fmt.Sprintf(echoFormatFileLine, now.DateTime(), level, file, line, fmt.Sprintf(format, args...))
				writer.Write([]byte(color(str)))
				return
			}
			str := fmt.Sprintf(echoFormatFileLineTrace, now.DateTime(), level, l.traceId, l.spanId, file, line, fmt.Sprintf(format, args...))
			writer.Write([]byte(color(str)))
		case Format_Json:
			Json(_json{Level: level, LogMsg: fmt.Sprintf(format, args...), File: file, Line: fmt.Sprintf("%d", line), Time: now.DateTime(), TraceId: l.traceId, SpanId: l.spanId})
		}
		return
	}
	switch formatType {
	case Format_Normal:
		if l.traceId == "" {
			str := fmt.Sprintf(echoFormat, now.DateTime(), level, fmt.Sprintf(format, args...))
			writer.Write([]byte(color(str)))
			return
		}
		str := fmt.Sprintf(echoFormatTrace, now.DateTime(), level, l.traceId, l.spanId, fmt.Sprintf(format, args...))
		writer.Write([]byte(color(str)))
	case Format_Json:
		Json(_json{Level: level, LogMsg: fmt.Sprintf(format, args...), Time: now.DateTime(), TraceId: l.traceId, SpanId: l.spanId})
	}
}

func (l *Log) WithTraceId(traceId string) *Log {
	l.traceId = traceId
	return l
}

func (l *Log) WithSpanId(spanId string) *Log {
	l.spanId = spanId
	return l
}

func (l *Log) Info(format string, args ...any) {
	l.echo(format, Debug_Info, args...)
}

func (l *Log) Erro(format string, args ...any) {
	l.echo(format, Debug_Erro, args...)
}

func (l *Log) Warn(format string, args ...any) {
	l.echo(format, Debug_Warn, args...)
}

func (l *Log) Dbug(format string, args ...any) {
	l.echo(format, Debug_Dbug, args...)
}

func (l *Log) Test(format string, args ...any) {
	l.echo(format, Debug_Test, args...)
}

func (l *Log) Json(data any) {
	if buff, err := json.Marshal(data); err == nil {
		writer.Write(buff)
		writer.Write([]byte("\r\n"))
	}
}

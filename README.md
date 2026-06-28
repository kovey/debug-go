# debug-go

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.20-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/license-Apache%202.0-green)](./LICENSE)

A lightweight, feature-rich debugging and logging library for Go, designed for development and production environments. It provides color-coded console output, structured JSON logging, distributed tracing support, asynchronous file output with daily rotation, and panic recovery utilities — all with a minimal, zero-dependency API.

## Features

- **Five log levels** — `info`, `dbug`, `warn`, `erro`, `test` with configurable filtering
- **Color-coded terminal output** — each level has a distinct ANSI color for instant visual scanning
- **Dual output formats** — human-readable text or structured JSON for log aggregation
- **Source location** — optionally include file name and line number in every log entry
- **Distributed tracing** — attach `trace_id` and `span_id` to correlate logs across services
- **Asynchronous file logging** — non-blocking writes with configurable buffer size and automatic daily rotation
- **Panic recovery** — `run.RunCo` captures panics and logs full stack traces
- **Custom writers** — redirect output to any `io.Writer` (files, sockets, buffers)
- **Zero external dependencies** — relies only on the Go standard library

## Installation

```bash
go get -u github.com/kovey/debug-go
```

Requires Go 1.20 or later.

## Quick Start

```go
package main

import "github.com/kovey/debug-go/debug"

func main() {
    // Set the minimum log level (messages below this level are suppressed)
    debug.SetLevel(debug.Debug_Info)

    debug.Info("server started on port %d", 8080)
    debug.Dbug("processing request from %s", clientIP)
    debug.Warn("retry attempt %d/%d", attempt, maxRetries)
    debug.Erro("failed to connect to database: %s", err)
}
```

**Output:**

```
[2024-01-15 10:30:45][info] server started on port 8080
[2024-01-15 10:30:45][dbug] processing request from 192.168.1.1
[2024-01-15 10:30:46][warn] retry attempt 2/3
[2024-01-15 10:30:47][erro] failed to connect to database: connection refused
```

## Log Levels

| Level      | Constant            | Value | Typical Use                              | Color    |
|------------|---------------------|-------|------------------------------------------|----------|
| Info       | `debug.Debug_Info`  | 1     | General operational messages             | Default  |
| Debug      | `debug.Debug_Dbug`  | 2     | Diagnostic details during development    | Magenta  |
| Warning    | `debug.Debug_Warn`  | 3     | Potentially harmful situations           | Yellow   |
| Error      | `debug.Debug_Erro`  | 4     | Error events that may still allow continued operation | Red      |
| Test       | `debug.Debug_Test`  | 5     | Test-specific output                     | Green    |

Set the threshold with `SetLevel` — only messages at or above that level are emitted:

```go
debug.SetLevel(debug.Debug_Warn) // Only warn, erro, and test messages appear
```

## Output Formats

### Text Format (default)

Human-readable, color-coded output with timestamps:

```go
debug.Info("user %s logged in", username)
// [2024-01-15 10:30:45][info] user alice logged in
```

### JSON Format

Structured output suitable for log aggregators (ELK, Loki, Datadog, etc.):

```go
debug.UseJsonFormat()

debug.Info("user %s logged in", username)
// {"level":"info","log_msg":"user alice logged in","trace_id":"","span_id":"","file":"","line":"","time":"2024-01-15 10:30:45"}
```

Check the current format at runtime:

```go
if debug.FormatIsJson() {
    // JSON mode
}
```

## Source File and Line Numbers

Enable file and line number annotations to pinpoint exactly where each log line originated:

```go
debug.SetFileLine(debug.File_Line_Show)

debug.Erro("unexpected value: %v", value)
// [2024-01-15 10:30:47][erro] /app/handler.go(42): unexpected value: nil

debug.SetFileLine(debug.File_Line_Off) // Turn off (default)
```

## Distributed Tracing

Attach trace and span identifiers for request-level correlation across distributed systems:

```go
// Create a contextual logger with trace identifiers
logger := debug.LogWith("trace-abc123", "span-xyz789")

logger.Info("querying database")
// [2024-01-15 10:30:45][info][trace-abc123][span-xyz789] querying database

logger.Erro("query timeout after %s", elapsed)
// [2024-01-15 10:30:47][erro][trace-abc123][span-xyz789] query timeout after 30s
```

You can also chain `WithTraceId` and `WithSpanId` on an existing logger:

```go
logger := debug.LogWith("trace-abc123", "")
logger.WithSpanId("span-xyz789").Info("message")
```

## JSON Output for Arbitrary Data

Log structured data directly as JSON:

```go
debug.Json(map[string]any{
    "event":   "payment_received",
    "amount":  99.95,
    "user_id": 1234,
})
// {"event":"payment_received","amount":99.95,"user_id":1234}
```

## Custom Writer

Redirect log output to any destination that implements `io.Writer`:

```go
f, _ := os.Create("/var/log/app.log")
debug.SetWriter(f)
debug.Info("logging to file")
```

## Asynchronous File Logging

For production workloads, use the `async` package for non-blocking file output with automatic daily rotation:

```go
package main

import (
    "context"
    "os/signal"
    "syscall"

    "github.com/kovey/debug-go/async"
    "github.com/kovey/debug-go/debug"
)

func main() {
    // Start async logging: writes to /var/log/app/ with a 1024-entry buffer
    if err := async.Start("/var/log/app", 1024); err != nil {
        panic(err)
    }
    defer async.Close()

    // All debug output is now written asynchronously to daily log files
    debug.Info("application started")

    // Start the listener goroutine that flushes buffered logs to disk
    ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer cancel()
    go async.Listen(ctx)

    // ... application logic ...

    // On shutdown, Close() drains the buffer and closes the current file
}
```

Log files are named by date and automatically roll over at midnight:

```
/var/log/app/2024-01-15.log
/var/log/app/2024-01-16.log
/var/log/app/2024-01-17.log
```

## Panic Recovery

Safely execute goroutines with automatic panic capture and stack trace logging:

```go
go run.RunCo(func() {
    // If this function panics, the stack trace is logged via debug.Erro
    // and the panic is recovered — the process does not crash.
    riskyOperation()
})
```

You can also use `run.Panic` directly with `recover()`:

```go
defer func() {
    if run.Panic(recover()) {
        // Panic was captured and logged
    }
}()
```

## API Reference

### `debug` package

| Function | Description |
|----------|-------------|
| `SetLevel(t DebugType)` | Set the minimum log level threshold |
| `SetWriter(w io.Writer)` | Redirect output to a custom writer |
| `SetFileLine(fl FileLine)` | Enable/disable file and line number annotations |
| `UseJsonFormat()` | Switch to JSON output format |
| `FormatIsJson() bool` | Check if JSON format is active |
| `Info(format, args...)` | Log at info level |
| `Dbug(format, args...)` | Log at debug level |
| `Warn(format, args...)` | Log at warning level |
| `Erro(format, args...)` | Log at error level |
| `Test(format, args...)` | Log at test level |
| `Json(data any)` | Marshal and write arbitrary data as JSON |
| `LogWith(traceId, spanId string) *Log` | Create a contextual logger with trace identifiers |

### `*Log` methods

| Method | Description |
|--------|-------------|
| `Info(format, args...)` | Log at info level with trace context |
| `Dbug(format, args...)` | Log at debug level with trace context |
| `Warn(format, args...)` | Log at warning level with trace context |
| `Erro(format, args...)` | Log at error level with trace context |
| `Test(format, args...)` | Log at test level with trace context |
| `Json(data any)` | Marshal and write data as JSON |
| `WithTraceId(traceId string) *Log` | Set the trace ID (chainable) |
| `WithSpanId(spanId string) *Log` | Set the span ID (chainable) |

### `async` package

| Function | Description |
|----------|-------------|
| `Start(logDir string, length int) error` | Initialize async file logging with buffer size |
| `Listen(ctx context.Context)` | Start the flush loop (run in a goroutine) |
| `Close()` | Drain buffer, close file, and stop logging |

### `run` package

| Function | Description |
|----------|-------------|
| `RunCo(f func())` | Execute `f` and recover from panics, logging the stack trace |
| `Panic(err any) bool` | Log and recover from a panic value; returns true if a panic occurred |

## License

Apache 2.0 — see [LICENSE](./LICENSE) for details.

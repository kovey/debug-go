package async

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/kovey/debug-go/debug"
)

type logFile struct {
	logDir string
	data   chan []byte
	date   string
	ticker *time.Ticker
	file   *file
	closed atomic.Bool
}

func newLogFile(logDir string, length int) *logFile {
	return &logFile{logDir: logDir, data: make(chan []byte, length), date: time.Now().Format(time.DateOnly), ticker: time.NewTicker(1 * time.Second), file: &file{}}
}

func (l *logFile) Write(p []byte) (n int, err error) {
	if l.closed.Load() {
		return 0, fmt.Errorf("log is closed")
	}

	select {
	case l.data <- p:
		return len(p), nil
	default:
		return 0, fmt.Errorf("log is full")
	}
}

func (l *logFile) reopen(date string) {
	if err := l.file.close(); err != nil {
		debug.Erro("close log file[%s] failure: %s", fmt.Sprintf("%s/%s.log", l.logDir, l.date), err)
		return
	}

	l.file.open(fmt.Sprintf("%s/%s.log", l.logDir, date))
	l.date = date
}

func (l *logFile) Start() error {
	stat, err := os.Stat(l.logDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		if err := os.MkdirAll(l.logDir, 0o755); err != nil && !os.IsExist(err) {
			return err
		}
	}

	if !stat.IsDir() {
		return fmt.Errorf("logDir[%s] is file", l.logDir)
	}

	return l.file.open(fmt.Sprintf("%s/%s.log", l.logDir, l.date))
}

func (l *logFile) Listen(ctx context.Context) {
	defer l.ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case now := <-l.ticker.C:
			date := now.Format(time.DateOnly)
			if date != l.date {
				l.reopen(date)
			}
		case logData, ok := <-l.data:
			if !ok {
				if logData != nil {
					if err := l.file.write(logData); err != nil {
						fmt.Println(err.Error())
					}
				}
				return
			}

			if err := l.file.write(logData); err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}

func (l *logFile) Close() {
	if !l.closed.CompareAndSwap(false, true) {
		return
	}
	defer l.file.close()
	close(l.data)
	for logData := range l.data {
		if err := l.file.write(logData); err != nil {
			fmt.Println(err.Error())
		}
	}
}

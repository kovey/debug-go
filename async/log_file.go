package async

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/kovey/debug-go/debug"
)

type logFile struct {
	logDir string
	data   chan []byte
	date   string
	ticker *time.Ticker
	file   *file
	mu     sync.RWMutex
	closed bool
	wait   sync.WaitGroup
}

func newLogFile(logDir string, length int) *logFile {
	return &logFile{logDir: logDir, data: make(chan []byte, length), date: time.Now().Format(time.DateOnly), ticker: time.NewTicker(1 * time.Second), file: &file{}, wait: sync.WaitGroup{}}
}

func (l *logFile) Write(p []byte) (n int, err error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.closed {
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
	newPath := fmt.Sprintf("%s/%s.log", l.logDir, date)
	newFile := &file{}
	if err := newFile.open(newPath); err != nil {
		debug.Erro("open new log file[%s] failure: %s", newPath, err)
		return
	}

	if err := l.file.close(); err != nil {
		debug.Erro("close old log file[%s/%s.log] failure: %s", l.logDir, l.date, err)
	}

	l.file = newFile
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

		// Re-stat after directory creation to get valid stat info
		stat, err = os.Stat(l.logDir)
		if err != nil {
			return err
		}
	}

	if !stat.IsDir() {
		return fmt.Errorf("logDir[%s] is file", l.logDir)
	}

	return l.file.open(fmt.Sprintf("%s/%s.log", l.logDir, l.date))
}

func (l *logFile) Listen(ctx context.Context) {
	l.wait.Add(1)
	defer l.ticker.Stop()
	defer l.wait.Done()

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
				return
			}

			if err := l.file.write(logData); err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}

func (l *logFile) Close() {
	l.mu.Lock()
	if l.closed {
		l.mu.Unlock()
		return
	}
	l.closed = true
	l.mu.Unlock()

	defer l.file.close()
	close(l.data)
	for logData := range l.data {
		if err := l.file.write(logData); err != nil {
			fmt.Println(err.Error())
		}
	}
	l.wait.Wait()
}

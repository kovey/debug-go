package async

import (
	"fmt"
	"os"
	"time"

	"github.com/kovey/debug-go/debug"
)

type logFile struct {
	logDir string
	data   chan []byte
	sig    chan bool
	date   string
	ticker *time.Ticker
	file   *file
}

func newLogFile(logDir string, length int) *logFile {
	return &logFile{logDir: logDir, data: make(chan []byte, length), date: time.Now().Format(time.DateOnly), sig: make(chan bool, 1), ticker: time.NewTicker(1 * time.Second), file: &file{}}
}

func (l *logFile) Write(p []byte) (n int, err error) {
	l.data <- p
	return len(p), nil
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
		if os.IsNotExist(err) {
			return os.MkdirAll(l.logDir, 0755)
		}

		return err
	}

	if !stat.IsDir() {
		return fmt.Errorf("logDir[%s] is file", l.logDir)
	}

	return l.file.open(fmt.Sprintf("%s/%s.log", l.logDir, l.date))
}

func (l *logFile) Listen() {
	defer l.ticker.Stop()
	defer l.file.close()

	for {
		select {
		case <-l.sig:
			return
		case now := <-l.ticker.C:
			date := now.Format(time.DateOnly)
			if date != l.date {
				l.reopen(date)
			}
		case logData := <-l.data:
			if err := l.file.write(logData); err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}

func (l *logFile) Stop() {
	l.sig <- true
}

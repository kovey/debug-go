package async

import (
	"errors"
	"fmt"
	"os"
	"sync"
)

type file struct {
	mu   sync.Mutex
	f    *os.File
	path string
}

func (f *file) close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.closeLocked()
}

func (f *file) closeLocked() error {
	if f.f == nil {
		return nil
	}
	return f.f.Close()
}

func (f *file) open(path string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.openLocked(path)
}

func (f *file) openLocked(path string) error {
	fi, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	f.f = fi
	f.path = path
	return nil
}

func (f *file) write(data []byte) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.writeLocked(data)
}

func (f *file) writeLocked(data []byte) error {
	if f.f == nil {
		if f.path == "" {
			return fmt.Errorf("file is not opened")
		}
		if err := f.openLocked(f.path); err != nil {
			return err
		}
	}

	_, err := f.f.Write(data)
	if errors.Is(err, os.ErrClosed) {
		if f.path == "" {
			return err
		}

		if err := f.openLocked(f.path); err != nil {
			return err
		}

		_, err = f.f.Write(data)
	}

	return err
}

package async

import "os"

type file struct {
	f    *os.File
	path string
}

func (f *file) close() error {
	if f.f == nil {
		return nil
	}

	return f.f.Close()
}

func (f *file) open(path string) error {
	fi, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	f.f = fi
	f.path = path
	return nil
}

func (f *file) write(data []byte) error {
	_, err := f.f.Write(data)
	if err == os.ErrClosed {
		if f.path == "" {
			return err
		}

		if err := f.open(f.path); err != nil {
			return err
		}

		_, err = f.f.Write(data)
	}

	return err
}

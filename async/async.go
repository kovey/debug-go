package async

import (
	"context"

	"github.com/kovey/debug-go/debug"
)

var lf *logFile

func Start(logDir string, length int) error {
	lf = newLogFile(logDir, length)
	if err := lf.Start(); err != nil {
		return err
	}

	debug.SetWriter(lf)
	return nil
}

func Listen(ctx context.Context) {
	if lf == nil {
		return
	}

	lf.Listen(ctx)
}

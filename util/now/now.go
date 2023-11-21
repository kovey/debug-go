package now

import (
	"time"

	"github.com/kovey/debug-go/util/date"
)

func DateTime() string {
	return date.DateTime(time.Now())
}

func Date() string {
	return date.Date(time.Now())
}

func Time() string {
	return date.Time(time.Now())
}

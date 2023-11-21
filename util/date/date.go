package date

import "time"

func DateTime(t time.Time) string {
	return t.Format(time.DateTime)
}

func Date(t time.Time) string {
	return t.Format(time.DateOnly)
}

func Time(t time.Time) string {
	return t.Format(time.TimeOnly)
}

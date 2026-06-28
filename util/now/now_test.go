package now

import (
	"strings"
	"testing"
	"time"
)

func TestDateTime(t *testing.T) {
	result := DateTime()

	// Should match the format "2006-01-02 15:04:05"
	if len(result) != 19 {
		t.Errorf("DateTime() = %q, length should be 19, got %d", result, len(result))
	}

	// Verify it parses correctly
	_, err := time.Parse(time.DateTime, result)
	if err != nil {
		t.Errorf("DateTime() returned unparseable value %q: %v", result, err)
	}
}

func TestDateTimeIsCurrent(t *testing.T) {
	result := DateTime()
	now := time.Now().Format(time.DateTime)

	// They should match within a second (the same format applied to roughly the same time)
	if result[:16] != now[:16] {
		t.Errorf("DateTime() = %q, expected ~ %q (minute-level match)", result, now)
	}
}

func TestDate(t *testing.T) {
	result := Date()

	if len(result) != 10 {
		t.Errorf("Date() = %q, length should be 10, got %d", result, len(result))
	}

	_, err := time.Parse(time.DateOnly, result)
	if err != nil {
		t.Errorf("Date() returned unparseable value %q: %v", result, err)
	}
}

func TestDateIsCurrent(t *testing.T) {
	result := Date()
	expected := time.Now().Format(time.DateOnly)

	if result != expected {
		t.Errorf("Date() = %q, want %q", result, expected)
	}
}

func TestTime(t *testing.T) {
	result := Time()

	if len(result) != 8 {
		t.Errorf("Time() = %q, length should be 8, got %d", result, len(result))
	}

	_, err := time.Parse(time.TimeOnly, result)
	if err != nil {
		t.Errorf("Time() returned unparseable value %q: %v", result, err)
	}
}

func TestTimeFormat(t *testing.T) {
	result := Time()
	parts := strings.Split(result, ":")
	if len(parts) != 3 {
		t.Errorf("Time() = %q, should have 3 colon-separated parts", result)
	}
}

func TestAllReturnValues(t *testing.T) {
	dt := DateTime()
	d := Date()
	tm := Time()

	// The date part of DateTime should match Date
	if dt[:10] != d {
		t.Errorf("DateTime() date part %q should match Date() %q", dt[:10], d)
	}

	// The time part of DateTime should match Time
	if dt[11:] != tm {
		t.Errorf("DateTime() time part %q should match Time() %q", dt[11:], tm)
	}
}

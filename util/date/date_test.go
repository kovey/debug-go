package date

import (
	"testing"
	"time"
)

func TestDateTime(t *testing.T) {
	tm := time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC)
	result := DateTime(tm)
	expected := "2024-01-15 10:30:45"
	if result != expected {
		t.Errorf("DateTime() = %q, want %q", result, expected)
	}
}

func TestDateTimeMidnight(t *testing.T) {
	tm := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	result := DateTime(tm)
	expected := "2024-12-31 00:00:00"
	if result != expected {
		t.Errorf("DateTime() = %q, want %q", result, expected)
	}
}

func TestDate(t *testing.T) {
	tm := time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC)
	result := Date(tm)
	expected := "2024-01-15"
	if result != expected {
		t.Errorf("Date() = %q, want %q", result, expected)
	}
}

func TestDateYearBoundary(t *testing.T) {
	tm := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	result := Date(tm)
	expected := "2024-12-31"
	if result != expected {
		t.Errorf("Date() = %q, want %q", result, expected)
	}
}

func TestTime(t *testing.T) {
	tm := time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC)
	result := Time(tm)
	expected := "10:30:45"
	if result != expected {
		t.Errorf("Time() = %q, want %q", result, expected)
	}
}

func TestTimeEdge(t *testing.T) {
	tm := time.Date(2024, 1, 15, 23, 59, 59, 0, time.UTC)
	result := Time(tm)
	expected := "23:59:59"
	if result != expected {
		t.Errorf("Time() = %q, want %q", result, expected)
	}
}

func TestToTime(t *testing.T) {
	tm, err := ToTime("2024-01-15 10:30:45")
	if err != nil {
		t.Fatalf("ToTime() unexpected error: %v", err)
	}

	if tm.Year() != 2024 {
		t.Errorf("ToTime() year = %d, want 2024", tm.Year())
	}
	if tm.Month() != 1 {
		t.Errorf("ToTime() month = %d, want 1", tm.Month())
	}
	if tm.Day() != 15 {
		t.Errorf("ToTime() day = %d, want 15", tm.Day())
	}
	if tm.Hour() != 10 {
		t.Errorf("ToTime() hour = %d, want 10", tm.Hour())
	}
	if tm.Minute() != 30 {
		t.Errorf("ToTime() minute = %d, want 30", tm.Minute())
	}
	if tm.Second() != 45 {
		t.Errorf("ToTime() second = %d, want 45", tm.Second())
	}
}

func TestToTimeInvalid(t *testing.T) {
	_, err := ToTime("not-a-timestamp")
	if err == nil {
		t.Errorf("ToTime() with invalid input should return an error")
	}
}

func TestToTimeEmpty(t *testing.T) {
	_, err := ToTime("")
	if err == nil {
		t.Errorf("ToTime() with empty string should return an error")
	}
}

func TestToTimeRoundTrip(t *testing.T) {
	original := time.Date(2024, 6, 15, 14, 30, 0, 0, time.Local)
	formatted := DateTime(original)
	parsed, err := ToTime(formatted)
	if err != nil {
		t.Fatalf("ToTime() round-trip error: %v", err)
	}
	if !original.Equal(parsed) {
		t.Errorf("ToTime() round-trip: %v != %v", original, parsed)
	}
}

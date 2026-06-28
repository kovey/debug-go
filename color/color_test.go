package color

import (
	"strings"
	"testing"
)

func TestGreen(t *testing.T) {
	result := Green("hello")
	if !strings.Contains(result, "\033[32m") {
		t.Errorf("Green() should contain green ANSI code, got: %q", result)
	}
	if !strings.Contains(result, "hello") {
		t.Errorf("Green() should contain the message, got: %q", result)
	}
	if !strings.Contains(result, "\033[0m") {
		t.Errorf("Green() should contain reset ANSI code, got: %q", result)
	}
}

func TestYellow(t *testing.T) {
	result := Yellow("warning")
	if !strings.Contains(result, "\033[33m") {
		t.Errorf("Yellow() should contain yellow ANSI code, got: %q", result)
	}
	if !strings.Contains(result, "warning") {
		t.Errorf("Yellow() should contain the message, got: %q", result)
	}
}

func TestRed(t *testing.T) {
	result := Red("error")
	if !strings.Contains(result, "\033[31m") {
		t.Errorf("Red() should contain red ANSI code, got: %q", result)
	}
	if !strings.Contains(result, "error") {
		t.Errorf("Red() should contain the message, got: %q", result)
	}
}

func TestWhite(t *testing.T) {
	result := White("plain")
	if !strings.Contains(result, "\033[37m") {
		t.Errorf("White() should contain white ANSI code, got: %q", result)
	}
	if !strings.Contains(result, "plain") {
		t.Errorf("White() should contain the message, got: %q", result)
	}
}

func TestBlue(t *testing.T) {
	result := Blue("info")
	if !strings.Contains(result, "\033[34m") {
		t.Errorf("Blue() should contain blue ANSI code, got: %q", result)
	}
	if !strings.Contains(result, "info") {
		t.Errorf("Blue() should contain the message, got: %q", result)
	}
}

func TestMagenta(t *testing.T) {
	result := Magenta("debug")
	if !strings.Contains(result, "\033[35m") {
		t.Errorf("Magenta() should contain magenta ANSI code, got: %q", result)
	}
	if !strings.Contains(result, "debug") {
		t.Errorf("Magenta() should contain the message, got: %q", result)
	}
}

func TestCyan(t *testing.T) {
	result := Cyan("trace")
	if !strings.Contains(result, "\033[36m") {
		t.Errorf("Cyan() should contain cyan ANSI code, got: %q", result)
	}
	if !strings.Contains(result, "trace") {
		t.Errorf("Cyan() should contain the message, got: %q", result)
	}
}

func TestFormatStructure(t *testing.T) {
	// The format is: <color_code><content><reset_code>
	tests := []struct {
		name     string
		fn       func(string) string
		colorCode string
	}{
		{"Green", Green, "\033[32m"},
		{"Yellow", Yellow, "\033[33m"},
		{"Red", Red, "\033[31m"},
		{"White", White, "\033[37m"},
		{"Blue", Blue, "\033[34m"},
		{"Magenta", Magenta, "\033[35m"},
		{"Cyan", Cyan, "\033[36m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn("test")
			expectedPrefix := tt.colorCode
			expectedSuffix := "\033[0m"

			if !strings.HasPrefix(result, expectedPrefix) {
				t.Errorf("%s() should start with %q, got: %q", tt.name, expectedPrefix, result)
			}
			if !strings.HasSuffix(result, expectedSuffix) {
				t.Errorf("%s() should end with %q, got: %q", tt.name, expectedSuffix, result)
			}
		})
	}
}

func TestColorWithEmptyString(t *testing.T) {
	// Empty content should still produce color codes
	result := Red("")
	if !strings.HasPrefix(result, "\033[31m") {
		t.Errorf("Red() with empty string should start with color code, got: %q", result)
	}
	if !strings.HasSuffix(result, "\033[0m") {
		t.Errorf("Red() with empty string should end with reset code, got: %q", result)
	}
}

func TestColorUnknown(t *testing.T) {
	result := color(Color(999))
	if result != "" {
		t.Errorf("color() with unknown Color should return empty string, got: %q", result)
	}
}

func TestColorKnown(t *testing.T) {
	result := color(Color_Green)
	if result != "\033[32m" {
		t.Errorf("color(Color_Green) = %q, want %q", result, "\033[32m")
	}
}

func TestBgKnown(t *testing.T) {
	result := Bg(Bg_Red)
	if result != "\033[97;41m" {
		t.Errorf("Bg(Bg_Red) = %q, want %q", result, "\033[97;41m")
	}
}

func TestBgUnknown(t *testing.T) {
	result := Bg(Color(999))
	if result != "" {
		t.Errorf("Bg() with unknown Color should return empty string, got: %q", result)
	}
}

func TestBgNone(t *testing.T) {
	result := Bg(Bg_None)
	if result != "" {
		t.Errorf("Bg(Bg_None) should return empty string, got: %q", result)
	}
}

func TestColorNone(t *testing.T) {
	result := color(Color_None)
	if result != "" {
		t.Errorf("color(Color_None) should return empty string, got: %q", result)
	}
}

func TestColorReset(t *testing.T) {
	result := color(Color_Reset)
	if result != "\033[0m" {
		t.Errorf("color(Color_Reset) = %q, want %q", result, "\033[0m")
	}
}

func TestAllBgValues(t *testing.T) {
	// Verify all defined background colors return non-empty strings
	bgColors := map[Color]string{
		Bg_Green:   "\033[97;42m",
		Bg_White:   "\033[90;47m",
		Bg_Yellow:  "\033[90;43m",
		Bg_Red:     "\033[97;41m",
		Bg_Blue:    "\033[97;44m",
		Bg_Megenta: "\033[97;45m",
		Bg_Cyan:    "\033[97;46m",
	}

	for c, expected := range bgColors {
		result := Bg(c)
		if result != expected {
			t.Errorf("Bg(%d) = %q, want %q", c, result, expected)
		}
	}
}

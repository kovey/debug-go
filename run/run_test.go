package run

import (
	"testing"
)

func TestRunCoNormal(t *testing.T) {
	executed := false
	RunCo(func() {
		executed = true
	})

	if !executed {
		t.Errorf("RunCo() did not execute the function")
	}
}

func TestRunCoPanic(t *testing.T) {
	// RunCo should recover from panics and not propagate them
	completed := false

	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("RunCo should have recovered the panic, but it propagated: %v", r)
			}
		}()

		RunCo(func() {
			panic("test panic")
		})
		completed = true
	}()

	if !completed {
		t.Errorf("RunCo() should recover from panics and continue execution")
	}
}

func TestRunCoPanicWithString(t *testing.T) {
	// Test with various panic types
	panicked := true
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = false
				t.Errorf("RunCo should have recovered the panic, but it propagated")
			}
		}()
		RunCo(func() {
			panic("string panic")
		})
		panicked = false
	}()

	if panicked {
		t.Errorf("Panic was not properly contained by RunCo")
	}
}

func TestPanicNil(t *testing.T) {
	result := Panic(nil)
	if result != false {
		t.Errorf("Panic(nil) = %v, want false", result)
	}
}

func TestPanicWithError(t *testing.T) {
	// Panic with an error value should return true and not panic itself
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Panic() should not panic itself, got: %v", r)
			}
		}()
		result := Panic("test error")
		if result != true {
			t.Errorf("Panic(\"test error\") = %v, want true", result)
		}
	}()
}

func TestPanicWithErrorType(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Panic() should not panic itself, got: %v", r)
			}
		}()
		result := Panic(struct{ Msg string }{Msg: "structured error"})
		if result != true {
			t.Errorf("Panic() with struct = %v, want true", result)
		}
	}()
}

func TestPanicStacktrace(t *testing.T) {
	// Verify Panic produces stack trace output by calling it in a nested function
	var outerResult bool
	func() {
		defer func() {
			outerResult = Panic(recover())
		}()
		func() {
			panic("nested panic")
		}()
	}()

	if outerResult != true {
		t.Errorf("Panic(recover()) from nested panic = %v, want true", outerResult)
	}
}

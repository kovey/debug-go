package util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsFileExists(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	if !IsFile(tmpFile.Name()) {
		t.Errorf("IsFile(%q) = false, want true for existing file", tmpFile.Name())
	}
}

func TestIsFileNotExists(t *testing.T) {
	if IsFile("/nonexistent/path/that/does/not/exist/file.txt") {
		t.Errorf("IsFile() = true, want false for nonexistent file")
	}
}

func TestIsFileEmpty(t *testing.T) {
	if IsFile("") {
		t.Errorf("IsFile(\"\") = true, want false for empty path")
	}
}

func TestIsFileDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "testdir_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.Remove(tmpDir)

	// A directory exists on the filesystem, so os.Stat succeeds (err == nil)
	// and IsFile returns true. This documents the current behavior.
	if !IsFile(tmpDir) {
		t.Errorf("IsFile(%q) = false, want true for directory (os.Stat succeeds)", tmpDir)
	}
}

func TestRunDir(t *testing.T) {
	dir := RunDir()
	if dir == "" {
		t.Errorf("RunDir() returned empty string")
	}
}

func TestRunDirIsAbsolute(t *testing.T) {
	dir := RunDir()
	if !filepath.IsAbs(dir) {
		t.Errorf("RunDir() = %q, should be an absolute path", dir)
	}
}

func TestCurrentDir(t *testing.T) {
	dir := CurrentDir()
	if dir == "" {
		t.Errorf("CurrentDir() returned empty string")
	}
}

func TestCurrentDirIsAbsolute(t *testing.T) {
	dir := CurrentDir()
	if !filepath.IsAbs(dir) {
		t.Errorf("CurrentDir() = %q, should be an absolute path", dir)
	}
}

func TestUnitConstants(t *testing.T) {
	if Unit_Second != 1 {
		t.Errorf("Unit_Second = %d, want 1", Unit_Second)
	}
	if Unit_Minute != 60 {
		t.Errorf("Unit_Minute = %d, want 60", Unit_Minute)
	}
	if Unit_Hour != 3600 {
		t.Errorf("Unit_Hour = %d, want 3600", Unit_Hour)
	}
	if Unit_Day != 86400 {
		t.Errorf("Unit_Day = %d, want 86400", Unit_Day)
	}
}

func TestUnitRelationships(t *testing.T) {
	// Verify unit relationships
	if Unit_Minute != 60*Unit_Second {
		t.Errorf("Unit_Minute should equal 60 * Unit_Second")
	}
	if Unit_Hour != 60*Unit_Minute {
		t.Errorf("Unit_Hour should equal 60 * Unit_Minute")
	}
	if Unit_Day != 24*Unit_Hour {
		t.Errorf("Unit_Day should equal 24 * Unit_Hour")
	}
}

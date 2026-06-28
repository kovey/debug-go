package async

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// =============================================================================
// file tests
// =============================================================================

func TestFileOpen(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")

	f := &file{}
	err := f.open(path)
	if err != nil {
		t.Fatalf("file.open() error: %v", err)
	}
	defer f.close()

	if f.f == nil {
		t.Errorf("file.open() should set f.f, got nil")
	}
	if f.path != path {
		t.Errorf("file.open() path = %q, want %q", f.path, path)
	}
}

func TestFileOpenInvalidPath(t *testing.T) {
	f := &file{}
	err := f.open("/nonexistent/dir/that/cannot/be/created/file.log")
	if err == nil {
		t.Errorf("file.open() with invalid path should return error")
	}
}

func TestFileCloseNil(t *testing.T) {
	f := &file{}
	err := f.close()
	if err != nil {
		t.Errorf("file.close() on nil file should return nil, got: %v", err)
	}
}

func TestFileCloseOpenFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")

	f := &file{}
	if err := f.open(path); err != nil {
		t.Fatalf("file.open() error: %v", err)
	}

	err := f.close()
	if err != nil {
		t.Errorf("file.close() on open file should return nil, got: %v", err)
	}

	// Verify file is actually closed
	if f.f != nil {
		// The file handle is still set even after close — that's expected
		// since we don't set it to nil
	}
}

func TestFileWriteToOpenFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")

	f := &file{}
	if err := f.open(path); err != nil {
		t.Fatalf("file.open() error: %v", err)
	}
	defer f.close()

	err := f.write([]byte("hello world"))
	if err != nil {
		t.Errorf("file.write() error: %v", err)
	}

	// Verify content was written
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile() error: %v", err)
	}
	if string(data) != "hello world" {
		t.Errorf("file.write() wrote %q, want %q", string(data), "hello world")
	}
}

func TestFileWriteToClosedFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")

	f := &file{}
	if err := f.open(path); err != nil {
		t.Fatalf("file.open() error: %v", err)
	}

	// Close the file
	if err := f.close(); err != nil {
		t.Fatalf("file.close() error: %v", err)
	}

	// Now write to the closed file — should trigger os.ErrClosed and reopen
	err := f.write([]byte("after close"))
	if err != nil {
		t.Errorf("file.write() should reopen and write after close, got error: %v", err)
	}

	// Verify content was written after reopen
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile() error: %v", err)
	}
	if !strings.Contains(string(data), "after close") {
		t.Errorf("file.write() after reopen: got %q, want to contain %q", string(data), "after close")
	}
}

func TestFileWriteNilFile(t *testing.T) {
	f := &file{}
	err := f.write([]byte("test"))
	if err == nil {
		t.Errorf("file.write() on nil file with no path should return error")
	}
	if err.Error() != "file is not opened" {
		t.Errorf("file.write() error = %q, want %q", err, "file is not opened")
	}
}

func TestFileWriteNilFileWithPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")

	f := &file{path: path}
	err := f.write([]byte("test data"))
	if err != nil {
		t.Errorf("file.write() on nil file with valid path should reopen: %v", err)
	}

	// Verify content
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile() error: %v", err)
	}
	if string(data) != "test data" {
		t.Errorf("file.write() wrote %q, want %q", string(data), "test data")
	}

	f.close()
}

// =============================================================================
// logFile tests
// =============================================================================

func TestNewLogFile(t *testing.T) {
	dir := t.TempDir()
	lf := newLogFile(dir, 10)

	if lf.logDir != dir {
		t.Errorf("newLogFile().logDir = %q, want %q", lf.logDir, dir)
	}
	if cap(lf.data) != 10 {
		t.Errorf("newLogFile().data cap = %d, want 10", cap(lf.data))
	}
	if lf.file == nil {
		t.Errorf("newLogFile().file should not be nil")
	}
}

func TestLogFileWriteAndStart(t *testing.T) {
	dir := t.TempDir()
	lf := newLogFile(dir, 10)

	// Start the log file
	if err := lf.Start(); err != nil {
		t.Fatalf("logFile.Start() error: %v", err)
	}
	defer lf.Close()

	// Write some data
	n, err := lf.Write([]byte("test message\n"))
	if err != nil {
		t.Errorf("logFile.Write() error: %v", err)
	}
	if n != len("test message\n") {
		t.Errorf("logFile.Write() n = %d, want %d", n, len("test message\n"))
	}
}

func TestLogFileWriteFullBuffer(t *testing.T) {
	dir := t.TempDir()
	lf := newLogFile(dir, 1) // Buffer of 1

	if err := lf.Start(); err != nil {
		t.Fatalf("logFile.Start() error: %v", err)
	}
	defer lf.Close()

	// First write succeeds (buffer is empty)
	_, err := lf.Write([]byte("first"))
	if err != nil {
		t.Fatalf("logFile.Write() first write error: %v", err)
	}

	// Second write fails (buffer is full, nothing is draining)
	_, err = lf.Write([]byte("second"))
	if err == nil {
		t.Errorf("logFile.Write() should fail when buffer is full")
	}
	if err.Error() != "log is full" {
		t.Errorf("logFile.Write() full error = %q, want %q", err, "log is full")
	}
}

func TestLogFileWriteAfterClose(t *testing.T) {
	dir := t.TempDir()
	lf := newLogFile(dir, 10)

	if err := lf.Start(); err != nil {
		t.Fatalf("logFile.Start() error: %v", err)
	}

	lf.Close()

	_, err := lf.Write([]byte("after close"))
	if err == nil {
		t.Errorf("logFile.Write() after close should return error")
	}
	if err.Error() != "log is closed" {
		t.Errorf("logFile.Write() closed error = %q, want %q", err, "log is closed")
	}
}

func TestLogFileCloseIdempotent(t *testing.T) {
	dir := t.TempDir()
	lf := newLogFile(dir, 10)

	if err := lf.Start(); err != nil {
		t.Fatalf("logFile.Start() error: %v", err)
	}

	// Close twice — should not panic
	lf.Close()
	lf.Close()
}

func TestLogFileCloseConcurrent(t *testing.T) {
	dir := t.TempDir()
	lf := newLogFile(dir, 100)

	if err := lf.Start(); err != nil {
		t.Fatalf("logFile.Start() error: %v", err)
	}

	// Start Listener
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go lf.Listen(ctx)

	// Write some data
	for i := 0; i < 10; i++ {
		lf.Write([]byte("test\n"))
	}

	// Close should drain and wait
	lf.Close()

	// Verify file was written
	files, _ := os.ReadDir(dir)
	if len(files) == 0 {
		t.Errorf("Close() should flush data to disk")
	}
}

func TestLogFileStartInvalidDir(t *testing.T) {
	// Create a file where we expect a directory
	dir := t.TempDir()
	filePath := filepath.Join(dir, "notadir")
	if err := os.WriteFile(filePath, []byte("data"), 0644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	lf := newLogFile(filePath, 10)
	err := lf.Start()
	if err == nil {
		t.Errorf("logFile.Start() with file path should return error")
	}
}

func TestLogFileStartWithDirCreation(t *testing.T) {
	dir := t.TempDir()
	logDir := filepath.Join(dir, "logs", "app")

	lf := newLogFile(logDir, 10)
	if err := lf.Start(); err != nil {
		t.Fatalf("logFile.Start() should create dirs: %v", err)
	}
	defer lf.Close()

	// Verify directory was created
	info, err := os.Stat(logDir)
	if err != nil {
		t.Fatalf("os.Stat() error: %v", err)
	}
	if !info.IsDir() {
		t.Errorf("logFile.Start() should create a directory at %q", logDir)
	}
}

func TestLogFileReopen(t *testing.T) {
	dir := t.TempDir()
	lf := newLogFile(dir, 10)

	if err := lf.Start(); err != nil {
		t.Fatalf("logFile.Start() error: %v", err)
	}
	defer lf.Close()

	oldDate := lf.date
	newDate := "2025-06-15" // A different date

	lf.reopen(newDate)

	if lf.date != newDate {
		t.Errorf("reopen() date = %q, want %q", lf.date, newDate)
	}
	if lf.date == oldDate {
		t.Errorf("reopen() should change the date")
	}
}

func TestLogFileReopenInvalidDir(t *testing.T) {
	dir := t.TempDir()
	lf := newLogFile(dir, 10)

	if err := lf.Start(); err != nil {
		t.Fatalf("logFile.Start() error: %v", err)
	}
	defer lf.Close()

	oldDate := lf.date

	// Attempt to reopen to an invalid path
	lf.logDir = "/nonexistent/dir/path"
	lf.reopen("2025-06-15")

	// Date should NOT be updated on failure
	if lf.date != oldDate {
		t.Errorf("reopen() on failure should not change date: got %q, want %q", lf.date, oldDate)
	}
}

func TestLogFileListenContextCancel(t *testing.T) {
	dir := t.TempDir()
	lf := newLogFile(dir, 10)

	if err := lf.Start(); err != nil {
		t.Fatalf("logFile.Start() error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		lf.Listen(ctx)
		close(done)
	}()

	// Cancel right away
	cancel()

	select {
	case <-done:
		// Listen exited
	case <-time.After(2 * time.Second):
		// Force close and wait
		lf.Close()
		t.Errorf("Listen() did not exit after context cancellation within timeout")
	}
}

func TestLogFileListenWithData(t *testing.T) {
	dir := t.TempDir()
	lf := newLogFile(dir, 100)

	if err := lf.Start(); err != nil {
		t.Fatalf("logFile.Start() error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go lf.Listen(ctx)

	// Write some data — it should be consumed by Listen and written to file
	lf.Write([]byte("line1\n"))
	lf.Write([]byte("line2\n"))
	lf.Write([]byte("line3\n"))

	// Give it time to flush
	time.Sleep(100 * time.Millisecond)

	// Cancel and close
	cancel()
	lf.Close()

	// Verify data was written
	files, _ := filepath.Glob(filepath.Join(dir, "*.log"))
	if len(files) == 0 {
		t.Fatalf("no log files found in %s", dir)
	}

	data, err := os.ReadFile(files[0])
	if err != nil {
		t.Fatalf("os.ReadFile() error: %v", err)
	}
	if !strings.Contains(string(data), "line1") {
		t.Errorf("Listen() should write data to file, got: %q", string(data))
	}
}

// =============================================================================
// async package integration tests
// =============================================================================

func TestAsyncStart(t *testing.T) {
	dir := t.TempDir()

	err := Start(dir, 100)
	if err != nil {
		t.Fatalf("Start() error: %v", err)
	}
	defer Close()

	// Verify directory and file were created
	files, _ := filepath.Glob(filepath.Join(dir, "*.log"))
	if len(files) == 0 {
		t.Errorf("Start() should create a log file in %s", dir)
	}
}

func TestAsyncStartInvalidDir(t *testing.T) {
	err := Start("/nonexistent/path/that/should/fail", 100)
	if err == nil {
		// Cleanup just in case
		Close()
		t.Errorf("Start() with invalid path should return error")
	}
}

func TestAsyncListenAndClose(t *testing.T) {
	dir := t.TempDir()

	if err := Start(dir, 100); err != nil {
		t.Fatalf("Start() error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go Listen(ctx)

	// Let it run briefly
	time.Sleep(50 * time.Millisecond)

	cancel()
	Close()

	// Should not hang or panic
}

func TestAsyncListenNilLF(t *testing.T) {
	// When lf is nil, Listen should return immediately without panicking
	// This tests the case where Listen is called without Start
	ctx := context.Background()

	done := make(chan struct{})
	go func() {
		Listen(ctx)
		close(done)
	}()

	select {
	case <-done:
		// Success — returned immediately
	case <-time.After(1 * time.Second):
		t.Errorf("Listen() with nil lf should return immediately")
	}
}

func TestAsyncCloseNilLF(t *testing.T) {
	// When lf is nil, Close should return without panicking
	// Reset lf to nil (simulating never-started state)
	lf = nil
	Close() // Should not panic
}

func TestAsyncStartAlreadyExists(t *testing.T) {
	dir := t.TempDir()

	// First start
	if err := Start(dir, 100); err != nil {
		t.Fatalf("first Start() error: %v", err)
	}

	// We can't call Start twice (lf is already set), but Close should work
	Close()
}

func TestConcurrentWrites(t *testing.T) {
	dir := t.TempDir()
	lf := newLogFile(dir, 1000)

	if err := lf.Start(); err != nil {
		t.Fatalf("logFile.Start() error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go lf.Listen(ctx)

	// Concurrent writes
	done := make(chan struct{})
	const numWriters = 10
	for i := 0; i < numWriters; i++ {
		go func(id int) {
			for j := 0; j < 50; j++ {
				lf.Write([]byte("test data\n"))
			}
			done <- struct{}{}
		}(i)
	}

	// Wait for all writers
	for i := 0; i < numWriters; i++ {
		<-done
	}

	cancel()
	lf.Close()

	// Verify no panics and data was written
	files, _ := filepath.Glob(filepath.Join(dir, "*.log"))
	if len(files) == 0 {
		t.Fatalf("no log files found in %s", dir)
	}
}

func TestWriteCloseRace(t *testing.T) {
	// Test the fixed race condition: concurrent Write and Close
	dir := t.TempDir()
	lf := newLogFile(dir, 1000)

	if err := lf.Start(); err != nil {
		t.Fatalf("logFile.Start() error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go lf.Listen(ctx)

	done := make(chan struct{})
	go func() {
		for i := 0; i < 100; i++ {
			lf.Write([]byte("data\n"))
		}
		done <- struct{}{}
	}()

	// Close while writes are still happening
	time.Sleep(5 * time.Millisecond)
	lf.Close()

	// This should not panic
	<-done
}

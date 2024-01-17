package debug_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"testing"

	"github.com/FollowTheProcess/debug"
	"github.com/FollowTheProcess/test"
)

var (
	anInt   = 2
	aBool   = true
	aString = "hello world"
)

func TestDebug(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not get runtime caller")
	}

	runDebug := func() {
		debug.Debug(anInt)
		debug.Debug(aBool)
		debug.Debug(aString)
	}

	out := captureStderr(t, runDebug)

	want := fmt.Sprintf(`DEBUG: [%[1]s:28:3] anInt = 2
DEBUG: [%[1]s:29:3] aBool = true
DEBUG: [%[1]s:30:3] aString = hello world
`, file)

	test.Diff(t, out, want)
}

func captureStderr(t *testing.T, printer func()) string {
	t.Helper()
	old := os.Stderr // Backup of the real one
	defer func() {
		os.Stderr = old // Set it back even if we error later
	}()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() returned an error: %v", err)
	}

	// Set stdout to our new pipe
	os.Stderr = w

	capture := make(chan string)
	// Copy in a goroutine so printing can't block forever
	go func() {
		buf := new(bytes.Buffer)
		io.Copy(buf, r) //nolint: errcheck
		capture <- buf.String()
	}()

	// Call our test function that prints to stdout
	printer()

	// Close the writer
	w.Close()
	captured := <-capture

	return captured
}

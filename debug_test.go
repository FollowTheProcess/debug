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

// Keep this on the top of this file so it doesn't mess up the line numbers in the testcases.
func TestDebug(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not get runtime caller")
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			run := func() {
				debug.Debug(tt.arg)
			}
			got := captureStderr(t, run)
			want := fmt.Sprintf(tt.want, file) + "\n"
			test.Equal(t, got, want)
		})
	}
}

// testcase is a single debug test case.
type testcase struct {
	name string // Identifier for debugging
	arg  any    // What to pass to debug.Debug
	want string // The expected debug output
}

// Person is a fake exported struct for testing.
type Person struct {
	Exported    string
	notExported string
}

// List of testcases for TestDebug, keep this below TestDebug as that way adding more cases
// won't change any of the line numbers in existing tests.
//
// want contains a single %s fmt print verb to be replaced by the filename of this test file
// during each test.
var testcases = []testcase{
	{
		name: "int",
		arg:  2,
		want: "DEBUG: [%s:25:5] tt.arg = 2",
	},
	{
		name: "float",
		arg:  3.14159,
		want: "DEBUG: [%s:25:5] tt.arg = 3.14159",
	},
	{
		name: "bool",
		arg:  true,
		want: "DEBUG: [%s:25:5] tt.arg = true",
	},
	{
		name: "string",
		arg:  "hello world",
		want: `DEBUG: [%s:25:5] tt.arg = "hello world"`,
	},
	{
		name: "slice",
		arg:  []int{1, 2, 3, 4},
		want: "DEBUG: [%s:25:5] tt.arg = []int{1, 2, 3, 4}",
	},
	{
		name: "anonymous struct unexported fields",
		arg:  struct{ name string }{name: "dave"},
		want: `DEBUG: [%s:25:5] tt.arg = struct { name string }{name:"dave"}`,
	},
	{
		name: "anonymous struct exported fields",
		arg:  struct{ Name string }{Name: "dave"},
		want: `DEBUG: [%s:25:5] tt.arg = struct { Name string }{Name:"dave"}`,
	},
	{
		name: "struct with mixed fields",
		arg: Person{
			Exported:    "yes",
			notExported: "no",
		},
		want: `DEBUG: [%s:25:5] tt.arg = debug_test.Person{Exported:"yes", notExported:"no"}`,
	},
	{
		name: "map",
		arg:  map[string]bool{"good": true, "bad": false},
		want: `DEBUG: [%s:25:5] tt.arg = map[string]bool{"bad":false, "good":true}`,
	},
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

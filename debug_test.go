package debug_test

import (
	"testing"

	"github.com/FollowTheProcess/debug"
)

func TestHello(t *testing.T) {
	got := debug.Hello()
	want := "Hello debug"

	if got != want {
		t.Errorf("got %s, wanted %s", got, want)
	}
}

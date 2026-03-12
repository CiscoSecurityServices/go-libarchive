package guestfs

import (
	"testing"
)

func TestVersion(t *testing.T) {
	if VersionString() == "" {
		t.Fatal("Bad version String")
	}
}

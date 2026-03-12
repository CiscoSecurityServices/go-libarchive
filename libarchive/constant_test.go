package libarchive

import (
	"testing"
)

func TestVersion(t *testing.T) {
	if LibArchiveVersion == "" {
		t.Fatal("Wrong version String")
	}
	if LibArchiveVersionInt < 3000000 {
		t.Fatal("Wrong version number")
	}
}

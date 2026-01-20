package libarchive

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"
)

type archive_test_data struct {
	path    string
	name    string
	symlink string
	size    int64
	mode    os.FileMode
	data    []byte
}

func assertArchivesData(t *testing.T, testFile io.Reader, lastErr error, expected []archive_test_data) {
	reader, err := NewReader(testFile)
	if err != nil {
		t.Fatalf("Error on creating Archive from a io.Reader:\n %s", err)
	}
	defer func() {
		err := reader.Close()
		if err != nil {
			t.Fatalf("Error on reader Close:\n %s", err)
		}
	}()

	for i, expectedEntry := range expected {
		entry, err := reader.Next()

		if err != nil {
			t.Fatalf("%d - got error on reader.Next(): %s", i, err)
		}

		if reader.IsRaw() {
			t.Errorf("%d - expected archive data to NOT be raw", i)
		}

		name := entry.PathName()
		if name != expectedEntry.path {
			t.Errorf("%d - got %s expected %s as PathName", i, name, expectedEntry.path)
		}
		symlinkToNothing := entry.Symlink()
		if symlinkToNothing != expectedEntry.symlink {
			t.Errorf("%d - got %s expected %s as Symlink", i, symlinkToNothing, expectedEntry.symlink)
		}
		infoA := entry.Stat()
		if infoA.Name() != expectedEntry.name {
			t.Errorf("%d - got %s expected %s as Name", i, infoA.Name(), expectedEntry.name)
		}
		if infoA.Size() != expectedEntry.size {
			t.Errorf("%d - got %d expected %d as Size", i, infoA.Size(), expectedEntry.size)
		}
		if infoA.Mode() != expectedEntry.mode {
			t.Errorf("%d - got %v expected %v as Mode", i, infoA.Mode(), expectedEntry.mode)
		}

		allBytes, err := io.ReadAll(reader)
		if err != nil {
			t.Errorf("%d - got error on reader.Read():\n%s", i, err)
		}
		if int64(len(allBytes)) != expectedEntry.size {
			t.Errorf("%d - got %d as size of the read but expected %d", i, len(allBytes), expectedEntry.size)
		}

		if expectedEntry.data != nil {
			if !bytes.Equal(allBytes, expectedEntry.data) {
				t.Errorf("%d - The contents:\n [%s] are not the expectedContent:\n [%s]", i, allBytes, expectedEntry.data)
			}
		}

		if t.Failed() {
			t.FailNow()
		}
	}

	if lastErr == nil {
		lastErr = ErrArchiveEOF
	}
	_, err = reader.Next()
	if !errors.Is(err, lastErr) {
		t.Fatalf("Last reader.Next(): got %v expected %v", err, lastErr)
	}
}

func assertDefaultCompressed(t *testing.T, file string) {
	testFile, err := os.Open(file)
	if err != nil {
		t.Fatalf("Error while reading fixture file %s ", err)
	}

	reader, err := NewReader(testFile)
	if err != nil {
		t.Fatalf("Error while creating NewReader %s ", err)
	}

	defer func() {
		err := reader.Close()
		if err != nil {
			t.Fatalf("Error on reader Close:\n %s", err)
		}
	}()
	//--------------a-------------
	_, err = reader.Next()
	if err != nil {
		t.Fatalf("got error on reader.Next() first:\n%s", err)
	}
	if !reader.IsRaw() {
		t.Fatalf("expected compressed data to be raw")
	}

	b := make([]byte, 512)
	size, err := reader.Read(b)
	if err != nil {
		t.Fatalf("got error on reader.Read():\n%s", err)
	}
	if size != 14 {
		t.Fatalf("got %d as size of the read but expected %d", size, 14)
	}

	expectedContent := []byte("Sha lalal lal\n")
	if !bytes.Equal((b[:size]), expectedContent) {
		t.Fatalf("The contents:\n [%s] are not the expectedContent:\n [%s]", b[:size], expectedContent)
	}
	//--------------a-------------

	_, err = reader.Next()
	if err != ErrArchiveEOF {
		t.Fatalf("Expected EOF on second reader.Next() got err :\n %s", err)
	}
}

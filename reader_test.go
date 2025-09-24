package libarchive

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"
)

func TestNewArchive(t *testing.T) {
	assertArchive(t, "./fixtures/test.tar")
}

func TestNewArchiveCompressed(t *testing.T) {
	assertArchive(t, "./fixtures/test.tar.gz")
}

func TestCompressedGz(t *testing.T) {
	assertCompressed(t, "./fixtures/a.gz")
}

func TestCompressedBz2(t *testing.T) {
	assertCompressed(t, "./fixtures/a.bz2")
}

func TestFooA(t *testing.T) {
	testFile, err := os.Open("./fixtures/foo.a")
	if err != nil {
		t.Fatalf("Error while reading fixture file %s ", err)
	}

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

	v, err := reader.Next()
	if err != nil {
		t.Fatalf("Error getting next header:\n %s ", err)
	}
	if v.PathName() != "foo1" {
		t.Fatalf("Expected foo1 got %s", v.PathName())
	}
	if v.Stat().Size() != 5 {
		t.Fatalf("Expected size 5 got %d", v.Stat().Size())
	}

	v, err = reader.Next()
	if err != nil {
		t.Fatalf("Error getting next header:\n %s ", err)
	}
	if v.PathName() != "foo2" {
		t.Fatalf("Expected foo2 got %s", v.PathName())
	}
	if v.Stat().Size() != 5 {
		t.Fatalf("Expected size 5 got %d", v.Stat().Size())
	}

	_, err = reader.Next()
	if !errors.Is(err, ErrArchiveEOF) {
		t.Fatalf("Expected EOF got %s", err)
	}
}

func TestSomeA(t *testing.T) {
	testFile, err := os.ReadFile("./fixtures/foo.a")
	if err != nil {
		t.Fatalf("Error while reading fixture file %s ", err)
	}
	testFile[133] = 0xff // Corrupt the file a bit
	fmt.Println("Length: ", len(testFile))

	reader, err := NewReader(bytes.NewReader(testFile))
	if err != nil {
		t.Fatalf("Error on creating Archive from a io.Reader:\n %s", err)
	}
	defer func() {
		err := reader.Close()
		if err != nil {
			t.Fatalf("Error on reader Close:\n %s", err)
		}
	}()

	v, err := reader.Next()
	if err != nil {
		t.Fatalf("Error getting next header:\n %s ", err)
	}
	if v.PathName() != "foo1" {
		t.Fatalf("Expected foo1 got %s", v.PathName())
	}
	if v.Stat().Size() != 5 {
		t.Fatalf("Expected size 5 got %d", v.Stat().Size())
	}

	_, err = reader.Next()
	if !errors.Is(err, ErrInvalidHeaderSignature) {
		t.Fatalf("Expected ErrInvalidHeaderSignature got %s", err)
	}
}

func TestTwoReaders(t *testing.T) {
	testFile, err := os.Open("./fixtures/test.tar")
	if err != nil {
		t.Fatalf("Error while reading fixture file %s ", err)
	}

	_, err = NewReader(testFile)
	if err != nil {
		t.Fatalf("Error creating Archive from a io.Reader 1:\n %s ", err)
	}

	testFile2, err := os.Open("./fixtures/test2.tar")
	if err != nil {
		t.Fatalf("Error while reading fixture file %s ", err)
	}

	_, err = NewReader(testFile2)
	if err != nil {
		t.Fatalf("Error on creating Archive from a io.Reader 2:\n %s", err)
	}
}

// func TestBadArchive(t *testing.T) {
// 	testFile, err := os.Open("./fixtures/bad.tar")
// 	if err != nil {
// 		t.Fatalf("Error while reading fixture file %s ", err)
// 	}

// 	reader, err := NewReader(testFile)
// 	if err != nil {
// 		t.Fatalf("Error creating Archive from a io.Reader 1:\n %s ", err)
// 	}

// 	warn := false
// 	for {
// 		_, err := reader.Next()
// 		if errors.Is(err, ErrArchiveEOF) {
// 			break
// 		}
// 		if errors.Is(err, ErrArchiveWarn) {
// 			warn = true
// 			fmt.Println(err)
// 			continue
// 		}
// 		if err != nil {
// 			t.Fatalf("Error getting next header:\n %s ", err)
// 		}
// 	}
// 	if !warn {
// 		t.Fatalf("Expected a warning but got none")
// 	}
// }

func assertCompressed(t *testing.T, file string) {
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
	// Should use any of the data below if its raw
	// nameA := entryA.PathName()
	// if nameA != "a" {
	// 	t.Fatalf("got %s expected %s as Name of the first entry", nameA, "a")
	// }
	// symlinkToNothing := entryA.Symlink()
	// if symlinkToNothing != "" {
	// 	t.Fatalf("got %s expected %s as Symlink of the first entry", symlinkToNothing, "a")
	// }
	// infoA := entryA.Stat()
	// if infoA.Name() != "a" {
	// 	t.Fatalf("got %s expected %s as Name of the first entry", infoA.Name(), "a")
	// }
	// if infoA.Size() != 0 {
	// 	t.Fatalf("got %d expected %d as Size of the first entry", infoA.Size(), 0)
	// }
	// if infoA.Mode() != 0644 {
	// 	t.Fatalf("got %v expected %v as Mode of the first entry", infoA.Mode(), 0644)
	// }

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

func assertArchive(t *testing.T, file string) {
	testFile, err := os.Open(file)
	if err != nil {
		t.Fatalf("Error while reading fixture file %s ", err)
	}

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
	//--------------a-------------
	entryA, err := reader.Next()
	if err != nil {
		t.Fatalf("got error on reader.Next() first:\n%s", err)
	}
	if reader.IsRaw() {
		t.Fatalf("expected archive data to NOT be raw")
	}
	nameA := entryA.PathName()
	if nameA != "a" {
		t.Fatalf("got %s expected %s as Name of the first entry", nameA, "a")
	}
	symlinkToNothing := entryA.Symlink()
	if symlinkToNothing != "" {
		t.Fatalf("got %s expected %s as Symlink of the first entry", symlinkToNothing, "a")
	}
	infoA := entryA.Stat()
	if infoA.Name() != "a" {
		t.Fatalf("got %s expected %s as Name of the first entry", infoA.Name(), "a")
	}
	if infoA.Size() != 14 {
		t.Fatalf("got %d expected %d as Size of the first entry", infoA.Size(), 14)
	}
	if infoA.Mode() != 0664 {
		t.Fatalf("got %v expected %v as Mode of the first entry", infoA.Mode(), 0664)
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

	//--------------b-------------
	// should be symlink
	entryB, err := reader.Next()
	if err != nil {
		t.Fatalf("got error on reader.Next() second:\n%s", err)
	}
	if reader.IsRaw() {
		t.Fatalf("expected archive data to NOT be raw")
	}
	nameB := entryB.PathName()
	if nameB != "b" {
		t.Fatalf("got %s expected %s as Name of the second entry", nameB, "b")
	}
	symlinkToA := entryB.Symlink()
	if symlinkToA != "a" {
		t.Fatalf("got %s expected %s as Symlink of the second entry", symlinkToA, "a")
	}
	infoB := entryB.Stat()
	if infoB.Name() != "b" {
		t.Fatalf("got %s expected %s as Name of the second entry", infoB.Name(), "a")
	}
	if infoB.Size() != 0 {
		t.Fatalf("got %d expected %d as Size of the second entry", infoB.Size(), 0)
	}
	eModeB := 0777 | os.ModeSymlink
	if infoB.Mode() != eModeB {
		t.Fatalf("got %v expected %v as Mode of the second entry", infoB.Mode(), eModeB)
	}
	//--------------b-------------

	_, err = reader.Next()
	if err != ErrArchiveEOF {
		t.Fatalf("Expected EOF on second reader.Next() got err :\n %s", err)
	}
}

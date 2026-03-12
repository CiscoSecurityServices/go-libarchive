package libarchive

import (
	"bytes"
	"os"
	"testing"
)

func TestNewArchive(t *testing.T) {
	testData := []archive_test_data{
		{path: "a", name: "a", size: 14, mode: 0644, data: []byte("Sha lalal lal\n")},
		{path: "b", name: "b", size: 0, symlink: "a", mode: 0755 | os.ModeSymlink, data: nil},
	}

	for _, file := range []archive_test{
		{path: "test.tar"},
		{path: "test.tar.gz"},
		{path: "test.cpio"},
		{path: "test.tar", truncate: 1024 + 512},
		{path: "test.cpio", expErr: ErrUnexpectedEOF, truncate: 360},
	} {
		t.Run("Testing "+file.path, func(t *testing.T) {
			data, err := os.ReadFile("./fixtures/" + file.path)
			if err != nil {
				t.Fatalf("Error while reading fixture file %s ", err)
			}
			if file.truncate > 0 {
				data = data[:file.truncate]
			}

			assertArchivesData(t, bytes.NewReader(data), file.expErr, testData)
		})
	}
}

func TestCompressedGz(t *testing.T) {
	assertDefaultCompressed(t, "./fixtures/a.gz")
}

func TestCompressedBz2(t *testing.T) {
	assertDefaultCompressed(t, "./fixtures/a.bz2")
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

func TestFooA(t *testing.T) {
	testFile, err := os.Open("./fixtures/foo.a")
	if err != nil {
		t.Fatalf("Error while reading fixture file %s ", err)
	}

	testData := []archive_test_data{
		{path: "foo1", name: "foo1", size: 5, mode: 0644, data: []byte("foo1\n")},
		{path: "foo2", name: "foo2", size: 5, mode: 0644, data: []byte("foo2\n")},
	}
	assertArchivesData(t, testFile, nil, testData)
}

func TestCorruptA(t *testing.T) {
	data, err := os.ReadFile("./fixtures/foo.a")
	if err != nil {
		t.Fatalf("Error while reading fixture file %s ", err)
	}
	data[133] = 0xff // Corrupt the file a bit

	testData := []archive_test_data{
		{path: "foo1", name: "foo1", size: 5, mode: 0644, data: []byte("foo1\n")},
	}
	assertArchivesData(t, bytes.NewReader(data), ErrInvalidHeaderSignature, testData)
}

func TestIso(t *testing.T) {
	file, err := os.Open("./fixtures/foo.iso")
	if err != nil {
		t.Fatalf("Error opening fixture file %s ", err)
	}
	defer file.Close()

	testData := []archive_test_data{
		{path: ".", name: ".", size: 0, mode: 0755 | os.ModeDir},
		{path: "bar", name: "bar", size: 4, mode: 0644, data: []byte("foo\n")},
		{path: "bar-hl", name: "bar-hl", size: 0, mode: 0644, hardlink: "bar"},
		{path: "baz", name: "baz", size: 4, mode: 0644, data: []byte("baz\n")},
		{path: "bar-sl", name: "bar-sl", size: 0, mode: 0777 | os.ModeSymlink, symlink: "bar"},
	}
	assertArchivesData(t, file, nil, testData)
}

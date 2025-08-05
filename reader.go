package libarchive

/*
#cgo pkg-config: libarchive
#include <archive.h>
#include <archive_entry.h>
#include <stdlib.h>
#include "reader.h"
*/
import "C"
import (
	"io"
	"unsafe"
)

// Reader represents libarchive archive
type Reader struct {
	archive *C.struct_archive
	reader  io.Reader // the io.Reader from which we Read
	buffer  []byte    // buffer for the raw reading
	index   int64     // current reading index
}

// NewReader returns new Archive by calling archive_read_open
func NewReader(reader io.Reader) (r *Reader, err error) {
	r = &Reader{
		archive: C.archive_read_new(),
		reader:  reader,
		buffer:  make([]byte, 1024),
	}
	C.archive_read_support_filter_all(r.archive)
	C.archive_read_support_format_all(r.archive)
	C.archive_read_support_format_raw(r.archive)

	e := C.go_libarchive_open(r.archive, (*C.char)(unsafe.Pointer(r)))

	err = codeToError(r.archive, int(e))
	return
}

//export myopen
func myopen(archive *C.struct_archive, client_data unsafe.Pointer) C.int {
	// actually write something
	return ARCHIVE_OK
}

//export myclose
func myclose(archive *C.struct_archive, client_data unsafe.Pointer) C.int {
	// actually write something
	return ARCHIVE_OK
}

//export myread
func myread(archive *C.struct_archive, client_data *C.char, block unsafe.Pointer) C.size_t {
	reader := (*Reader)(unsafe.Pointer(client_data))
	read, err := reader.reader.Read(reader.buffer)
	if err != nil && err != ErrArchiveEOF {
		// set error
		read = -1
	}

	*(*uintptr)(block) = uintptr(unsafe.Pointer(&reader.buffer[0]))

	return C.size_t(read)
}

// Next calls archive_read_next_header and returns an
// interpretation of the ArchiveEntry which is a wrapper around
// libarchive's archive_entry, or Err.
//
// ErrArchiveEOF is returned when there
// is no more to be read from the archive
func (r *Reader) Next() (ArchiveEntry, error) {
	e := new(entryImpl)

	errno := int(C.archive_read_next_header(r.archive, &e.entry))
	err := codeToError(r.archive, errno)

	if err != nil {
		e = nil
	}

	return e, err
}

// func (r *Reader) Position() int64 {
// 	return int64(C.archive_position(r.archive))
// }

// Must be called after Next
func (r *Reader) IsRaw() bool {
	format := C.archive_format(r.archive)
	return int(format) == 589824
}

// Read calls archive_read_data which reads the current archive_entry.
// It acts as io.Reader.Read in any other aspect
func (r *Reader) Read(b []byte) (n int, err error) {
	n = int(C.archive_read_data(r.archive, unsafe.Pointer(&b[0]), C.size_t(cap(b))))
	if n == 0 {
		err = ErrArchiveEOF
	} else if 0 > n { // err
		err = codeToError(r.archive, ARCHIVE_FAILED)
		n = 0
	}
	r.index += int64(n)
	return
}

// Size returns compressed size of the current archive entry
func (r *Reader) Size() int {
	return int(C.archive_filter_bytes(r.archive, C.int(0)))
}

// Free frees the resources the underlying libarchive archive is using
// calling archive_read_free
// Note this calls
func (r *Reader) ReadFree() error {
	if C.archive_read_free(r.archive) == ARCHIVE_FATAL {
		return ErrArchiveFatalClosing
	}
	return nil
}

// Close closes the underlying libarchive archive
// calling archive read_cloe
func (r *Reader) ReadClose() error {
	if C.archive_read_close(r.archive) == ARCHIVE_FATAL {
		return ErrArchiveFatalClosing
	}
	return nil
}

// Close closes the underlying libarchive archive and frees it, changing the names since its more common in go to always call Close
func (r *Reader) Close() error {
	err := r.ReadClose()
	if err != nil {
		return err
	}
	return r.ReadFree()
}

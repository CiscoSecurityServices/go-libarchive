package libarchive

// #include <archive.h>
import "C"
import (
	"fmt"
	"io"
	"strings"
)

const (
	ARCHIVE_EOF    = C.ARCHIVE_EOF
	ARCHIVE_OK     = C.ARCHIVE_OK
	ARCHIVE_RETRY  = C.ARCHIVE_RETRY
	ARCHIVE_WARN   = C.ARCHIVE_WARN
	ARCHIVE_FAILED = C.ARCHIVE_FAILED
	ARCHIVE_FATAL  = C.ARCHIVE_FATAL
)

var (
	ErrArchiveEOF             = io.EOF
	ErrArchiveFatal           = wrapError{level: "FATAL"}
	ErrArchiveFailed          = wrapError{level: "FAILED"}
	ErrArchiveRetry           = wrapError{level: "RETRY"}
	ErrArchiveWarn            = wrapError{level: "WARN"}
	ErrArchiveFatalClosing    = ErrArchiveFatal.wrap("critical error, archive closing")
	ErrInvalidHeaderSignature = fmt.Errorf("invalid header signature")
)

func codeToError(archive *C.struct_archive, e int) error {
	switch e {
	case ARCHIVE_EOF:
		return ErrArchiveEOF
	case ARCHIVE_FATAL:
		errString := errorString(archive)
		if errString == "Incorrect file header signature" {
			return ErrInvalidHeaderSignature
		}
		return ErrArchiveFatal.wrap(errString)
	case ARCHIVE_FAILED:
		return ErrArchiveFailed.wrap(errorString(archive))
	case ARCHIVE_RETRY:
		return ErrArchiveRetry.wrap(errorString(archive))
	case ARCHIVE_WARN:
		return ErrArchiveWarn.wrap(errorString(archive))
	}
	return nil
}

func errorString(archive *C.struct_archive) string {
	return C.GoString(C.archive_error_string(archive))
}

// -------------------------------Error Wrapper---------------------------------
type wrapError struct {
	err   string
	level string
}

func (err wrapError) Error() string {
	return fmt.Sprintf("%v [%s]", err.prefix(), err.err)
}
func (err wrapError) wrap(inner string) error {
	return wrapError{level: err.level, err: inner}
}
func (err wrapError) prefix() string {
	return fmt.Sprintf("libarchive: %v", err.level)
}

//	func (err wrapError) Unwrap() error {
//		return err.level
//	}
func (err wrapError) Is(target error) bool {
	ts := target.Error()
	return strings.HasPrefix(ts, err.prefix())
}

//-------------------------------Error Wrapper---------------------------------

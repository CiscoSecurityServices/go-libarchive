package libarchive

/*
#cgo pkg-config: libarchive
#include <archive.h>
#include <archive_entry.h>
#include <stdlib.h>
*/
import "C"

import (
	"os"
	"path/filepath"
	"syscall"
)

// ArchiveEntry represents an libarchive archive_entry
type ArchiveEntry interface {
	// FileInfo describing archive_entry
	Stat() os.FileInfo
	// The name of the entry
	PathName() string
	Symlink() string
	Hardlink() string
	IsHardlink() bool
}

type entryImpl struct {
	entry *C.struct_archive_entry
}

type entryInfo struct {
	stat  *C.struct_stat
	entry *C.struct_archive_entry
	name  string
}

func (h *entryImpl) Stat() os.FileInfo {
	info := &entryInfo{}
	info.stat = C.archive_entry_stat(h.entry)
	info.entry = h.entry
	info.name = filepath.Base(h.PathName())
	return info
}

// PathName returns the path name of the entry
func (h *entryImpl) PathName() string {
	name := C.archive_entry_pathname(h.entry)

	return C.GoString(name)
}

// Symlink returns the symlink name of the entry
// returns empty string if no symlink is set
func (h *entryImpl) Symlink() string {
	name := C.archive_entry_symlink(h.entry)

	return C.GoString(name)
}

// Hardlink returns the hardlink name of the entry
// returns empty string if no hardlink is set
func (h *entryImpl) Hardlink() string {
	name := C.archive_entry_hardlink(h.entry)

	return C.GoString(name)
}

func (h *entryImpl) IsHardlink() bool {
	return C.archive_entry_hardlink(h.entry) != nil
}

// Name returns the base name of the entry
func (e *entryInfo) Name() string {
	return e.name
}
func (e *entryInfo) Size() int64 {
	return int64(C.archive_entry_size(e.entry))
}
func (e *entryInfo) Mode() os.FileMode {
	// Use archive_entry_perm for permissions
	mode := os.FileMode(C.archive_entry_perm(e.entry))

	// Use archive_entry_filetype for file type (more reliable for all formats)
	fileType := C.archive_entry_filetype(e.entry)
	switch fileType {
	case syscall.S_IFLNK:
		mode |= os.ModeSymlink
	case syscall.S_IFSOCK:
		mode |= os.ModeSocket
	case syscall.S_IFCHR:
		mode |= os.ModeDevice | os.ModeCharDevice
	case syscall.S_IFBLK:
		mode |= os.ModeDevice
	case syscall.S_IFDIR:
		mode |= os.ModeDir
	case syscall.S_IFIFO:
		mode |= os.ModeNamedPipe
	}
	return mode
}
func (e *entryInfo) IsDir() bool {
	fileType := C.archive_entry_filetype(e.entry)
	return fileType == syscall.S_IFDIR
}
func (e *entryInfo) Sys() interface{} {
	return e.stat
}

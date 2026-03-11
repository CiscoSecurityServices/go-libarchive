package guestfs

/*
#cgo CFLAGS:  -DGUESTFS_PRIVATE=1 -DGUESTFS_NO_WARN_DEPRECATED -UGUESTFS_NO_DEPRECATED
#cgo LDFLAGS: -lguestfs

#include <guestfs.h>
#include <stdio.h>
#include <string.h>

const void get_guestfs_version_string(char* buffer, size_t buffer_size) {
    // Create a guestfs handle
    guestfs_h *g = guestfs_create();
    if (g == NULL) {
        return;
    }

    // Retrieve the version information
    struct guestfs_version *v = guestfs_version(g);

    if (v != NULL) {
        // Ensure the buffer is large enough to hold the version string
        if (buffer_size > 0) {
            // Copy the extra string safely
            char extra[39];
            strncpy(extra, v->extra, sizeof(extra) - 1);
            extra[sizeof(extra) - 1] = '\0'; // Ensure null-termination

            // Format the version string into the buffer
            snprintf(buffer, buffer_size, "%ld.%ld.%ld%s", v->major, v->minor, v->release, extra);
        }
        guestfs_free_version(v);
    }

    guestfs_close(g);
    return;
}
*/
import "C"
import "unsafe"

func VersionString() string {
	buffer := make([]byte, 100)
	C.get_guestfs_version_string((*C.char)(unsafe.Pointer(&buffer[0])), C.size_t(len(buffer)))
	return C.GoString((*C.char)(unsafe.Pointer(&buffer[0])))
}

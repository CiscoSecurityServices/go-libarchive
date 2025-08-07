#include <archive.h>
#include <archive_entry.h>
ssize_t go_libarchive_reader(struct archive *a, void *client_data, const void **block);
int64_t go_libarchive_seek(struct archive * a, void *client_data, int64_t request, int whence);
ssize_t go_libarchive_open(struct archive *a, char *client_data);

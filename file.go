// pdgo File API - unified CGO implementation

package pdgo

/*
#include <stdint.h>

// File API
void* pd_file_open(const char* path, int mode);
int pd_file_close(void* file);
int pd_file_read(void* file, void* buf, uint32_t len);
int pd_file_write(void* file, const void* buf, uint32_t len);
int pd_file_flush(void* file);
int pd_file_tell(void* file);
int pd_file_seek(void* file, int pos, int whence);
int pd_file_mkdir(const char* path);
int pd_file_unlink(const char* path, int recursive);
int pd_file_rename(const char* from, const char* to);
*/
import "C"
import "unsafe"

// SDFile represents an open file
type SDFile struct {
	ptr unsafe.Pointer
}

// FileStat contains file statistics
type FileStat struct {
	IsDir bool
	Size  uint32
	MTime uint32
}

// File provides access to file system functions
type File struct{}

func newFile() *File {
	return &File{}
}

// FileOptions for opening files
type FileOptions int32

const (
	FileRead     FileOptions = 1 << 0
	FileReadData FileOptions = 1 << 1
	FileWrite    FileOptions = 1 << 2
	FileAppend   FileOptions = 1 << 3
)

// Open opens a file
func (f *File) Open(path string, mode FileOptions) (*SDFile, error) {
	cpath := make([]byte, len(path)+1)
	copy(cpath, path)
	ptr := C.pd_file_open((*C.char)(unsafe.Pointer(&cpath[0])), C.int(mode))
	if ptr != nil {
		return &SDFile{ptr: ptr}, nil
	}
	return nil, &fileError{op: "open", path: path}
}

// Close closes a file
func (f *File) Close(file *SDFile) error {
	if file != nil && file.ptr != nil {
		if C.pd_file_close(file.ptr) == 0 {
			file.ptr = nil
			return nil
		}
	}
	return &fileError{op: "close"}
}

// Read reads from a file
func (f *File) Read(file *SDFile, buf []byte) (int, error) {
	if file != nil && file.ptr != nil && len(buf) > 0 {
		n := C.pd_file_read(file.ptr, unsafe.Pointer(&buf[0]), C.uint32_t(len(buf)))
		if n >= 0 {
			return int(n), nil
		}
	}
	return 0, &fileError{op: "read"}
}

// Write writes to a file
func (f *File) Write(file *SDFile, buf []byte) (int, error) {
	if file != nil && file.ptr != nil && len(buf) > 0 {
		n := C.pd_file_write(file.ptr, unsafe.Pointer(&buf[0]), C.uint32_t(len(buf)))
		if n >= 0 {
			return int(n), nil
		}
	}
	return 0, &fileError{op: "write"}
}

// Flush flushes file buffers
func (f *File) Flush(file *SDFile) error {
	if file != nil && file.ptr != nil {
		if C.pd_file_flush(file.ptr) == 0 {
			return nil
		}
	}
	return &fileError{op: "flush"}
}

// Tell returns current position in file
func (f *File) Tell(file *SDFile) (int, error) {
	if file != nil && file.ptr != nil {
		pos := C.pd_file_tell(file.ptr)
		if pos >= 0 {
			return int(pos), nil
		}
	}
	return 0, &fileError{op: "tell"}
}

// Seek seeks to position in file
func (f *File) Seek(file *SDFile, pos int, whence int) error {
	if file != nil && file.ptr != nil {
		if C.pd_file_seek(file.ptr, C.int(pos), C.int(whence)) == 0 {
			return nil
		}
	}
	return &fileError{op: "seek"}
}

// Stat returns file statistics
func (f *File) Stat(path string) (*FileStat, error) {
	// TODO: implement stat with proper struct handling
	return nil, &fileError{op: "stat", path: path}
}

// Mkdir creates a directory
func (f *File) Mkdir(path string) error {
	cpath := make([]byte, len(path)+1)
	copy(cpath, path)
	if C.pd_file_mkdir((*C.char)(unsafe.Pointer(&cpath[0]))) == 0 {
		return nil
	}
	return &fileError{op: "mkdir", path: path}
}

// Unlink deletes a file
func (f *File) Unlink(path string, recursive bool) error {
	cpath := make([]byte, len(path)+1)
	copy(cpath, path)
	var r C.int
	if recursive {
		r = 1
	}
	if C.pd_file_unlink((*C.char)(unsafe.Pointer(&cpath[0])), r) == 0 {
		return nil
	}
	return &fileError{op: "unlink", path: path}
}

// Rename renames a file
func (f *File) Rename(from, to string) error {
	cfrom := make([]byte, len(from)+1)
	copy(cfrom, from)
	cto := make([]byte, len(to)+1)
	copy(cto, to)
	if C.pd_file_rename((*C.char)(unsafe.Pointer(&cfrom[0])), (*C.char)(unsafe.Pointer(&cto[0]))) == 0 {
		return nil
	}
	return &fileError{op: "rename", path: from}
}

type fileError struct {
	op   string
	path string
}

func (e *fileError) Error() string {
	if e.path != "" {
		return "file " + e.op + " failed: " + e.path
	}
	return "file " + e.op + " failed"
}

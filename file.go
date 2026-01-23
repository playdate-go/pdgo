//go:build !tinygo

package pdgo

/*
#cgo CFLAGS: -DTARGET_EXTENSION=1
#include "pd_api.h"
#include <stdlib.h>

// File API helper functions

static const char* file_geterr(const struct playdate_file* file) {
    return file->geterr();
}

static int file_stat(const struct playdate_file* file, const char* path, FileStat* stat) {
    return file->stat(path, stat);
}

static int file_mkdir(const struct playdate_file* file, const char* path) {
    return file->mkdir(path);
}

static int file_unlink(const struct playdate_file* file, const char* name, int recursive) {
    return file->unlink(name, recursive);
}

static int file_rename(const struct playdate_file* file, const char* from, const char* to) {
    return file->rename(from, to);
}

static void* file_open(const struct playdate_file* file, const char* name, FileOptions mode) {
    return file->open(name, mode);
}

static int file_close(const struct playdate_file* file, void* f) {
    return file->close(f);
}

static int file_read(const struct playdate_file* file, void* f, void* buf, unsigned int len) {
    return file->read(f, buf, len);
}

static int file_write(const struct playdate_file* file, void* f, const void* buf, unsigned int len) {
    return file->write(f, buf, len);
}

static int file_flush(const struct playdate_file* file, void* f) {
    return file->flush(f);
}

static int file_tell(const struct playdate_file* file, void* f) {
    return file->tell(f);
}

static int file_seek(const struct playdate_file* file, void* f, int pos, int whence) {
    return file->seek(f, pos, whence);
}

// List files callback wrapper
typedef struct {
    void** paths;
    int count;
    int capacity;
} ListFilesData;

static void listFilesCallback(const char* path, void* userdata) {
    ListFilesData* data = (ListFilesData*)userdata;
    if (data->count >= data->capacity) {
        int newCapacity = data->capacity == 0 ? 16 : data->capacity * 2;
        void** newPaths = realloc(data->paths, newCapacity * sizeof(void*));
        if (newPaths == NULL) return;
        data->paths = newPaths;
        data->capacity = newCapacity;
    }
    char* copy = strdup(path);
    if (copy != NULL) {
        data->paths[data->count++] = copy;
    }
}

static int file_listfiles(const struct playdate_file* file, const char* path, void*** outPaths, int* outCount, int showhidden) {
    ListFilesData data = {NULL, 0, 0};
    int result = file->listfiles(path, listFilesCallback, &data, showhidden);
    *outPaths = data.paths;
    *outCount = data.count;
    return result;
}

static void file_freeListFiles(void** paths, int count) {
    for (int i = 0; i < count; i++) {
        free(paths[i]);
    }
    free(paths);
}
*/
import "C"
import (
	"errors"
	"io"
	"unsafe"
)

// FileOptions represents file open options
type FileOptions int

const (
	FileRead     FileOptions = C.kFileRead
	FileReadData FileOptions = C.kFileReadData
	FileWrite    FileOptions = C.kFileWrite
	FileAppend   FileOptions = C.kFileAppend
)

// SeekWhence represents seek origin
type SeekWhence int

const (
	SeekSet SeekWhence = 0 // Seek from beginning
	SeekCur SeekWhence = 1 // Seek from current position
	SeekEnd SeekWhence = 2 // Seek from end
)

// FileStat contains file information
type FileStat struct {
	IsDir     bool
	Size      uint
	ModYear   int
	ModMonth  int
	ModDay    int
	ModHour   int
	ModMinute int
	ModSecond int
}

// SDFile wraps a Playdate file handle
type SDFile struct {
	ptr  unsafe.Pointer
	file *File
}

// File wraps the playdate_file API
type File struct {
	ptr *C.struct_playdate_file
}

func newFile(ptr *C.struct_playdate_file) *File {
	return &File{ptr: ptr}
}

// GetError returns the last file error
func (f *File) GetError() string {
	err := C.file_geterr(f.ptr)
	if err == nil {
		return ""
	}
	return goString(err)
}

// Stat returns file information
func (f *File) Stat(path string) (*FileStat, error) {
	cpath := cString(path)
	defer freeCString(cpath)

	var stat C.FileStat
	result := C.file_stat(f.ptr, cpath, &stat)
	if result != 0 {
		err := f.GetError()
		if err != "" {
			return nil, errors.New(err)
		}
		return nil, errors.New("failed to stat file")
	}

	return &FileStat{
		IsDir:     stat.isdir != 0,
		Size:      uint(stat.size),
		ModYear:   int(stat.m_year),
		ModMonth:  int(stat.m_month),
		ModDay:    int(stat.m_day),
		ModHour:   int(stat.m_hour),
		ModMinute: int(stat.m_minute),
		ModSecond: int(stat.m_second),
	}, nil
}

// Mkdir creates a directory
func (f *File) Mkdir(path string) error {
	cpath := cString(path)
	defer freeCString(cpath)

	result := C.file_mkdir(f.ptr, cpath)
	if result != 0 {
		err := f.GetError()
		if err != "" {
			return errors.New(err)
		}
		return errors.New("failed to create directory")
	}
	return nil
}

// Unlink removes a file or directory
func (f *File) Unlink(name string, recursive bool) error {
	cname := cString(name)
	defer freeCString(cname)

	rec := 0
	if recursive {
		rec = 1
	}

	result := C.file_unlink(f.ptr, cname, C.int(rec))
	if result != 0 {
		err := f.GetError()
		if err != "" {
			return errors.New(err)
		}
		return errors.New("failed to unlink file")
	}
	return nil
}

// Rename renames a file or directory
func (f *File) Rename(from, to string) error {
	cfrom := cString(from)
	defer freeCString(cfrom)
	cto := cString(to)
	defer freeCString(cto)

	result := C.file_rename(f.ptr, cfrom, cto)
	if result != 0 {
		err := f.GetError()
		if err != "" {
			return errors.New(err)
		}
		return errors.New("failed to rename file")
	}
	return nil
}

// ListFiles lists files in a directory
func (f *File) ListFiles(path string, showHidden bool) ([]string, error) {
	cpath := cString(path)
	defer freeCString(cpath)

	hidden := 0
	if showHidden {
		hidden = 1
	}

	var outPaths *unsafe.Pointer
	var outCount C.int
	result := C.file_listfiles(f.ptr, cpath, &outPaths, &outCount, C.int(hidden))
	if result != 0 {
		err := f.GetError()
		if err != "" {
			return nil, errors.New(err)
		}
		return nil, errors.New("failed to list files")
	}

	if outCount == 0 {
		return []string{}, nil
	}

	defer C.file_freeListFiles(outPaths, outCount)

	paths := make([]string, int(outCount))
	cPaths := (*[1 << 20]unsafe.Pointer)(unsafe.Pointer(outPaths))[:outCount:outCount]
	for i := 0; i < int(outCount); i++ {
		paths[i] = goString((*C.char)(cPaths[i]))
	}

	return paths, nil
}

// Open opens a file
func (f *File) Open(name string, mode FileOptions) (*SDFile, error) {
	cname := cString(name)
	defer freeCString(cname)

	ptr := C.file_open(f.ptr, cname, C.FileOptions(mode))
	if ptr == nil {
		err := f.GetError()
		if err != "" {
			return nil, errors.New(err)
		}
		return nil, errors.New("failed to open file")
	}

	return &SDFile{ptr: ptr, file: f}, nil
}

// Close closes a file
func (sf *SDFile) Close() error {
	if sf.ptr == nil {
		return nil
	}
	result := C.file_close(sf.file.ptr, sf.ptr)
	sf.ptr = nil
	if result != 0 {
		err := sf.file.GetError()
		if err != "" {
			return errors.New(err)
		}
		return errors.New("failed to close file")
	}
	return nil
}

// Read reads from a file
func (sf *SDFile) Read(buf []byte) (int, error) {
	if sf.ptr == nil {
		return 0, errors.New("file is closed")
	}
	if len(buf) == 0 {
		return 0, nil
	}

	n := C.file_read(sf.file.ptr, sf.ptr, unsafe.Pointer(&buf[0]), C.uint(len(buf)))
	if n < 0 {
		err := sf.file.GetError()
		if err != "" {
			return 0, errors.New(err)
		}
		return 0, errors.New("read error")
	}
	if n == 0 {
		return 0, io.EOF
	}
	return int(n), nil
}

// Write writes to a file
func (sf *SDFile) Write(buf []byte) (int, error) {
	if sf.ptr == nil {
		return 0, errors.New("file is closed")
	}
	if len(buf) == 0 {
		return 0, nil
	}

	n := C.file_write(sf.file.ptr, sf.ptr, unsafe.Pointer(&buf[0]), C.uint(len(buf)))
	if n < 0 {
		err := sf.file.GetError()
		if err != "" {
			return 0, errors.New(err)
		}
		return 0, errors.New("write error")
	}
	return int(n), nil
}

// Flush flushes the file
func (sf *SDFile) Flush() error {
	if sf.ptr == nil {
		return errors.New("file is closed")
	}

	result := C.file_flush(sf.file.ptr, sf.ptr)
	if result != 0 {
		err := sf.file.GetError()
		if err != "" {
			return errors.New(err)
		}
		return errors.New("failed to flush file")
	}
	return nil
}

// Tell returns the current file position
func (sf *SDFile) Tell() (int, error) {
	if sf.ptr == nil {
		return 0, errors.New("file is closed")
	}

	pos := C.file_tell(sf.file.ptr, sf.ptr)
	if pos < 0 {
		err := sf.file.GetError()
		if err != "" {
			return 0, errors.New(err)
		}
		return 0, errors.New("failed to get file position")
	}
	return int(pos), nil
}

// Seek sets the file position
func (sf *SDFile) Seek(offset int, whence SeekWhence) (int, error) {
	if sf.ptr == nil {
		return 0, errors.New("file is closed")
	}

	result := C.file_seek(sf.file.ptr, sf.ptr, C.int(offset), C.int(whence))
	if result != 0 {
		err := sf.file.GetError()
		if err != "" {
			return 0, errors.New(err)
		}
		return 0, errors.New("failed to seek")
	}

	return sf.Tell()
}

// ReadAll reads the entire file
func (sf *SDFile) ReadAll() ([]byte, error) {
	var result []byte
	buf := make([]byte, 4096)
	for {
		n, err := sf.Read(buf)
		if n > 0 {
			result = append(result, buf[:n]...)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

// WriteString writes a string to the file
func (sf *SDFile) WriteString(s string) (int, error) {
	return sf.Write([]byte(s))
}

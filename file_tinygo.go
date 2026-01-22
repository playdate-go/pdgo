//go:build tinygo

// TinyGo implementation of File API

package pdgo

// SDFile represents an open file
type SDFile struct {
	ptr uintptr
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
	if bridgeFileOpen != nil {
		buf := make([]byte, len(path)+1)
		copy(buf, path)
		ptr := bridgeFileOpen(&buf[0], int32(mode))
		if ptr != 0 {
			return &SDFile{ptr: ptr}, nil
		}
	}
	return nil, &fileError{op: "open", path: path}
}

// Close closes a file
func (f *File) Close(file *SDFile) error {
	if bridgeFileClose != nil && file != nil && file.ptr != 0 {
		if bridgeFileClose(file.ptr) == 0 {
			file.ptr = 0
			return nil
		}
	}
	return &fileError{op: "close"}
}

// Read reads from a file
func (f *File) Read(file *SDFile, buf []byte) (int, error) {
	if bridgeFileRead != nil && file != nil && file.ptr != 0 && len(buf) > 0 {
		n := bridgeFileRead(file.ptr, &buf[0], uint32(len(buf)))
		if n >= 0 {
			return int(n), nil
		}
	}
	return 0, &fileError{op: "read"}
}

// Write writes to a file
func (f *File) Write(file *SDFile, buf []byte) (int, error) {
	if bridgeFileWrite != nil && file != nil && file.ptr != 0 && len(buf) > 0 {
		n := bridgeFileWrite(file.ptr, &buf[0], uint32(len(buf)))
		if n >= 0 {
			return int(n), nil
		}
	}
	return 0, &fileError{op: "write"}
}

// Flush flushes file buffers
func (f *File) Flush(file *SDFile) error {
	if bridgeFileFlush != nil && file != nil && file.ptr != 0 {
		if bridgeFileFlush(file.ptr) == 0 {
			return nil
		}
	}
	return &fileError{op: "flush"}
}

// Tell returns current position in file
func (f *File) Tell(file *SDFile) (int, error) {
	if bridgeFileTell != nil && file != nil && file.ptr != 0 {
		pos := bridgeFileTell(file.ptr)
		if pos >= 0 {
			return int(pos), nil
		}
	}
	return 0, &fileError{op: "tell"}
}

// Seek seeks to position in file
func (f *File) Seek(file *SDFile, pos int, whence int) error {
	if bridgeFileSeek != nil && file != nil && file.ptr != 0 {
		if bridgeFileSeek(file.ptr, int32(pos), int32(whence)) == 0 {
			return nil
		}
	}
	return &fileError{op: "seek"}
}

// Stat returns file statistics
func (f *File) Stat(path string) (*FileStat, error) {
	if bridgeFileStat != nil {
		buf := make([]byte, len(path)+1)
		copy(buf, path)
		var isDir, size, mtime int32
		if bridgeFileStat(&buf[0], &isDir, &size, &mtime) == 0 {
			return &FileStat{
				IsDir: isDir != 0,
				Size:  uint32(size),
				MTime: uint32(mtime),
			}, nil
		}
	}
	return nil, &fileError{op: "stat", path: path}
}

// Mkdir creates a directory
func (f *File) Mkdir(path string) error {
	if bridgeFileMkdir != nil {
		buf := make([]byte, len(path)+1)
		copy(buf, path)
		if bridgeFileMkdir(&buf[0]) == 0 {
			return nil
		}
	}
	return &fileError{op: "mkdir", path: path}
}

// Unlink deletes a file
func (f *File) Unlink(path string, recursive bool) error {
	if bridgeFileUnlink != nil {
		buf := make([]byte, len(path)+1)
		copy(buf, path)
		var r int32
		if recursive {
			r = 1
		}
		if bridgeFileUnlink(&buf[0], r) == 0 {
			return nil
		}
	}
	return &fileError{op: "unlink", path: path}
}

// Rename renames a file
func (f *File) Rename(from, to string) error {
	if bridgeFileRename != nil {
		fromBuf := make([]byte, len(from)+1)
		copy(fromBuf, from)
		toBuf := make([]byte, len(to)+1)
		copy(toBuf, to)
		if bridgeFileRename(&fromBuf[0], &toBuf[0]) == 0 {
			return nil
		}
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

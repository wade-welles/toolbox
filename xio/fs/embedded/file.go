package embedded

import (
	"bytes"
	"os"
	"path/filepath"
	"time"
)

// File holds the data for an embedded file.
type File struct {
	*bytes.Reader
	name    string
	size    int64
	modTime time.Time
	isDir   bool
	files   []os.FileInfo
	data    []byte
}

// NewFile creates a new embedded file.
func NewFile(name string, modTime time.Time, data []byte) File {
	return File{
		name:    filepath.Base(name),
		size:    int64(len(data)),
		modTime: modTime,
		data:    data,
	}
}

// Close the file. Does nothing and always returns nil. Implements the
// io.Closer interface.
func (f *File) Close() error {
	return nil
}

// Readdir reads a directory and returns information about its contents.
// Implements the http.File interface.
func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	if f.isDir {
		return f.files, nil
	}
	return nil, os.ErrNotExist
}

// Stat returns information about the file. Implements the http.File
// interface.
func (f *File) Stat() (os.FileInfo, error) {
	return f, nil
}

// Name returns the base name of the file. Implements the os.FileInfo
// interface.
func (f *File) Name() string {
	return f.name
}

// Size returns the size of the file in bytes. Implements the os.FileInfo
// interface.
func (f *File) Size() int64 {
	return f.size
}

// Mode returns the file mode bits. Implements the os.FileInfo interface.
func (f *File) Mode() os.FileMode {
	if f.isDir {
		return 0555
	}
	return 0444
}

// ModTime returns the file modification time. Implements the os.FileInfo
// interface.
func (f *File) ModTime() time.Time {
	return f.modTime
}

// IsDir returns true if this represents a directory. Implements the
// os.FileInfo interface.
func (f *File) IsDir() bool {
	return f.isDir
}

// Sys returns nil. Implements the os.FileInfo interface.
func (f *File) Sys() interface{} {
	return nil
}
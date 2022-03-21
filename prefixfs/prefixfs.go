package prefixfs

import (
	"os"

	"github.com/AndrusGerman/vfs"
)

// A FS that prefixes the path in each vfs.Filesystem operation.
type FS struct {
	fs vfs.Filesystem

	// Prefix is used to prefix the path in each vfs.Filesystem operation.
	prefix string
}

// Create returns a file system that prefixes all paths and forwards to root.
func Create(root vfs.Filesystem, prefix string) *FS {
	return &FS{fs: root, prefix: prefix}
}

// prefixPath returns path with the prefix prefixed.
func (fs *FS) prefixPath(path string) string {
	return fs.prefix + string(fs.PathSeparator()) + path
}

// PathSeparator implements vfs.Filesystem.
func (fs *FS) PathSeparator() uint8 { return fs.fs.PathSeparator() }

// OpenFile implements vfs.Filesystem.
func (fs *FS) OpenFile(name string, flag int, perm os.FileMode) (vfs.File, error) {
	return fs.fs.OpenFile(fs.prefixPath(name), flag, perm)
}

// Open implements vfs.Filesystem.
func (fs *FS) Open(name string) (vfs.File, error) {
	return fs.fs.Open(fs.prefixPath(name))
}

// Remove implements vfs.Filesystem.
func (fs *FS) Remove(name string) error {
	return fs.fs.Remove(fs.prefixPath(name))
}

// Rename implements vfs.Filesystem.
func (fs *FS) Rename(oldpath, newpath string) error {
	return fs.fs.Rename(fs.prefixPath(oldpath), fs.prefixPath(newpath))
}

// Mkdir implements vfs.Filesystem.
func (fs *FS) Mkdir(name string, perm os.FileMode) error {
	return fs.fs.Mkdir(fs.prefixPath(name), perm)
}

// Symlink implements vfs.Filesystem.
func (fs *FS) Symlink(oldname, newname string) error {
	return fs.fs.Symlink(fs.prefixPath(oldname), fs.prefixPath(newname))
}

// Stat implements vfs.Filesystem.
func (fs *FS) Stat(name string) (os.FileInfo, error) {
	return fs.fs.Stat(fs.prefixPath(name))
}

// Lstat implements vfs.Filesystem.
func (fs *FS) Lstat(name string) (os.FileInfo, error) {
	return fs.fs.Lstat(fs.prefixPath(name))
}

// ReadDir implements vfs.Filesystem.
func (fs *FS) ReadDir(path string) ([]os.FileInfo, error) {
	return fs.fs.ReadDir(fs.prefixPath(path))
}

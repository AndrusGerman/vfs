package replication

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/AndrusGerman/vfs"
)

type ReplicationFS struct {
	primary   vfs.Filesystem
	secondary []vfs.Filesystem
}

func NewReplication(primary vfs.Filesystem, secondary ...vfs.Filesystem) *ReplicationFS {
	return &ReplicationFS{
		primary:   primary,
		secondary: secondary,
	}
}

func (ctx *ReplicationFS) OpenFile(name string, flag int, perm os.FileMode) (vfs.File, error) {
	return newreplicationFileCreate(name, flag, perm, ctx.primary, ctx.secondary...)
}

func (ctx *ReplicationFS) ReadDir(path string) ([]fs.FileInfo, error) {
	return ctx.primary.ReadDir(path)
}

func (ctx *ReplicationFS) Mkdir(name string, perm os.FileMode) error {
	err := ctx.primary.Mkdir(name, perm)
	if err != nil && !errIsFileExist(err) {
		return err
	}
	for _, f := range ctx.secondary {
		err := f.Mkdir(name, perm)
		if err != nil && !errIsFileExist(err) {
			return err
		}
	}
	return nil
}

func (ctx *ReplicationFS) Rename(oldpath string, newpath string) error {
	err := ctx.primary.Rename(oldpath, newpath)
	if err != nil {
		return err
	}
	for _, f := range ctx.secondary {
		err := f.Rename(oldpath, newpath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ctx *ReplicationFS) Remove(name string) error {
	err := ctx.primary.Remove(name)
	if err != nil {
		return err
	}
	for _, f := range ctx.secondary {
		err := f.Remove(name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ctx *ReplicationFS) PathSeparator() uint8 {
	return ctx.primary.PathSeparator()
}

func (ctx *ReplicationFS) Stat(name string) (fs.FileInfo, error) {
	return ctx.primary.Stat(name)

}

func (ctx *ReplicationFS) Lstat(name string) (fs.FileInfo, error) {
	return ctx.primary.Lstat(name)
}

func (ctx *ReplicationFS) RReadDir(path string) ([]fs.FileInfo, [][]fs.FileInfo, error) {
	files, err := ctx.primary.ReadDir(path)
	if err != nil {
		return nil, nil, err
	}
	var secondary = make([][]fs.FileInfo, len(ctx.secondary))
	for i, f := range ctx.secondary {
		Sfiles, err := f.ReadDir(path)
		if err != nil {
			return nil, nil, err
		}
		secondary[i] = Sfiles
	}
	return files, secondary, nil
}

type replicationFile struct {
	primary   vfs.File
	secondary []vfs.File
}

func newreplicationFileCreate(name string, flag int, perm os.FileMode, primary vfs.Filesystem, secondary ...vfs.Filesystem) (*replicationFile, error) {
	var r = new(replicationFile)
	file, err := primary.OpenFile(name, flag, perm)
	if err != nil {
		return nil, fmt.Errorf("replication error: primary %s, name: %s", err.Error(), name)
	}
	r.primary = file
	for i, f := range secondary {
		fileS, err := f.OpenFile(name, flag, perm)
		if err != nil {
			return nil, fmt.Errorf("replication error: secondary %s, index: %d, name: %s", err.Error(), i, name)
		}
		r.secondary = append(r.secondary, fileS)
	}
	return r, err
}

func (fi replicationFile) Sync() error {
	return fi.primary.Sync()
}

func (fi replicationFile) Close() error {
	if fi.primary == nil {
		return ErrPrimaryFileIsNull
	}
	for _, f := range fi.secondary {
		if f == nil {
			continue
		}
		f.Close()
	}
	return fi.primary.Close()
}
func (fi replicationFile) Truncate(n int64) error {
	for _, f := range fi.secondary {
		f.Truncate(n)
	}
	return fi.primary.Truncate(n)
}
func (fi replicationFile) Name() string {
	return fi.primary.Name()
}
func (fi replicationFile) Write(p []byte) (n int, err error) {
	if fi.primary == nil {
		return 0, ErrPrimaryFileIsNull
	}
	for i, f := range fi.secondary {
		if f == nil {
			logSecondaryFileIsNull(fi.Name(), i)
			continue
		}
		f.Write(p)
	}
	return fi.primary.Write(p)
}

func (fi replicationFile) Read(p []byte) (n int, err error) {
	return fi.primary.Read(p)
}

func (fi replicationFile) ReadAt(p []byte, off int64) (n int, err error) {
	return fi.primary.ReadAt(p, off)
}

func (fi replicationFile) Seek(offset int64, whence int) (int64, error) {
	for _, f := range fi.secondary {
		f.Seek(offset, whence)
	}
	return fi.primary.Seek(offset, whence)
}

func errIsFileExist(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "file exists")
}

package replicationfs

import (
	"io/fs"
	"os"

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

func (ctx *ReplicationFS) Open(name string) (vfs.File, error) {
	return newreplicationFileCreate(name, os.O_RDONLY, 0, ctx.primary, ctx.secondary...)
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

func (ctx *ReplicationFS) Symlink(oldname, newname string) error {
	err := ctx.primary.Symlink(oldname, newname)
	if err != nil {
		return err
	}
	for _, f := range ctx.secondary {
		err := f.Symlink(oldname, newname)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ctx *ReplicationFS) Remove(name string) error {
	err := ctx.primary.Remove(name)
	if err != nil && !errIsNotFileExist(err) {
		return err
	}
	for _, f := range ctx.secondary {
		err := f.Remove(name)
		if err != nil && !errIsNotFileExist(err) {
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

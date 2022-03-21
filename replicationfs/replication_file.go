package replicationfs

import (
	"fmt"
	"log"
	"os"

	"github.com/AndrusGerman/vfs"
)

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

func (fi *replicationFile) Sync() error {
	if fi == nil {
		return ErrPrimaryFileIsNull
	}
	for _, f := range fi.secondary {
		if f == nil {
			continue
		}
		f.Sync()
	}
	return fi.primary.Sync()
}

func (fi *replicationFile) Close() error {
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
func (fi *replicationFile) Truncate(n int64) error {
	if fi.primary == nil {
		return ErrPrimaryFileIsNull
	}
	for _, f := range fi.secondary {
		if f == nil {
			continue
		}
		f.Truncate(n)
	}
	return fi.primary.Truncate(n)
}
func (fi *replicationFile) Name() string {
	if fi.primary == nil {
		log.Println(ErrPrimaryFileIsNull)
		return ""
	}
	return fi.primary.Name()
}
func (fi *replicationFile) Write(p []byte) (n int, err error) {
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
	//w := io.MultiWriter(buf1, buf2)
	return fi.primary.Write(p)
}

func (fi *replicationFile) Read(p []byte) (n int, err error) {
	if fi.primary == nil {
		return 0, ErrPrimaryFileIsNull
	}
	return fi.primary.Read(p)
}

func (fi *replicationFile) ReadAt(p []byte, off int64) (n int, err error) {
	if fi.primary == nil {
		return 0, ErrPrimaryFileIsNull
	}
	return fi.primary.ReadAt(p, off)
}

func (fi *replicationFile) Seek(offset int64, whence int) (int64, error) {
	if fi.primary == nil {
		return 0, ErrPrimaryFileIsNull
	}
	for i, f := range fi.secondary {
		if f == nil {
			logSecondaryFileIsNull(fi.Name(), i)
			continue
		}
		f.Seek(offset, whence)
	}
	return fi.primary.Seek(offset, whence)
}

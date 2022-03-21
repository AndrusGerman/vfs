package dumpfs

import (
	"encoding/gob"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"

	"github.com/AndrusGerman/vfs"
)

type dumpManager struct {
	fs   vfs.Filesystem
	buff io.Writer
}

func NewDumpfs(fs vfs.Filesystem, buff io.Writer) error {
	// New Manager
	var dm = &dumpManager{fs: fs, buff: buff}
	root, err := dm.createDump()
	if err != nil {
		return err
	}
	// Encode Dumpfs
	return dm.encode(root)
}

func (dm *dumpManager) createDump() (*DumpFileInfo, error) {
	fileRoot, err := dm.fs.Lstat("/")
	if err != nil {
		return nil, err
	}
	return dm.dumpFileInfo(fileRoot, "/")
}

func (dm *dumpManager) encode(root *DumpFileInfo) error {
	return gob.NewEncoder(dm.buff).Encode(root)

}
func (dm *dumpManager) dumpFileInfo(thisFile fs.FileInfo, parentDir string) (*DumpFileInfo, error) {
	if thisFile == nil {
		return nil, nil
	}
	var name = path.Join(parentDir, thisFile.Name())

	var vfi = &DumpFileInfo{
		Name: thisFile.Name(),
		Dir:  thisFile.IsDir(),
		Mode: thisFile.Mode(),
	}

	if !thisFile.IsDir() {
		file, err := dm.fs.OpenFile(name, os.O_RDONLY, 0777)
		if err != nil {
			return nil, err
		}
		bt, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		vfi.Buf = bt
	}

	if thisFile.IsDir() {
		childs, err := dm.createChildsDump(name)
		if err != nil {
			return nil, err
		}
		vfi.Childs = childs
	}

	return vfi, nil

}

func (dm *dumpManager) createChildsDump(childName string) (map[string]*DumpFileInfo, error) {
	files, err := dm.fs.ReadDir(childName)
	if err != nil {
		return nil, nil
	}
	var childs = make(map[string]*DumpFileInfo)
	for _, file := range files {
		joinName := path.Join(childName, file.Name())
		value, err := dm.dumpFileInfo(file, childName)
		if err != nil {
			return nil, err
		}
		childs[joinName] = value
	}
	return childs, nil
}

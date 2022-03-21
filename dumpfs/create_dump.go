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

func NewDumpfs(fs vfs.Filesystem, buff io.Writer) error {
	fileRoot, err := fs.Lstat("/")
	if err != nil {
		return err
	}
	root, err := DumpadoFileInfo(fs, fileRoot, "/")
	if err != nil {
		return err
	}
	return gob.NewEncoder(buff).Encode(root)
}

func DumpadoFileInfo(fs vfs.Filesystem, thisFile fs.FileInfo, parentDir string) (*DumpFileInfo, error) {
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
		file, err := fs.OpenFile(name, os.O_RDONLY, 0777)
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
		fs.ReadDir(name)
		childs, err := createChildsNewDump(fs, name)
		if err != nil {
			return nil, err
		}
		vfi.Childs = childs
	}

	return vfi, nil

}

func createChildsNewDump(mem vfs.Filesystem, childName string) (map[string]*DumpFileInfo, error) {
	files, err := mem.ReadDir(childName)
	if err != nil {
		return nil, nil
	}
	var childs = make(map[string]*DumpFileInfo)
	for _, file := range files {
		joinName := path.Join(childName, file.Name())
		value, err := DumpadoFileInfo(mem, file, childName)
		if err != nil {
			return nil, err
		}
		childs[joinName] = value
	}
	return childs, nil
}

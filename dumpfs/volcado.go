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

type DumpFileInfo struct {
	Name      string
	Dir       bool
	Mode      os.FileMode
	ParentDir string
	Childs    map[string]*DumpFileInfo
	Buf       []byte
}

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

func createChilds(mem vfs.Filesystem, childName string) (map[string]*DumpFileInfo, error) {
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
		childs, err := createChilds(fs, name)
		if err != nil {
			return nil, err
		}
		vfi.Childs = childs
	}

	return vfi, nil

}

func GetDumpfs(buff io.Reader, fs vfs.Filesystem) error {
	var data = new(DumpFileInfo)
	err := gob.NewDecoder(buff).Decode(data)
	if err != nil {
		return err
	}
	return getFilesystem(data, fs)
}

func getFilesystem(vmfs *DumpFileInfo, fs vfs.Filesystem) error {
	for _, dumpFile := range vmfs.Childs {
		err := createFileByDump(fs, dumpFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func createFileByDump(fs vfs.Filesystem, dumpFile *DumpFileInfo) error {
	var name = path.Join(dumpFile.ParentDir, dumpFile.Name)
	if dumpFile.Dir {
		err := fs.Mkdir(name, dumpFile.Mode)
		if err != nil {
			return nil
		}
	}
	if !dumpFile.Dir {
		file, err := fs.OpenFile(name, os.O_CREATE|os.O_RDWR, 0777)
		if err != nil {
			return err
		}
		defer file.Close()
		file.Write(dumpFile.Buf)
	}
	for _, dumpFile := range dumpFile.Childs {
		err := createFileByDump(fs, dumpFile)
		if err != nil {
			return nil
		}
	}
	return nil
}

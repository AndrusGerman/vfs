package volcado

import (
	"encoding/gob"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"

	"github.com/AndrusGerman/vfs"
)

type VolcadoFileInfo struct {
	Name      string
	Dir       bool
	Mode      os.FileMode
	ParentDir string
	Childs    map[string]*VolcadoFileInfo
	Buf       []byte
}

func NewVolcado(filesystem vfs.Filesystem, buff io.Writer) error {
	fileRoot, err := filesystem.Lstat("/")
	if err != nil {
		return err
	}
	root, err := newVolcadoFileInfo(filesystem, fileRoot, "/")
	if err != nil {
		return err
	}
	return gob.NewEncoder(buff).Encode(root)
}

func createChilds(mem vfs.Filesystem, fileIn fs.FileInfo, name string) (map[string]*VolcadoFileInfo, error) {
	files, err := mem.ReadDir(name)
	if err != nil {
		return nil, nil
	}
	var childs = make(map[string]*VolcadoFileInfo)
	for _, fi := range files {
		joinName := path.Join(name, fi.Name())
		value, err := newVolcadoFileInfo(mem, fi, name)
		if err != nil {
			return nil, err
		}
		childs[joinName] = value
	}
	return childs, nil
}

func newVolcadoFileInfo(mem vfs.Filesystem, fileIn fs.FileInfo, parentDir string) (*VolcadoFileInfo, error) {
	if fileIn == nil {
		return nil, nil
	}
	var name = path.Join(parentDir, fileIn.Name())

	var vfi = &VolcadoFileInfo{
		Name: fileIn.Name(),
		Dir:  fileIn.IsDir(),
		Mode: fileIn.Mode(),
	}

	if !fileIn.IsDir() {
		file, err := mem.OpenFile(name, os.O_RDONLY, 0777)
		if err != nil {
			return nil, err
		}
		bt, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		vfi.Buf = bt
	}

	if fileIn.IsDir() {
		mem.ReadDir(name)
		childs, err := createChilds(mem, fileIn, name)
		if err != nil {
			return nil, err
		}
		vfi.Childs = childs
	}

	return vfi, nil

}

func get_buf(bf *[]byte) []byte {
	if bf == nil {
		return nil
	}
	return *bf
}

func ResolveVolcado(buff io.Reader, fs vfs.Filesystem) error {
	var data = new(VolcadoFileInfo)
	err := gob.NewDecoder(buff).Decode(data)
	if err != nil {
		return err
	}
	return get_filesystem(data, fs)
}

func get_filesystem(vmfs *VolcadoFileInfo, fs vfs.Filesystem) error {
	for _, vfi := range vmfs.Childs {
		err := create_file(fs, vfi)
		if err != nil {
			return err
		}
	}

	return nil
}

func create_file(fs vfs.Filesystem, vs *VolcadoFileInfo) error {
	var name = path.Join(vs.ParentDir, vs.Name)
	if vs.Dir {
		err := fs.Mkdir(name, vs.Mode)
		if err != nil {
			return nil
		}
	}
	if !vs.Dir {
		file, err := fs.OpenFile(name, os.O_CREATE|os.O_RDWR, 0777)
		if err != nil {
			return err
		}
		defer file.Close()
		file.Write(vs.Buf)
	}
	for _, vfi := range vs.Childs {
		err := create_file(fs, vfi)
		if err != nil {
			return nil
		}
	}
	return nil
}

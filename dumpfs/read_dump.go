package dumpfs

import (
	"encoding/gob"
	"io"
	"os"
	"path"

	"github.com/AndrusGerman/vfs"
)

func GetDumpfs(buff io.Reader, fs vfs.Filesystem) error {
	var data = new(DumpFileInfo)
	err := gob.NewDecoder(buff).Decode(data)
	if err != nil {
		return err
	}
	return getFilesByDump(fs, data)
}

func getFilesByDump(dst vfs.Filesystem, src *DumpFileInfo) error {
	for _, dumpFile := range src.Childs {
		err := createFileByDump(dst, dumpFile)
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

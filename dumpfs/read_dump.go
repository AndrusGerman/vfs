package dumpfs

import (
	"encoding/gob"
	"io"
	"os"
	"path"

	"github.com/AndrusGerman/vfs"
)

type readDumpManager struct {
	fs   vfs.Filesystem
	buff io.Reader
}

func GetDumpfs(buff io.Reader, fs vfs.Filesystem) error {
	var rdm = &readDumpManager{fs: fs, buff: buff}
	dfi, err := rdm.decode()
	if err != nil {
		return err
	}
	return rdm.getFilesByDump(dfi)
}

func (rdm *readDumpManager) decode() (*DumpFileInfo, error) {
	var dfi = new(DumpFileInfo)
	err := gob.NewDecoder(rdm.buff).Decode(dfi)
	return dfi, err
}

func (rdm *readDumpManager) getFilesByDump(src *DumpFileInfo) error {
	for _, dumpFile := range src.Childs {
		err := rdm.createFileByDump(dumpFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func (rdm *readDumpManager) createFileByDump(dumpFile *DumpFileInfo) error {
	var name = path.Join(dumpFile.ParentDir, dumpFile.Name)
	if dumpFile.Dir {
		err := rdm.fs.Mkdir(name, dumpFile.Mode)
		if err != nil {
			return nil
		}
	}
	if !dumpFile.Dir {
		file, err := rdm.fs.OpenFile(name, os.O_CREATE|os.O_RDWR, 0777)
		if err != nil {
			return err
		}
		defer file.Close()
		file.Write(dumpFile.Buf)
	}
	for _, dumpFile := range dumpFile.Childs {
		err := rdm.createFileByDump(dumpFile)
		if err != nil {
			return nil
		}
	}
	return nil
}

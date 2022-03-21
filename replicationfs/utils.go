package replicationfs

import (
	"bytes"
	"io"
	"io/fs"
	"log"
	"os"
	"path"

	"github.com/AndrusGerman/vfs"
)

type UtilsSync struct {
	ReplaceDifferences bool
	DeleteNonExisting  bool
	primary            vfs.Filesystem
	secondary          []vfs.Filesystem
}

func Sync(option UtilsSync, primary vfs.Filesystem, secondary ...vfs.Filesystem) error {
	log.Println("dev note: This function is not complete")
	option.primary = primary
	option.secondary = secondary
	option.recursiveReadPrimary("/")
	return nil
}

func (us *UtilsSync) recursiveReadPrimary(folder string) error {
	files, err := us.primary.ReadDir(folder)
	if err != nil {
		return err
	}
	for _, fi := range files {
		pathJoin := path.Join(folder, fi.Name())
		us.primaryToSecondary(fi, pathJoin)
		if fi.IsDir() {
			us.recursiveReadPrimary(pathJoin)
		}

	}
	return nil
}

func (us *UtilsSync) primaryToSecondary(fi fs.FileInfo, path string) {

	// Read File
	var bufferFile *bytes.Buffer
	if !fi.IsDir() {
		file, err := us.primary.OpenFile(path, os.O_RDONLY, 0777)
		if err != nil {
			log.Panicln(err)
			return
		}
		btFile, err := io.ReadAll(file)
		if err != nil {
			log.Panicln(err)
			return
		}
		bufferFile = bytes.NewBuffer(btFile)
	}
	// Create file Function

	var createFile = func(secondaryFs vfs.Filesystem) {
		file, err := secondaryFs.OpenFile(path, os.O_CREATE|os.O_RDWR, fi.Mode().Perm())
		if err != nil {
			log.Panicln(err)
			return
		}
		defer file.Close()
		io.Copy(file, bufferFile)
	}

	// View All Secondary
	for _, secondaryFs := range us.secondary {
		_, err := secondaryFs.Stat(path)

		// Is Dir
		if os.IsNotExist(err) && fi.IsDir() {
			secondaryFs.Mkdir(path, fi.Mode().Perm())
			continue
		}

		// Is File
		if os.IsNotExist(err) && !fi.IsDir() {
			createFile(secondaryFs)
			continue
		}
		if err != nil {
			log.Panicln(err)
		}

		// This file is exist
		if !fi.IsDir() && !os.IsNotExist(err) && us.ReplaceDifferences {
			log.Println("dev note: replace file verification, by hash")
			secondaryFs.Remove(path)
			createFile(secondaryFs)
		}
	}
}

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
			us.countsFolders++
			continue
		}

		// Is File
		if os.IsNotExist(err) && !fi.IsDir() {
			createFile(secondaryFs)
			us.countsFiles++
			continue
		}
		if err != nil {
			log.Panicln(err)
		}

		// This file is exist
		if !fi.IsDir() && !os.IsNotExist(err) && us.ReplaceDifferencesFiles {
			secondaryFs.Remove(path)
			createFile(secondaryFs)
		}
	}
}

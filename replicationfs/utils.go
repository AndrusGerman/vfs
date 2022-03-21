package replicationfs

import (
	"log"

	"github.com/AndrusGerman/vfs"
)

type UtilsSync struct {
	// Delete files with different content
	ReplaceDifferencesFiles bool
	// Delete secondary files that do not exist in the primary file system
	DeleteNotExistingFiles bool
	primary                vfs.Filesystem
	secondary              []vfs.Filesystem
	countsFolders          uint
	countsFiles            uint
	removeFiles            uint
}

func Sync(option *UtilsSync, primary vfs.Filesystem, secondary ...vfs.Filesystem) error {
	if option == nil {
		option = &UtilsSync{}
	}
	log.Println("dev note: This function is not complete")
	option.primary = primary
	option.secondary = secondary

	err := option.recursiveReadSecondary("/")
	if err != nil {
		return err
	}
	err = option.recursiveReadPrimary("/")
	if err != nil {
		return err
	}
	log.Printf("dev note: add (%d) folders, add (%d) files, remove (%d) ", option.countsFolders, option.countsFiles, option.removeFiles)
	return nil
}

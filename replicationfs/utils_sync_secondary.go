package replicationfs

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/fs"
	"os"
	"path"

	"github.com/AndrusGerman/vfs"
)

// Find Files and Remove not valids
func (us *UtilsSync) recursiveReadSecondary(folder string) error {
	if !us.DeleteNotExistingFiles && !us.ReplaceDifferencesFiles {
		return nil
	}
	// Read SecondaryFS
	for _, fSecondary := range us.secondary {
		// Read Files
		err := us.recursiveReadSecondaryFiles(folder, fSecondary)
		if err != nil {
			return err
		}
	}
	return nil
}

func (us *UtilsSync) recursiveReadSecondaryFiles(folder string, fSecondary vfs.Filesystem) error {

	// Read secondary Files
	secondaryFiles, err := fSecondary.ReadDir(folder)
	if err != nil {
		return err
	}

	for _, fileSecondary := range secondaryFiles {
		pathJoin := path.Join(folder, fileSecondary.Name())

		//  Manage secondary Files
		cont, err := us.secondaryToPrimary(fSecondary, fileSecondary, pathJoin)
		if err != nil {
			return nil
		}
		// Find Internal Files
		if cont && fileSecondary.IsDir() {
			us.recursiveReadPrimary(pathJoin)
		}

	}
	return nil
}

func (us *UtilsSync) secondaryToPrimary(fSecondary vfs.Filesystem, fileSecondary fs.FileInfo, path string) (continueB bool, err error) {

	infPrimary, errPrimary := us.primary.Stat(path)

	// This File is Exists
	if !os.IsNotExist(errPrimary) {

		// Types Diferentes
		if infPrimary.IsDir() != fileSecondary.IsDir() {
			return false, vfs.RemoveAll(fSecondary, path)
		}

		// Hash Verification, is file
		if !infPrimary.IsDir() && us.ReplaceDifferencesFiles {
			// Is Not Valid
			valid := us.hashFileVerificationSame(us.primary, fSecondary, infPrimary, fileSecondary, path)
			if valid {
				return true, nil
			}
			// Is Valid
			if !valid {
				err = fSecondary.Remove(path)
				us.removeFiles++
			}
			return false, err
		}

		return true, err
	}

	// This File is not Exists
	if os.IsNotExist(errPrimary) {
		err = fSecondary.Remove(path)
		us.removeFiles++
		return false, err
	}

	return false, err
}

func (us *UtilsSync) hashFileVerificationSame(a, b vfs.Filesystem, aFile, bFile fs.FileInfo, path string) bool {
	if aFile.Size() != bFile.Size() {
		return false
	}

	aMD5, err := us.getMD5File(a, path)
	if err != nil {
		panic(err)
	}
	bMD5, err := us.getMD5File(b, path)
	if err != nil {
		panic(err)
	}
	return aMD5 == bMD5
}

func (us *UtilsSync) getMD5File(fs vfs.Filesystem, path string) (string, error) {
	hsA := md5.New()
	fileA, err := fs.OpenFile(path, os.O_RDONLY, 0777)
	if err != nil {
		panic(err)
	}
	io.Copy(hsA, fileA)
	return hex.EncodeToString(hsA.Sum(nil)), nil
}

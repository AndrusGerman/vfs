package replicationfs

import "strings"

func errIsFileExist(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "file exists")
}

func errIsNotFileExist(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "file does not exist")
}

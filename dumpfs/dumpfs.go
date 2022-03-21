package dumpfs

import (
	"os"
)

type DumpFileInfo struct {
	Name      string
	Dir       bool
	Mode      os.FileMode
	ParentDir string
	Childs    map[string]*DumpFileInfo
	Buf       []byte
}

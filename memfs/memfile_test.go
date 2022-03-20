package memfs

import (
	"testing"

	"github.com/AndrusGerman/vfs"
)

func TestFileInterface(t *testing.T) {
	_ = vfs.File(NewMemFile("", nil, nil))
}

package replication

import (
	"errors"
	"log"
)

var ErrPrimaryFileIsNull = errors.New("replication: primary file is null")

func logSecondaryFileIsNull(name string, count int) {
	log.Println("warning replication: in secondary vfs, file is null: index: ", count, " name: ", name)
}

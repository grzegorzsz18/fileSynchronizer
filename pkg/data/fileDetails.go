package data

import (
	"time"
)

type FileDetails struct {
	Name         string
	Modification time.Time
	Hash         uint32
}

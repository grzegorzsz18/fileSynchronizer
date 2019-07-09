package pkg

import (
	"fileSender/pkg/data"
	"hash/fnv"
	"os"
	"path/filepath"
)

func GetFilesList(filePath string) map[uint32]data.FileDetails {
	files := map[uint32]data.FileDetails{}
	_ = filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		files[calculateFileHash(info)] = data.FileDetails{
			Name:         info.Name(),
			Modification: info.ModTime(),
		}
		return nil
	})
	return files
}

func calculateFileHash(info os.FileInfo) uint32 {
	h := fnv.New32()
	h.Write([]byte(string(info.Size()) + info.Mode().String() + info.Name()))
	return h.Sum32()
}

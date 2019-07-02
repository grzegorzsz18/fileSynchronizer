package pkg

import (
	"fileSender/pkg/data"
	"hash/fnv"
	"os"
	"path/filepath"
)

func GetFilesList(filePath string) []data.FileDetails {
	var files []data.FileDetails
	_ = filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		files = append(files, data.FileDetails{
			Name:         info.Name(),
			Modification: info.ModTime(),
			Hash:         calculateFileHash(info),
		})
		return nil
	})
	return files
}

func calculateFileHash(info os.FileInfo) uint32 {
	h := fnv.New32()
	h.Write([]byte(string(info.Size()) + info.Mode().String() + info.Name()))
	return h.Sum32()
}

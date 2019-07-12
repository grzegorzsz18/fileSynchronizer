package pkg

import (
	"fileSender/pkg/data"
	"hash/fnv"
	"os"
	"path/filepath"
	"strings"
)

func GetFilesList(filePath string) map[uint32]data.FileDetails {

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return map[uint32]data.FileDetails{}
	}

	files := map[uint32]data.FileDetails{}
	_ = filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		pathDepth := calcPathDepth(filePath, path)
		if pathDepth == 1 {
			files[calculateFileHash(info)] = data.FileDetails{
				Name:         info.Name(),
				Modification: info.ModTime(),
				IsDirectory:  info.IsDir(),
			}
		}
		return nil
	})
	return files
}

func calcPathDepth(filePath string, path string) int {
	return len(strings.Split(path, "/")) - len(strings.Split(filePath, "/"))
}

func calculateFileHash(info os.FileInfo) uint32 {
	h := fnv.New32()
	h.Write([]byte(string(info.Size()) + info.Name()))
	return h.Sum32()
}

package files

import (
	"fileSender/client/internal"
	"fileSender/client/internal/data"
	"fileSender/pkg"
	"fmt"
	"sync"
)

func LocalFilesChecker(conf data.ClientConfig) {

	var wg sync.WaitGroup

	fileDetails, err := internal.RetrieveFilesInfoFromServer(conf)

	if err != nil {
		fmt.Printf("error while connecting to server %v", err)
		panic(1)
	}

	localFiles := pkg.GetFilesList(conf.DirectoryPath)

	for k, v := range localFiles {
		if _, ok := fileDetails[k]; !ok {
			wg.Add(1)
			go SendFileToServer(conf, v.Name, &wg)
		}
	}

	wg.Wait()
}

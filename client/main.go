package main

import (
	"fileSender/client/internal"
	"fileSender/client/internal/files"
	"fmt"
	"time"
)

func main() {

	var err error

	clientConfig, err := internal.ReadConfig()
	//
	if err != nil {
		fmt.Printf("error while reading config %v", err)
		panic(1)
	}

	for {

		files.LocalFilesChecker(clientConfig)

		fmt.Println("next iter")

		time.Sleep(time.Second * time.Duration(clientConfig.RefreshFilesTime))

	}
}

package main

import (
	"fileSender/client/internal"
	"fileSender/client/internal/data"
	"fileSender/client/internal/files"
	"fmt"
)

func main() {

	var clientConfig data.ClientConfig
	var err error
	err = internal.ReadConfig(&clientConfig)
	//
	if err != nil {
		fmt.Printf("error while reading config %v", err)
		panic(1)
	}
	//
	//fileDetails, err := internal.RetrieveFilesInfoFromServer(&clientConfig)
	//if err != nil {
	//	fmt.Printf("error while connecting to server %v", err)
	//	panic(1)
	//}
	//
	//for _ , f := range fileDetails{
	//	fmt.Println(f.Name, f.Modification, f.Hash)
	//}



	files.SendFileToServer(clientConfig, "name")
}
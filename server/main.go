package main

import (
	"fileSender/server/controller"
	"fileSender/server/databaseConnector/user"
	"fileSender/server/files"
	"sync"
)

func main() {

	var wg sync.WaitGroup
	wg.Add(2)
	user.ConnectToUserDatabase()

	go controller.ApiController()
	go files.HandleFilesReceining()

	wg.Wait()

}

//todo list
//checking user and password when sending file
//improve checking file hashes - modification instead? local.mod > remote.mod then sync

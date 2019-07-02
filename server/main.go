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

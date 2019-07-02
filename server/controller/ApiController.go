package controller

import (
	"encoding/json"
	"fileSender/pkg"
	"fileSender/pkg/data"
	"fileSender/server/databaseConnector/user"
	"fmt"
	"log"
	"net/http"
)

func ApiController() {
	http.HandleFunc("/files", handleFilesInfo)
	http.HandleFunc("/users", addNewUser)
	log.Fatal(http.ListenAndServe(":18080", nil))
}

func handleFilesInfo(rw http.ResponseWriter, req *http.Request) {

	dec := json.NewDecoder(req.Body)
	var userData = data.UserDetails{}
	_ = dec.Decode(&userData)

	userDB := user.GetUserDBConnection()
	if userDB.CheckUserCredentials(userData.Name, userData.Password) {
		files := pkg.GetFilesList(userData.Name)
		filesJson, _ := json.Marshal(files)
		fmt.Println("user got info")
		_, _ = rw.Write([]byte(filesJson))
	} else {
		fmt.Println("wrong user credentials")
		rw.WriteHeader(401)
	}

}

func addNewUser(rw http.ResponseWriter, req *http.Request) {
	dec := json.NewDecoder(req.Body)
	var u = data.UserDetails{}
	_ = dec.Decode(&u)

	userDb := user.GetUserDBConnection()
	if userDb.AddUserToDB(u.Name, u.Password) != nil {
		rw.WriteHeader(409)
	} else {
		rw.WriteHeader(201)
	}
}
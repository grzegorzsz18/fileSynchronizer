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
	http.HandleFunc("/files", files)
	http.HandleFunc("/users", users)
	log.Fatal(http.ListenAndServe(":18080", nil))
}

func files(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		dec := json.NewDecoder(req.Body)
		var userData = data.UserDetails{}
		err := dec.Decode(&userData)

		if err != nil {
			fmt.Print("error while decoding")
		}

		userDB := user.GetUserDBConnection()
		if userDB.CheckUserCredentials(userData.Name, userData.Password) {
			files := pkg.GetFilesList(userData.Name)
			err := json.NewEncoder(rw).Encode(files)
			if err != nil {
				fmt.Printf("error while encoding file list")
			}
			fmt.Println("user got info")
		} else {
			fmt.Println("wrong user credentials")
			rw.WriteHeader(401)
		}
	} else {
		rw.WriteHeader(405)
	}

}

func users(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		dec := json.NewDecoder(req.Body)
		var u = data.UserDetails{}
		_ = dec.Decode(&u)

		userDb := user.GetUserDBConnection()
		if userDb.AddUserToDB(u.Name, u.Password) != nil {
			rw.WriteHeader(409)
		} else {
			rw.WriteHeader(201)
		}
	} else {
		rw.WriteHeader(405)
	}
}

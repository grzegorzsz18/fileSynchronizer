package controller

import (
	"encoding/json"
	"fileSender/pkg"
	"fileSender/pkg/data"
	"fileSender/server/databaseConnector/user"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

func ApiController() {
	http.HandleFunc("/files", files)
	http.HandleFunc("/users", users)
	http.HandleFunc("/users/success", success)
	log.Fatal(http.ListenAndServe(":18080", nil))
}

func files(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		dec := json.NewDecoder(req.Body)
		var userData = data.UserDetails{}
		err := dec.Decode(&userData)

		if err != nil {
			fmt.Println("error while decoding")
		}

		userDB := user.GetUserDBConnection()
		if userDB.CheckUserCredentials(userData.Name, userData.Password) {
			var userPath string
			userPath = userData.Name + userData.Path
			files := pkg.GetFilesList(userPath)
			err := json.NewEncoder(rw).Encode(files)
			if err != nil {
				fmt.Println("error while encoding file list")
			}
		} else {
			fmt.Println("wrong user credentials")
			rw.WriteHeader(401)
		}
	} else {
		rw.WriteHeader(405)
	}

}

func users(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		{
			var u = data.UserDetails{}
			err := req.ParseForm()
			if err != nil {
				fmt.Printf("Error while adding new user %v \n", err)
			}

			u.Name = req.Form.Get("Name")
			u.Password = req.Form.Get("Password")

			userDb := user.GetUserDBConnection()
			if userDb.AddUserToDB(u.Name, u.Password) != nil {
				rw.WriteHeader(409)
				_, _ = rw.Write([]byte("User already exists"))
			} else {
				http.Redirect(rw, req, req.Header.Get("Referer")+"/success"+"?name="+u.Name, 301)
			}
		}
	case "GET":
		{
			wd, _ := os.Getwd()
			t := template.Must(template.ParseFiles(wd + "/static/addUser.gohtml"))
			_ = t.Execute(rw, nil)
		}
	default:
		rw.WriteHeader(405)
	}
}

func success(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		name := req.URL.Query().Get("name")

		wd, _ := os.Getwd()
		t := template.Must(template.ParseFiles(wd + "/static/success.gohtml"))
		_ = t.Execute(rw, name)
	}
}

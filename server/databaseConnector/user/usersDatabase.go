package user

import (
	"errors"
	"fileSender/pkg/data"
	"io/ioutil"
	"os"
)

type UserDB struct {
	users map[string]User
}

var userDb UserDB

func ConnectToUserDatabase() {
	users := make(map[string]User, 0)
	users["u"] = User{
		Nick:          "u",
		PasswordHash: data.EncodePassword("password"),
		EncryptionKey: data.EncodePassword("user" + "password"),
	}
	userDb = UserDB{
		users: users,
	}
}

func GetUserDBConnection() UserDB {
	return userDb
}

func (UserDB) CheckUserCredentials(userName string, userPassword string) bool {
	user, ok := userDb.users[userName]
	if !ok {
		return false
	} else {
		return userPassword == user.PasswordHash
	}
}

func (UserDB) AddUserToDB(userName string, password string) error {
	_, ok := userDb.users[userName]
	if ok {
		return errors.New("User already exists")
	} else {
		userDb.users[userName] = User{
			Nick:          userName,
			PasswordHash:  data.EncodePassword(password),
			EncryptionKey: data.EncodePassword(userName + password),
			// todo zmienic tworzenie katalogu usera
		}
		_ = os.Mkdir(userName, 777)
		_ = ioutil.WriteFile(userName+"/"+userName+".txt", []byte("ok"), 777)
	}

	return nil
}

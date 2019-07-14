package user

import (
	"database/sql"
	"errors"
	"fileSender/pkg/data"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

type UserDB struct {
	connection *sql.DB
}

var userDb UserDB

func ConnectToUserDatabase() {

}

func GetUserDBConnection() UserDB {
	return userDb
}

func (UserDB) CheckUserCredentials(userName string, userPassword string) bool {
	db, err := sql.Open("sqlite3", "./users.db")

	checkErr(err)

	defer db.Close()

	stmt, err := db.Query("SELECT * FROM 'users' WHERE Nick=? AND PasswordHash=?", userName, userPassword)

	checkErr(err)

	return stmt.Next()
}

func (UserDB) AddUserToDB(userName string, password string) error {

	db, err := sql.Open("sqlite3", "./users.db")

	checkErr(err)

	defer db.Close()

	stmt, err := db.Query("SELECT * FROM 'users' WHERE Nick=?", userName)

	checkErr(err)

	if stmt.Next() {
		return errors.New("User already exists")
	} else {
		_, err = db.Exec("INSERT into 'users' values(?,?)", userName, data.EncodePassword(password))

		checkErr(err)

		_ = os.Mkdir(userName, 0777)
	}

	return nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

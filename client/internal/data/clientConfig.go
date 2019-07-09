package data

import "fileSender/pkg/data"

type ClientConfig struct {
	UserName          string
	UserPassword      string
	userPasswordHash  string
	ServerHost        string
	ServerPortRest    string
	ServerPortTCP     string
	RefreshFilesTime  int
	userEncriptionKey string
	DirectoryPath     string
}

func (c ClientConfig) GetUserPasswordHash() string {
	return data.EncodePassword(c.UserPassword)
}

func (c ClientConfig) GetUserEncryptionKey() string {
	return data.EncodePassword(c.UserName + c.UserPassword)
}

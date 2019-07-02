package internal

import (
	"fileSender/client/internal/data"
	data2 "fileSender/pkg/data"
)

func ReadConfig(config *data.ClientConfig) error {

	config.RefreshFilesTime = 1
	config.ServerHost = "localhost"
	config.ServerPortRest = "18080"
	config.UserName = "user"
	config.UserPasswordHash = data2.EncodePassword("password")
	config.UserEncriptionKey = data2.EncodePassword("." + "password")
	config.ServerPortTCP = "22222"

	return nil
}

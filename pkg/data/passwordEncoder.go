package data

import (
	"crypto/md5"
	"encoding/hex"
)

func EncodePassword(password string) string {
	md5Hash := md5.New()
	md5Hash.Write([]byte(password))
	return hex.EncodeToString(md5Hash.Sum(nil))
}

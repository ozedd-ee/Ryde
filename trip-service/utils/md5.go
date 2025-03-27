package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func GenerateMD5Key(str1, str2 string) string {
	data := str1 + ":" + str2
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

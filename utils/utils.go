package utils

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

func MdHashing(input string) string {
	byteInput := []byte(input)
	md5Hash := md5.Sum(byteInput)
	return hex.EncodeToString(md5Hash[:])
}

func Expiration() time.Time {
	return time.Now().Add(24 * time.Hour)
}

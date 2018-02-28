package utils

import (
	"math/rand"
	"time"
)

const (
	PASSWORDDIC = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func RandomString(length int) string {
	rand.Seed(time.Now().Unix())
	stringarray := make([]byte, length)
	diclen := len(PASSWORDDIC)
	for i := 0; i < length; i++ {
		stringarray[i] = PASSWORDDIC[rand.Int()%diclen]
	}

	return string(stringarray)
}

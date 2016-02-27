package uuid

import (
	"math/rand"
	"time"
)

var (
	alphaNum []byte = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	size     int    = 10
)

func NewUuid() []byte {
	uid := make([]byte, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		uid[i] = alphaNum[rand.Intn(len(alphaNum))]
	}
	return uid
}

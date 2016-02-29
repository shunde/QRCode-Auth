package uuid

import (
	"math/rand"
)

const (
	alphaNum string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	size     int    = 10
)

func NewUuid() []byte {
	return randStr(size)
}

func NewUuidN(n int) []byte {
	if n <= 0 {
		return NewUuid()
	}
	return randStr(n)
}

func randStr(n int) []byte {
	uid := make([]byte, n)

	for i := 0; i < n; i++ {
		uid[i] = alphaNum[rand.Intn(len(alphaNum))]
	}
	return uid
}

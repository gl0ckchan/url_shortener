package random

import (
	"math/rand"
	"time"
)

var charset = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func NewRandomString(size int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	buffer 	   := make([]rune, size)

	for i := range buffer {
		buffer[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(buffer)
}
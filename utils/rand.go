package utils

import (
	"math/rand"
	"os"
	"strconv"
	"time"
)

//RandStringGen - Random string genrator for url
func RandStringGen() string {
	length, _ := strconv.Atoi(os.Getenv("URL_LENGTH")) //Length of randomnized url
	charset := "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" +
		"0123456789"

	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

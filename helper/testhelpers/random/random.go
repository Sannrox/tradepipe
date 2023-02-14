package random

import (
	"math/rand"
	"time"
)

func GenerateRandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func GenerateRandomYear() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(30) + 1990
}

func GenerateRandomMonth() time.Month {
	rand.Seed(time.Now().UnixNano())
	return time.Month(rand.Intn(12) + 1)
}

package random

import (
	"math/rand"
	"time"
)

var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateRandomString(n int) string {

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

func GenerateRandomeISIN() string {
	// 1 2  Countrie-code
	b := make([]rune, 12)
	b[0] = letter[rand.Intn(len(letter))]
	b[1] = letter[rand.Intn(len(letter))]
	// 3 11 National Securities Identifying Number
	for i := 2; i < 10; i++ {
		b[i] = letter[rand.Intn(len(letter))]
	}
	// 12 12 Check Digit
	b[11] = letter[rand.Intn(len(letter))]

	return string(b)
}

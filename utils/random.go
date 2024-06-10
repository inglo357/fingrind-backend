package utils

import "math/rand"

var alphabets string = "abcdefghijklmnopqrstuvwxyz1234567890"

func RandomString(r int) string{
	bits := []rune{}
	k := len(alphabets)

	for i := 0; i < r; i++{
		index := rand.Intn(k)
		bits = append(bits, rune(alphabets[index]))
	}

	return string(bits)
}

func RandomEmail() string{
	return RandomString(10) + "@asd.com"
}
package util

import "math/rand"

func GenarateRandomNumber(n int) string {
	var charsets = []rune("0123456789")
	numbers := make([]rune, n)
	for i := range numbers {
		numbers[i] = charsets[rand.Intn(len(charsets))]
	}

	return string(numbers)
}

package helper

import "math/rand"

func GenerateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, n)
	for i := 0; i < n; i++ {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

package cmd

import (
	"math/rand"
)

const letters = "abcdefghijklmnopqrstuvwxyz"

func randLowerCaseString(n int) string {
	b := make([]byte, n)
	l := int64(len(letters))
	for i := range b {
		b[i] = letters[rand.Int63()%l] //nolint:gosec
	}
	return string(b)
}

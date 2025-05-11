package util

import (
	"math/rand/v2"
	"strings"
)

const (
	printableStart = 32
	printableEnd   = 126
	alphabet       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alphanum       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func RandomPrintableString(n int) string {
	var sb strings.Builder

	for i := 0; i < n; i++ {
		c := byte(rand.IntN(printableEnd-printableStart+1) + printableStart)
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomAlphabetString(n int) string {
	var sb strings.Builder

	for i := 0; i < n; i++ {
		c := alphabet[rand.IntN(len(alphabet))]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomAlphaNumString(n int) string {
	var sb strings.Builder

	for i := 0; i < n; i++ {
		c := alphanum[rand.IntN(len(alphanum))]
		sb.WriteByte(c)
	}

	return sb.String()
}

package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomPrintableString(t *testing.T) {
	s := RandomPrintableString(100)
	ok := true

	for i := range s {
		if s[i] < printableStart || s[i] > printableEnd {
			ok = false
			break
		}
	}

	require.True(t, ok)
}

func TestRandomAlphabetString(t *testing.T) {
	s := RandomAlphabetString(100)
	ok := true

	for i := range s {
		isInAlphabetString := false
		for j := range alphabet {
			if s[i] == alphabet[j] {
				isInAlphabetString = true
				break
			}
		}
		if !isInAlphabetString {
			ok = false
			break
		}
	}

	require.True(t, ok)
}

func TestRandomAlphanumString(t *testing.T) {
	s := RandomAlphabetString(100)
	ok := true

	for i := range s {
		isInAlphabetString := false
		for j := range alphanum {
			if s[i] == alphanum[j] {
				isInAlphabetString = true
				break
			}
		}
		if !isInAlphabetString {
			ok = false
			break
		}
	}

	require.True(t, ok)
}

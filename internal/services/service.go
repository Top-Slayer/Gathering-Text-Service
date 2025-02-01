package services

import (
	"unicode"
)

func IsLaoText(text string) bool {
	for _, r := range text {
		if r != ' ' && !unicode.Is(unicode.Lao, r) {
			return false
		}
	}
	return true

}

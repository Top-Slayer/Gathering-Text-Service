package services

import (
	"strings"
	"unicode"
)

func ClearSpace(text *string) {
	*text = strings.TrimSpace(*text)
}

func IsLaoText(text string) bool {
	for _, r := range text {
		if r != ' ' && !unicode.Is(unicode.Lao, r) {
			return false
		}
	}
	return true
}

func CheckLaoFormat(text string) bool {
	runes := []rune(text)
	var t rune

	if (runes[0] >= 0x0eb0 && runes[0] <= 0x0ebd) ||
		(runes[0] >= 0x0ec8 && runes[0] <= 0x0ecd) ||
		len(runes) < 3 {
		return false
	}

	for _, c := range runes {
		if c == t {
			return false
		}
		t = c
	}

	return true
}

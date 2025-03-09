package services

import (
	"crypto/sha256"
	"encoding/hex"
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
	var t [2]rune

	if len(runes) == 0 {
		return false
	}

	if (runes[0] >= 0x0eb0 && runes[0] <= 0x0ebd) ||
		(runes[0] >= 0x0ec8 && runes[0] <= 0x0ecd) ||
		len(runes) < 2 {
		return false
	}

	for i, c := range runes {
		if c == t[0] && t[0] == t[1] {
			return false
		}
		t[0] = runes[i]
		if i > 0 {
			t_i := i - 1
			t[1] = runes[t_i]
		}
	}

	return true
}

func AutorizeAdmin(pass []byte) bool {
	hash := sha256.New()
	hash.Write(pass)
	if hex.EncodeToString(hash.Sum(nil)) == "f766dad97841c5b14ab7e88f4f9c60e94b251b37eaefddc94251860adf75cfd9" {
		return true
	} else {
		return false
	}
}

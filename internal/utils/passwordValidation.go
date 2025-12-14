package utils

import (
	"fmt"
	"strings"
	"unicode"
)

type Messages struct {
	Length  string
	Upper   string
	Lower   string
	Number  string
	Special string
}

func ValidatePassword(plaintext string) (bool, error) {
	password := strings.TrimSpace(plaintext)
	const minLength int = 8
	const maxLength int = 32
	var hasUpper bool
	var hasLower bool
	var hasDigit bool
	var hasSpecial bool
	var isValidPassword bool = false
	var errorString string

	msg := Messages{
		Length:  "password must be at least 8 to 32 characters long",
		Upper:   "it must have at least 1 uppercase character",
		Lower:   "it must have at least 1 lowercase character",
		Number:  "it must have at least 1 digit from 0 to 9",
		Special: "it must have at least 1 special character from @!£$%^&*_-+",
	}

	s := []rune(password)
	passLength := len(s)
	if !(passLength >= minLength && passLength <= maxLength) {
		errorString += msg.Length + ", also,\n\t- " + msg.Lower + "\n\t- " + msg.Upper + "\n\t- " + msg.Number + "\n\t- " + msg.Special
		return isValidPassword, fmt.Errorf("%v", errorString)

	}

	for _, ch := range s {
		switch {
		case unicode.IsLetter(ch) && unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLetter(ch) && unicode.IsLower(ch):
			hasLower = true
		case unicode.IsNumber(ch):
			hasDigit = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}

	}

	if !hasLower {
		errorString += ", " + msg.Lower
	}

	if !hasUpper {
		errorString += ", " + msg.Upper
	}

	if !hasDigit {
		errorString += ", " + msg.Number
	}
	if !hasSpecial {
		errorString += ", " + msg.Special
	}

	if len(errorString) != 0 {
		return isValidPassword, fmt.Errorf("final:%v", errorString)
	}

	isValidPassword = hasUpper && hasLower && hasDigit && hasSpecial
	return isValidPassword, nil
}

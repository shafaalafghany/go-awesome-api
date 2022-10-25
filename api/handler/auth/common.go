package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"unicode"
)

func ValidateEmail(email string) error {
	if strings.TrimSpace(email) == "" {
		return fmt.Errorf("email cannot be empty")
	}
	if !strings.Contains(email, "@") {
		return errors.New("email is missing @")
	}
	if !strings.Contains(email, ".") {
		return errors.New("email is missing . ")
	}
	atPos := strings.IndexByte(email, byte('@'))
	dotPos := strings.LastIndexByte(email, byte('.'))
	if atPos > dotPos {
		return fmt.Errorf("'@' must come before last '.'")
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) > 40 {
		return fmt.Errorf("password cannot exceed 40 characters")
	}
	if len(password) < 8 {
		return fmt.Errorf("password must be more than 8 characters")
	}

	hasUpperCase := false
	hasLowerCase := false
	hasNumber := false
	hasSpecialCharacter := false

	for _, v := range password {
		if unicode.IsNumber(v) {
			hasNumber = true
		}
		if unicode.IsLower(v) {
			hasLowerCase = true
		}
		if unicode.IsUpper(v) {
			hasUpperCase = true
		}
		if unicode.IsPunct(v) {
			hasSpecialCharacter = true
		}
	}

	if hasUpperCase && hasLowerCase && hasNumber && hasSpecialCharacter {
		return nil
	}
	if !hasLowerCase {
		return fmt.Errorf("must have a lowercase character")
	}
	if !hasUpperCase {
		return fmt.Errorf("must have a uppercase character")
	}
	if !hasNumber {
		return fmt.Errorf("must have a numerical character")
	}
	if !hasSpecialCharacter {
		return fmt.Errorf("must have a special character")
	}
	return nil
}

func ValidateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if len(name) > 128 {
		return fmt.Errorf("name cannot exceed 128 character")
	}

	return nil
}

func RandString(n int) (string, error) {
	const letter = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	generatedString := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letter))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random string")
		}
		generatedString[i] = letter[num.Int64()]
	}
	return string(generatedString), nil
}

package models

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

func ValidatePassword(password string) error {
	if len(password) < MinPasswordLength || len(password) > MaxPasswordLength {
		return fmt.Errorf(ErrPasswordLength, MinPasswordLength, MaxPasswordLength)
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return errors.New(ErrPasswordComplexity)
	}

	return nil
}

func Validate(u UserSignupRequest) error {
	if u.Name == "" || u.Email == "" || u.Password == "" {
		return errors.New( ErrRequiredFieldsEmpty)
	}

	if err := ValidatePassword(u.Password); err != nil {
		return err
	}

	if !strings.Contains(u.Email, "@") || !strings.Contains(u.Email, ".") {
		return errors.New(ErrInvalidEmailFormat)
	}

	if u.Age <= 0 {
		return errors.New(ErrNegativeAge)
	}

	return nil
}

package util

import (
	"fmt"
	"net/mail"
)

const (
	StudentRole = "student"
	AdminRole   = "admin"
)

func ValidateString(value string, minLength, maxLength int) error {
	if len(value) < minLength || len(value) > maxLength {
		return fmt.Errorf("string length should be between %d and %d", minLength, maxLength)
	}
	return nil
}

func ValidatePassword(value string) error {
	if err := ValidateString(value, 6, 50); err != nil {
		return err
	}
	return nil
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 200); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("email is not valid")
	}
	return nil
}

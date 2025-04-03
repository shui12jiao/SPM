package util

import (
	"fmt"
	"net/mail"
	"time"
)

const (
	StudentRole = "student"
	AdminRole   = "admin"
)

// start_time
func ValidateStartTime(value time.Time) error {
	// 要求开始时间在当前时间之后
	if value.Before(time.Now()) {
		return fmt.Errorf("开始时间必须在当前时间之后")
	}
	return nil
}

// password
func ValidatePassword(value string) error {
	if err := ValidateString(value, 6, 50); err != nil {
		return err
	}
	return nil
}

// email
func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 200); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("无效的电子邮件地址: %s", value)
	}
	return nil
}

func ValidateString(value string, minLength, maxLength int) error {
	if len(value) < minLength || len(value) > maxLength {
		return fmt.Errorf("字符串长度必须在 %d 到 %d 之间", minLength, maxLength)
	}
	return nil
}

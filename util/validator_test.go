package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateEmail(t *testing.T) {
	testCases := []struct {
		name  string
		email string
		valid bool
	}{
		{
			name:  "有效的邮箱",
			email: "test@example.com",
			valid: true,
		},
		{
			name:  "缺少域名",
			email: "test@",
			valid: false,
		},
		{
			name:  "缺少@符号",
			email: "testexample.com",
			valid: false,
		},
		{
			name:  "只有用户名",
			email: "test",
			valid: false,
		},
		{
			name:  "包含中文的邮箱",
			email: "测试@example.com",
			valid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateEmail(tc.email)
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	testCases := []struct {
		name     string
		password string
		valid    bool
	}{
		{
			name:     "有效的密码 - 6个字符",
			password: "123456",
			valid:    true,
		},
		{
			name:     "有效的密码 - 12个字符",
			password: "123456789012",
			valid:    true,
		},
		{
			name:     "密码太短",
			password: "12345",
			valid:    false,
		},
		{
			name:     "密码太长",
			password: "1234567890123456789012345678901234567890123456789012345",
			valid:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePassword(tc.password)
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

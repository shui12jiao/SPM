package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err, "生成密码哈希失败")
	require.NotEmpty(t, hashedPassword, "密码哈希为空")

	err = CheckPassword(password, hashedPassword)
	require.NoError(t, err, "密码验证失败")

	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, hashedPassword)
	require.Equal(t, err, bcrypt.ErrMismatchedHashAndPassword, "密码验证错误")

	hashedPasswordOther, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPasswordOther)
	require.NotEqual(t, hashedPassword, hashedPasswordOther, "不同时间生成的密码哈希不应相同")
}

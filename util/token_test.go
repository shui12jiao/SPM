package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(RandomString(32))
	require.NoError(t, err)

	userID := 1
	role := StudentRole
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(userID, role, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	require.Equal(t, userID, payload.UserID)
	require.Equal(t, role, payload.UserRole)
	require.NotZero(t, payload.ID)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(RandomString(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(1, StudentRole, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, "token已过期")
	require.Nil(t, payload)
}

func TestInvalidJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(RandomString(32))
	require.NoError(t, err)

	payload, err := maker.VerifyToken(RandomString(20))
	require.Error(t, err)
	require.Contains(t, err.Error(), "token无效")
	require.Nil(t, payload)
}

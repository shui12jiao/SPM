package util

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrExpiredToken = errors.New("token已过期")
	ErrInvalidToken = errors.New("token无效")
)

// token管理器
type Maker interface {
	CreateToken(userID int, userRole string, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}

type Payload struct {
	ID        uuid.UUID `json:"id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`

	UserID   int    `json:"user_id"`
	UserRole string `json:"user_role"`
}

func NewPayload(userID int, userRole string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),

		UserID:   userID,
		UserRole: userRole,
	}
	return payload, err
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

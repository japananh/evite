package tokenprovider

import (
	"app-invite-service/common"
	"errors"
	"time"
)

var (
	ErrNotFound = common.NewCustomError(
		errors.New("token not found"),
		"token not found",
		"ErrNotFound",
	)
	ErrInvalidToken = common.NewCustomError(
		errors.New("invalid token provided"),
		"invalid token provided",
		"ErrInvalidToken",
	)
)

type Token struct {
	Token   string    `json:"token"`
	Created time.Time `json:"created"`
	Expiry  int       `json:"expiry"` // milliseconds
}

type TokenPayload struct {
	UserId          int    `json:"user_id,omitempty"`
	InvitationToken string `json:"invite_token,omitempty"`
}

type TokenConfig struct {
	AccessTokenExpiry  int
	RefreshTokenExpiry int
}

func NewTokenConfig(atExpiry, rtExpiry int) (*TokenConfig, error) {
	return &TokenConfig{
		AccessTokenExpiry:  atExpiry,
		RefreshTokenExpiry: rtExpiry,
	}, nil
}

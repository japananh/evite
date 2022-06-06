package userbiz

import (
	"app-invite-service/common"
	"app-invite-service/component/tokenprovider"
	"app-invite-service/modules/user/usermodel"
	"context"
	crand "crypto/rand"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"io"
	"math/big"
	"math/rand"
	"strings"
	"time"
)

var (
	ErrInviteTokenNotExisted = common.NewCustomError(
		errors.New("invite token not existed"),
		"invite token not existed",
		"ErrInviteTokenNotExisted",
	)
	ErrInvalidInviteToken = common.NewCustomError(
		errors.New("invalid invite token"),
		"invalid invite token",
		"ErrInvalidInviteToken",
	)
)

// Adapted from https://elithrar.github.io/article/generating-secure-random-numbers-crypto-rand/
func init() {
	assertAvailablePRNG()
}

func assertAvailablePRNG() {
	// Assert that a cryptographically secure PRNG is available.
	// Panic otherwise.
	buf := make([]byte, 1)
	_, err := io.ReadFull(crand.Reader, buf)
	if err != nil {
		panic(fmt.Sprintf("crypto/rand is unavailable: Read() failed with %#v", err))
	}
}

// GenerateRandomString returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(min, max int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(max-min) + min
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := crand.Int(crand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}
	return string(ret), nil
}

type IRedis interface {
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd
}

type generateTokenBiz struct {
	redis IRedis
}

func NewGenerateTokenBiz(redis IRedis) *generateTokenBiz {
	return &generateTokenBiz{redis: redis}
}

func (biz *generateTokenBiz) GenerateToken(ctx context.Context) (*usermodel.InvitationToken, error) {
	var minTokenLen = 6
	var maxTokenLen = 12
	token, err := GenerateRandomString(minTokenLen, maxTokenLen)
	if err != nil {
		return nil, err
	}

	val, err := biz.redis.SetNX(ctx, token, token, common.InviteTokenExpirySecond*time.Second).Result()
	if err != nil || !val {
		return nil, err
	}

	return &usermodel.InvitationToken{Token: token}, nil
}

type LoginWithInviteTokenRedis interface {
	MGet(ctx context.Context, keys ...string) *redis.SliceCmd
}

type loginWithInviteTokenBiz struct {
	redis         LoginWithInviteTokenRedis
	tokenProvider tokenprovider.Provider
	hash          Hash
	tokenConfig   *tokenprovider.TokenConfig
}

func NewLoginWithInviteTokenBiz(
	redis LoginWithInviteTokenRedis,
	tokenProvider tokenprovider.Provider,
	hash Hash,
	tokenConfig *tokenprovider.TokenConfig,
) *loginWithInviteTokenBiz {
	return &loginWithInviteTokenBiz{
		redis:         redis,
		tokenProvider: tokenProvider,
		hash:          hash,
		tokenConfig:   tokenConfig,
	}
}

func (biz *loginWithInviteTokenBiz) LoginWithInviteToken(ctx context.Context, data *usermodel.UserLoginWithInviteToken) (*usermodel.Account, error) {
	if err := data.Validate(); err != nil {
		return nil, err
	}

	// validate redis token
	if result := biz.redis.MGet(ctx, data.InvitationToken); result == nil {
		return nil, ErrInviteTokenNotExisted
	}

	payload := tokenprovider.TokenPayload{
		InvitationToken: data.InvitationToken,
	}

	accessToken, err := biz.tokenProvider.Generate(payload, biz.tokenConfig.AccessTokenExpiry)
	if err != nil {
		return nil, common.ErrInternal(err)
	}

	refreshToken, err := biz.tokenProvider.Generate(payload, biz.tokenConfig.RefreshTokenExpiry)
	if err != nil {
		return nil, common.ErrInternal(err)
	}

	account := usermodel.NewAccount(accessToken, refreshToken)

	return account, nil
}

type ValidateInviteTokenRedis interface {
	MGet(ctx context.Context, keys ...string) *redis.SliceCmd
}

type validateInviteTokenBiz struct {
	redis ValidateInviteTokenRedis
}

func NewValidateInviteTokenBiz(redis ValidateInviteTokenRedis) *validateInviteTokenBiz {
	return &validateInviteTokenBiz{redis: redis}
}

func (biz *validateInviteTokenBiz) ValidateInviteToken(
	ctx context.Context,
	token usermodel.InvitationToken,
) error {
	if result := biz.redis.MGet(ctx, strings.TrimSpace(token.Token)); result == nil {
		return ErrInvalidInviteToken
	}

	return nil
}

package userbiz

import (
	"app-invite-service/common"
	"app-invite-service/component/tokenprovider"
	usermodel "app-invite-service/modules/user/usermodel"
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

// generate invitation token

type generateTokenBiz struct {
	redis *redis.Client
}

func NewGenerateTokenBiz(redis *redis.Client) *generateTokenBiz {
	return &generateTokenBiz{redis: redis}
}

func (biz *generateTokenBiz) GenerateToken(ctx context.Context) (*usermodel.InvitationToken, error) {
	var minTokenLen = 6
	var maxTokenLen = 12
	token, err := GenerateRandomString(minTokenLen, maxTokenLen)
	if err != nil {
		return nil, err
	}

	payload := usermodel.InvitationToken{Token: token, Status: 1}

	p, err := payload.MarshalBinary()
	if err != nil {
		return nil, err
	}

	val, err := biz.redis.SetNX(ctx, token, string(p), common.InviteTokenExpirySecond*time.Second).Result()
	if err != nil || !val {
		return nil, err
	}

	return &payload, nil
}

// Login with invitation token

type loginWithInviteTokenBiz struct {
	redis         *redis.Client
	tokenProvider tokenprovider.Provider
	hash          Hash
	tokenConfig   *tokenprovider.TokenConfig
}

func NewLoginWithInviteTokenBiz(
	redis *redis.Client,
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

	// check redis token existed
	tokenFromRedis := biz.redis.Get(ctx, data.InvitationToken)
	if tokenFromRedis.Val() == "" {
		return nil, ErrInviteTokenNotExisted
	}

	// check whether invitation token is disabled or not
	var foundToken usermodel.InvitationToken
	if err := foundToken.UnmarshalBinary([]byte(tokenFromRedis.Val())); err != nil {
		return nil, common.ErrInternal(err)
	}
	if foundToken.Status == 0 {
		return nil, ErrInvalidInviteToken
	}

	// create JWT token
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

// Validate invitation token

type validateInviteTokenBiz struct {
	redis *redis.Client
}

func NewValidateInviteTokenBiz(redis *redis.Client) *validateInviteTokenBiz {
	return &validateInviteTokenBiz{redis: redis}
}

func (biz *validateInviteTokenBiz) ValidateInvitationToken(
	ctx context.Context,
	token string,
) error {
	// check token existed
	tokenFromRedis := biz.redis.Get(ctx, strings.TrimSpace(token))
	if tokenFromRedis.Val() == "" {
		return ErrInviteTokenNotExisted
	}

	// check whether token disabled or not
	var foundToken usermodel.InvitationToken
	if err := foundToken.UnmarshalBinary([]byte(tokenFromRedis.Val())); err != nil {
		return common.ErrInternal(err)
	}
	if foundToken.Status == 0 {
		return ErrInvalidInviteToken
	}

	return nil
}

// list all invitation token

type listInvitationTokenBiz struct {
	redis *redis.Client
}

func NewListInvitationTokenBiz(redis *redis.Client) *listInvitationTokenBiz {
	return &listInvitationTokenBiz{redis: redis}
}

func (biz *listInvitationTokenBiz) ListInvitationToken(
	ctx context.Context,
	filter *usermodel.InvitationTokenFilter,
	_ *common.Paging,
) ([]usermodel.InvitationToken, error) {
	var listToken []usermodel.InvitationToken

	iter := biz.redis.Scan(ctx, 0, "prefix:*", 0).Iterator()
	for iter.Next(ctx) {
		token := &usermodel.InvitationToken{}
		if err := token.UnmarshalBinary([]byte(iter.Val())); err != nil {
			break
		}
		listToken = append(listToken, *token)
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	var filteredListToken []usermodel.InvitationToken
	for _, token := range listToken {
		if *filter.Status == token.Status {
			filteredListToken = append(filteredListToken, token)
		}
	}

	return filteredListToken, nil
}

// Update invitation token

type updateInvitationTokenBiz struct {
	redis *redis.Client
}

func NewUpdateInvitationTokenBiz(redis *redis.Client) *updateInvitationTokenBiz {
	return &updateInvitationTokenBiz{redis: redis}
}

func (biz *updateInvitationTokenBiz) UpdateInvitationToken(
	ctx context.Context,
	token string,
	data *usermodel.InvitationTokenUpdate,
) error {
	// check redis token existed
	tokenFromRedis := biz.redis.Get(ctx, strings.TrimSpace(token))
	if tokenFromRedis.Val() == "" {
		return ErrInviteTokenNotExisted
	}

	var foundToken usermodel.InvitationToken
	if err := foundToken.UnmarshalBinary([]byte(tokenFromRedis.Val())); err != nil {
		return common.ErrInternal(err)
	}

	// update token
	foundToken.Status = data.Status
	t, err := foundToken.MarshalBinary()
	if err != nil {
		return common.ErrInternal(err)
	}

	if _, err := biz.redis.SetXX(ctx, token, t, redis.KeepTTL).Result(); err != nil {
		return common.ErrInternal(err)
	}

	return nil
}

package component

import (
	"app-invite-service/component/tokenprovider"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type AppContext interface {
	SecretKey() string
	GetMainDBConnection() *gorm.DB
	GetRedisConnection() *redis.Client
	GetTokenConfig() *tokenprovider.TokenConfig
}

type appCtx struct {
	secretKey   string
	db          *gorm.DB
	redis       *redis.Client
	tokenConfig *tokenprovider.TokenConfig
}

func NewAppContext(
	db *gorm.DB,
	redis *redis.Client,
	secretKey string,
	tokenConfig *tokenprovider.TokenConfig,
) *appCtx {
	return &appCtx{secretKey: secretKey, db: db, redis: redis, tokenConfig: tokenConfig}
}

func (ctx *appCtx) GetMainDBConnection() *gorm.DB {
	return ctx.db
}

func (ctx *appCtx) GetRedisConnection() *redis.Client {
	return ctx.redis
}

func (ctx *appCtx) SecretKey() string {
	return ctx.secretKey
}

func (ctx *appCtx) GetTokenConfig() *tokenprovider.TokenConfig {
	return ctx.tokenConfig
}

package component

import (
	"app-invite-service/component/tokenprovider"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type AppContext interface {
	SecretKey() string
	GetDBConn() *gorm.DB
	GetRedisConn() *redis.Client
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
) AppContext {
	return &appCtx{secretKey: secretKey, db: db, redis: redis, tokenConfig: tokenConfig}
}

func (ctx *appCtx) GetDBConn() *gorm.DB {
	return ctx.db
}

func (ctx *appCtx) GetRedisConn() *redis.Client {
	return ctx.redis
}

func (ctx *appCtx) SecretKey() string {
	return ctx.secretKey
}

func (ctx *appCtx) GetTokenConfig() *tokenprovider.TokenConfig {
	return ctx.tokenConfig
}

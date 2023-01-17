package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App    `yaml:"app"`
		Logger `yaml:"logger"`
		MySQL  `yaml:"mysql"`
		Redis  `yaml:"redis"`
		//RMQ   `yaml:"rabbitmq"`
	}

	App struct {
		AtExpiry       int    `env-required:"true" yaml:"access_token_expiry"  env:"ACCESS_TOKEN_EXPIRY"`
		RtExpiry       int    `env-required:"true" yaml:"refresh_token_expiry" env:"REFRESH_TOKEN_EXPIRY"`
		Port           int    `env-required:"true" yaml:"port"                 env:"APP_PORT"`
		Name           string `env-required:"true" yaml:"name"                 env:"APP_NAME"`
		Version        string `env-required:"true" yaml:"version"              env:"APP_VERSION"`
		ENV            string `env-required:"true" yaml:"env"                  env:"APP_ENV"`
		AllowedOrigins string `env-required:"true" yaml:"allowed_origins"      env:"ALLOWED_ORIGINS"`
		SecretKey      string `env-required:"true"                             env:"APP_SECRET_KEY"`
	}

	Logger struct {
		Level string `env-required:"true" yaml:"log_level" env:"LOG_LEVEL"`
	}

	MySQL struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"MYSQL_POOL_MAX"`
		URL     string `env-required:"true"                 env:"MYSQL_URL"`
	}

	Redis struct {
		Port     int    `env-required:"true" yaml:"port"     env:"REDIS_PORT"`
		Host     string `env-required:"true" yaml:"host"     env:"REDIS_HOST"`
		Password string `env-required:"true" yaml:"password" env:"REDIS_PASSWORD"`
	}

	//RMQ struct {
	//	ServerExchange string `env-required:"true" yaml:"rpc_server_exchange" env:"RMQ_RPC_SERVER"`
	//	ClientExchange string `env-required:"true" yaml:"rpc_client_exchange" env:"RMQ_RPC_CLIENT"`
	//	URL            string `env-required:"true"                            env:"RMQ_URL"`
	//}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadConfig("./config/config.yml", cfg); err != nil {
		return nil, fmt.Errorf("read config in yml error: %w", err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read config in env error: %w", err)
	}

	return cfg, nil
}

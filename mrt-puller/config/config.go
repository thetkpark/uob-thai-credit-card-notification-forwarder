package config

import (
	"github.com/thetkpark/uob-thai-credit-card-notification-common/config"
	"github.com/thetkpark/uob-thai-credit-card-notification-common/kv"
	"github.com/thetkpark/uob-thai-credit-card-notification-common/logger"
	"github.com/thetkpark/uob-thai-credit-card-notification-common/publisher"
)

type Config struct {
	Log       logger.Config
	Mrt       MrtApiConfig
	CardId    string `env:"CARD_ID"`
	KV        kv.RedisKVConfig
	Publisher publisher.PubSubConfig
}

type MrtApiConfig struct {
	BaseURL    string `env:"MRT_API_BASE_URL,default=https://api.mangmoomemv.com/v1"`
	Email      string `env:"MRT_API_EMAIL"`
	Password   string `env:"MRT_API_PASSWORD"`
	FetchLimit int    `env:"MRT_API_FETCH_LIMIT,default=50"`
}

func Init() Config {
	var cfg Config
	config.LoadConfigFromENV(&cfg)
	return cfg
}

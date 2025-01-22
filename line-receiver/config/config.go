package config

import (
	"github.com/thetkpark/uob-thai-credit-card-notification-common/config"
	"github.com/thetkpark/uob-thai-credit-card-notification-common/logger"
	"github.com/thetkpark/uob-thai-credit-card-notification-common/publisher"
)

type Config struct {
	Log                  logger.Config
	LineChannelSecret    string   `env:"LINE_CHANNEL_SECRET"`
	WhiteListCardNumbers []string `env:"WHITELIST_CARD_NUMBERS"`
	PubSub               publisher.PubSubConfig
}

func Init() Config {
	var cfg Config
	config.LoadConfigFromENV(&cfg)
	return cfg
}

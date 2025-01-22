package config

import (
	"github.com/Netflix/go-env"
	"github.com/thetkpark/uob-thai-credit-card-notification-common/logger"
	"log"
)

type Config struct {
	Log         logger.Config
	TelegramBot TelegramBot
}

type TelegramBot struct {
	Token         string   `env:"TELEGRAM_BOT_TOKEN"`
	TargetChatIds []string `env:"TELEGRAM_BOT_TARGET_CHAT_IDS"`
}

func Init() Config {
	var cfg Config
	if _, err := env.UnmarshalFromEnviron(&cfg); err != nil {
		log.Fatal(err)
	}
	return cfg
}

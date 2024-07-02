package miniapp

import (
	"os"
	"time"

	initdata "github.com/telegram-mini-apps/init-data-golang"
)

var (
	expIn = 48 * time.Hour
)

func Verify(secret string) error {
	return initdata.Validate(secret, os.Getenv("BOT_TOKEN"), expIn)
}

func Parse(secret string) (initdata.InitData, error) {
	id, err := initdata.Parse(secret)
	if err != nil {
		return initdata.InitData{}, err
	}

	return id, nil
}

func ParseChatID(secret string) (int64, error) {
	id, err := initdata.Parse(secret)
	if err != nil {
		return 0, err
	}

	if id.Chat.ID != 0 {
		return id.Chat.ID, nil
	}
	if id.User.ID != 0 {
		return id.User.ID, nil
	}

	return 0, nil
}

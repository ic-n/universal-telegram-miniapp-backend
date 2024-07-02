package bot

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/ic-n/universal-telegram-miniapp-backend/server/pkg/records"
)

func Bot(ctx context.Context, rdb *records.Database) *bot.Bot {
	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}

	b, err := bot.New(os.Getenv("BOT_TOKEN"), opts...)
	if err != nil {
		log.Fatal(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypePrefix, start(rdb))

	return b
}

func start(rdb *records.Database) func(context.Context, *bot.Bot, *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		var arg string
		args := strings.Split(update.Message.Text, " ")
		if len(args) > 1 {
			arg = args[1]
		}

		switch arg {
		// case "upsell":
		// case "unsub":
		}

		tier := "free" // payments.GetTier(update.Message.Chat.ID)

		var text strings.Builder

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   text.String(),
			ReplyMarkup: models.InlineKeyboardMarkup{
				InlineKeyboard: [][]models.InlineKeyboardButton{{{
					Text: "Launch",
					WebApp: &models.WebAppInfo{
						URL: fmt.Sprintf("%s?from=activity&tier=%s", os.Getenv("MINIAPP_URL"), tier),
					},
				}}},
			},
		})
		if err != nil {
			log.Println(err)
		}
	}
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// wipe any "bot pinned this message" messages in chat.
	if update.Message.From.IsBot && update.Message.Text == "" && update.Message.PinnedMessage.Message != nil {
		b.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    update.Message.Chat.ID,
			MessageID: update.Message.ID,
		})
	}
}

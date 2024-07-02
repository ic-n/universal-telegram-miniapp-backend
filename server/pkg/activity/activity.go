package activity

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/ic-n/universal-telegram-miniapp-backend/server/pkg/miniapp"
	"github.com/ic-n/universal-telegram-miniapp-backend/server/pkg/records"
)

type ActivityRequest struct {
	Secret string
	// aditional fields
}

func Activity(b *bot.Bot, rdb *records.Database) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req ActivityRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := miniapp.Verify(req.Secret); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		chatID, err := miniapp.ParseChatID(req.Secret)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		thisChat, err := b.GetChat(ctx, &bot.GetChatParams{
			ChatID: chatID,
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// ignore if fails
		if thisChat.PinnedMessage != nil {
			b.DeleteMessage(ctx, &bot.DeleteMessageParams{
				ChatID:    chatID,
				MessageID: thisChat.PinnedMessage.ID,
			})
		}

		tier := "free" // payments.GetTier(update.Message.Chat.ID)

		var text strings.Builder
		// todo: text

		m, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:   chatID,
			Text:     text.String(),
			Entities: []models.MessageEntity{},
			ReplyMarkup: models.InlineKeyboardMarkup{
				InlineKeyboard: [][]models.InlineKeyboardButton{{{
					Text: "Resume",
					WebApp: &models.WebAppInfo{
						URL: fmt.Sprintf("%s?from=activity&tier=%s", os.Getenv("MINIAPP_URL"), tier),
					},
				}}},
			},
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err = b.PinChatMessage(ctx, &bot.PinChatMessageParams{
			ChatID:              chatID,
			MessageID:           m.ID,
			DisableNotification: false,
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
}

type ActivityStopRequest struct {
	Secret string
	Start  time.Time
	End    time.Time
}

func ActivityStop(b *bot.Bot, rdb *records.Database) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req ActivityStopRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := miniapp.Verify(req.Secret); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		chatID, err := miniapp.ParseChatID(req.Secret)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		thisChat, err := b.GetChat(ctx, &bot.GetChatParams{
			ChatID: chatID,
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// ignore if fails
		if thisChat.PinnedMessage != nil {
			var text strings.Builder
			text.WriteString(thisChat.PinnedMessage.Text) // todo: edit text

			b.EditMessageText(ctx, &bot.EditMessageTextParams{
				ChatID:      chatID,
				MessageID:   thisChat.PinnedMessage.ID,
				Text:        text.String(),
				ReplyMarkup: models.ReplyKeyboardRemove{},
			})

			b.SetMessageReaction(ctx, &bot.SetMessageReactionParams{
				ChatID:    chatID,
				MessageID: thisChat.PinnedMessage.ID,
				Reaction: []models.ReactionType{
					{
						Type: models.ReactionTypeTypeEmoji,
						ReactionTypeEmoji: &models.ReactionTypeEmoji{
							Type:  models.ReactionTypeTypeEmoji,
							Emoji: "🔥",
						},
					},
				},
			})
		}
	}
}

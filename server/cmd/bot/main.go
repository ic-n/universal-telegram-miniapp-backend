package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/ic-n/universal-telegram-miniapp-backend/server/pkg/api"
	"github.com/ic-n/universal-telegram-miniapp-backend/server/pkg/bot"
	"github.com/ic-n/universal-telegram-miniapp-backend/server/pkg/database"
	"github.com/ic-n/universal-telegram-miniapp-backend/server/pkg/records"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	db, err := database.New("/opt/mycoolapp/main.db")
	if err != nil {
		log.Fatal(err)
	}
	rdb, err := records.New(db)
	if err != nil {
		log.Fatal(err)
	}

	b := bot.Bot(ctx, rdb)

	go api.Server(ctx, b, rdb)

	b.Start(ctx)
}

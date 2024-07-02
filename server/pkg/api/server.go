package api

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/go-telegram/bot"
	"github.com/ic-n/universal-telegram-miniapp-backend/server/pkg/activity"
	"github.com/ic-n/universal-telegram-miniapp-backend/server/pkg/miniapp"
	"github.com/ic-n/universal-telegram-miniapp-backend/server/pkg/records"
)

func Server(ctx context.Context, b *bot.Bot, rdb *records.Database) {
	mux := gin.New()
	mux.Use(logger.SetLogger())

	// for testing
	mux.Use(func(ctx *gin.Context) {
		ctx.Header("Cache-Control", "no-cache")
		ctx.Next()
	})

	// static data
	mux.Static("/assets", "/opt/fn/app/assets")
	mux.Static("/games", "/opt/fn/app/games")
	mux.Static("/images", "/opt/fn/app/images")
	mux.Static("/landing", "/opt/fn/app/landing")
	mux.StaticFile("/miniapp-expand-cheat", "/opt/fn/app/load.html")
	mux.StaticFile("/", "/opt/fn/app/index.html")
	mux.StaticFile("/logo480.png", "/opt/fn/app/logo480.png")

	// api to store and retrieve apend only data
	mux.POST("/api/freestorage", rdb.HandlerRecords)
	mux.PUT("/api/freestorage", rdb.HandlerAddRecord)

	// API for managing activities
	mux.POST("/api/activity", activity.Activity(b, rdb))
	mux.POST("/api/activity/stop", activity.ActivityStop(b, rdb))

	// Strapi Proxy
	strapiAuth := "Bearer " + os.Getenv("STRAPI_TOKEN")
	mux.Any("/content/*path", ReverseProxy("http://localhost:1337", func(r *http.Request) {
		if err := miniapp.Verify(r.Header.Get("Telegram-Init-Data")); err != nil {
			// don't set auth header, not even proxy path rewrite
			return
		}
		r.URL.Path = "/v1" + strings.TrimPrefix(r.URL.Path, "/content")
		r.Header.Set("Authorization", strapiAuth)
	}))

	log.Fatal(autotls.RunWithContext(ctx, mux, "miniapp.myapp.com"))
}

func ReverseProxy(target string, mods ...func(*http.Request)) gin.HandlerFunc {
	url, err := url.Parse(target)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(url)
	return func(c *gin.Context) {
		r := c.Request
		for _, mod := range mods {
			mod(r)
		}

		proxy.ServeHTTP(c.Writer, r)
	}
}

package main

import (
	"log/slog"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"wordle/clients/telegram"
	"wordle/config"
	"wordle/image"
	"wordle/vocab/sqlite"

	"wordle/consumer"
	"wordle/events"
)

func main() {
	logger := slog.Default()

	// init config
	cfg := config.MustLoad()

	// init telegram
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		logger.Error("Can't connect to tg account: ", err)
		os.Exit(1)
	}
	logger.Info("Authorized", "account", bot.Self.UserName)

	// init storage
	vocabs, err := sqlite.New(cfg.SqliteStoragePath)
	if err != nil {
		logger.Error("Can't connect to database: ", err)
		os.Exit(1)
	}

	// init painter
	painter, err := image.New(cfg.Fontfile)
	if err != nil {
		logger.Error("Can't init painter: ", err)
		os.Exit(1)
	}
	// init processor
	eventsProcessor := events.New(
		*telegram.New(bot, cfg.OwnerChatID),
		*vocabs,
		events.NewDialogsTat(),
		*painter,
	)

	// init consumer
	consumer := consumer.New(eventsProcessor, bot)
	if err := consumer.Start(); err != nil {
		logger.Error("Service is stopped: ", err)
		os.Exit(1)
	}
}

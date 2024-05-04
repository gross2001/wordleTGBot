package consumer

import (
	"wordle/events"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Processor interface {
	HandleMessage(domain events.MessageInfo)
	HandleCommand(domain events.MessageInfo)
}

type Consumer struct {
	processor Processor
	tgbot     *tgbotapi.BotAPI
}

func New(processor Processor, tgbot *tgbotapi.BotAPI) *Consumer {
	return &Consumer{
		processor: processor,
		tgbot:     tgbot,
	}
}

func (consumer Consumer) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := consumer.tgbot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		domain := events.MessageInfo{
			ChatID:          update.Message.Chat.ID,
			UserMessage:     update.Message.Text,
			UserCommand:     update.Message.Command(),
			UserMessageID:   update.Message.MessageID,
			UserMessageDate: update.Message.Date,
		}

		if update.Message.IsCommand() {
			consumer.processor.HandleCommand(domain)
			continue
		}
		consumer.processor.HandleMessage(domain)
	}
	return nil
}

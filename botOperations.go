package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func sendMessageByUpdate(bot *tgbotapi.BotAPI, request tgbotapi.Message, message string) {
	msg := tgbotapi.NewMessage(request.Chat.ID, " ")
	msg.ReplyToMessageID = request.MessageID
	msg.Text = message
	bot.Send(msg)
}

func sendMessageByChatID(bot *tgbotapi.BotAPI, chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, " ")
	msg.Text = message
	bot.Send(msg)
}

func sendPhotoByChatID(bot *tgbotapi.BotAPI, chatID int64, imgToSend []byte) (tgbotapi.Message, error) {
	photoFileBytes := tgbotapi.FileBytes{
		Name:  "picture",
		Bytes: imgToSend,
	}
	return (bot.Send(tgbotapi.NewPhoto(chatID, photoFileBytes)))
}

package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Client struct {
	*tgbotapi.BotAPI
	ownerID int64
}

func New(bot *tgbotapi.BotAPI, ownerID int64) *Client {
	return &Client{
		BotAPI:  bot,
		ownerID: ownerID,
	}
}

func (client Client) SendMessageToReply(chatID int64, initialRequest int, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyToMessageID = initialRequest
	client.Send(msg)
}

func (client Client) SendMessageByChatID(chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	client.Send(msg)
}

func (client Client) SendMessageToOwner(message string) {
	msg := tgbotapi.NewMessage(client.ownerID, message)
	client.Send(msg)
}

func (client Client) SendPhotoByChatID(chatID int64, imgToSend []byte) (tgbotapi.Message, error) {
	photoFileBytes := tgbotapi.FileBytes{
		Name:  "picture",
		Bytes: imgToSend,
	}
	return (client.Send(tgbotapi.NewPhoto(chatID, photoFileBytes)))
}

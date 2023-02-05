package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func sendMessageToReply(chatID int64, initialRequest int, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyToMessageID = initialRequest
	bot.Send(msg)
}

func sendMessageByChatID(chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	bot.Send(msg)
}

func sendPhotoByChatID(chatID int64, imgToSend []byte) (tgbotapi.Message, error) {
	photoFileBytes := tgbotapi.FileBytes{
		Name:  "picture",
		Bytes: imgToSend,
	}
	return (bot.Send(tgbotapi.NewPhoto(chatID, photoFileBytes)))
}

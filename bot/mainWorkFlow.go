package bot

import (
	"log"
	"strconv"
	"strings"
	"wordle/vocab"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	bot     *tgbotapi.BotAPI
	dialogs Dialogs
	vocabs  vocab.Vocab
)

func init() {
	dialogs = GetDialogsTat()
	vocabs = vocab.TatVocab{}

	var err error
	if bot, err = tgbotapi.NewBotAPI(token); err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

}

func Start() {
	currentDay := startGameByNewDay()
	users := make(map[int64]usersAnswer, 0)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if (update.Message.Date / 86400) > currentDay.dayNumb {
			currentDay.startNewDay(update.Message.Date)
			sendMessageByChatID(ownerChatID, "number of users "+strconv.Itoa(len(users)))
			users = make(map[int64]usersAnswer, 0)
		}

		if update.Message.IsCommand() {
			handleCommand(update.Message.Command(), update.Message.Chat.ID)
			continue
		}

		userChatId := update.Message.Chat.ID
		userAnswers, _ := users[userChatId]
		userWord := strings.ToLower(update.Message.Text)
		//	log.Printf(".Message.From.chatID[%d] Message %s", userChatId, userWord)

		if userAnswers.IsgameEnd != 0 {
			sendMessageToReply(userChatId, update.Message.MessageID, dialogs.wordsEnded)
			continue
		}

		if ok, reason := checkWordIsOk(userWord); !ok {
			sendMessageToReply(userChatId, update.Message.MessageID, reason)
			continue
		}

		userAnswers.answers = append(userAnswers.answers, userWord)

		msg := sendImageAfterWord(userChatId, currentDay.currentWord, userAnswers)
		userAnswers.IsgameEnd = checkIsGameEnd(userWord, currentDay.currentWord, userAnswers)
		if userAnswers.IsgameEnd != 0 {
			sendFinalMessage(userChatId, update.Message.MessageID, currentDay.currentWord, userAnswers)
		}
		if userAnswers.lastMsgID != 0 {
			msgToDelete := tgbotapi.NewDeleteMessage(userChatId, userAnswers.lastMsgID)
			bot.Send(msgToDelete)
		}
		userAnswers.lastMsgID = msg.MessageID

		users[userChatId] = userAnswers
		//log.Println("Number of users", len(users))
	}
}

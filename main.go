package main

import (
	"log"
	"strconv"
	"strings"
	imgwordle "wordle/imagewordle"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	_ int8 = iota
	userWins
	userLoses
)

func main() {

	var currentDay dayInfo
	currentDay.dayNumb = timeOfStart / 86400
	currentDay.currentWordNumb = 0
	currentDay.currentWord = answers[currentDay.currentWordNumb]
	log.Println("Current day is ", currentDay.dayNumb)
	log.Println("Current word is ", currentDay.currentWord)

	users := make(map[int64]usersAnswer, 0)

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			userWord := strings.ToLower(update.Message.Text)
			userChatId := update.Message.Chat.ID
			log.Printf(".Message.From.chatID[%d] Message %s", userChatId, userWord)

			if (update.Message.Date / 86400) > currentDay.dayNumb {
				startNewDay(&currentDay, update.Message.Date)
				sendMessageByChatID(bot, ownerChatID, "number of users "+strconv.Itoa(len(users)))
				users = make(map[int64]usersAnswer, 0)
			}

			userAnswer, _ := users[userChatId]

			if userAnswer.IsgameEnd != 0 {
				sendMessageByUpdate(bot, *update.Message, wordsEnded)
				continue
			}

			if ok, reason := checkWordIsOk(userWord); !ok {
				sendMessageByUpdate(bot, *update.Message, reason)
				continue
			}

			userAnswer.answers = append(userAnswer.answers, userWord)

			imgToSend := imgwordle.FullImage(currentDay.currentWord, userAnswer.answers)
			sendPhotoByChatID(bot, userChatId, imgToSend)

			log.Println(userAnswer.answers)

			if currentDay.currentWord == userWord {
				userAnswer.IsgameEnd = userWins
			} else if len(userAnswer.answers) == 6 {
				userAnswer.IsgameEnd = userLoses
			}

			if userAnswer.IsgameEnd != 0 {
				var finalMessage string
				if userAnswer.IsgameEnd == userWins {
					finalMessage = youWin
				}
				if userAnswer.IsgameEnd == userLoses {
					finalMessage = correctWordIs + currentDay.currentWord
				}
				sendMessageByUpdate(bot, *update.Message, finalMessage)
				imgToSend := imgwordle.RectOnly(currentDay.currentWord, userAnswer.answers)
				result, err := sendPhotoByChatID(bot, userChatId, imgToSend)
				if err == nil {
					sendMessageByUpdate(bot, result, pictureToFriends)
				}
			}

			users[userChatId] = userAnswer
			log.Println("Number of users", len(users))
		}
	}
}

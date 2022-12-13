package main

import (
	"log"
	"strconv"
	"strings"
	imgwordle "wordle/imagewordle"
	"wordle/vocab"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type usersAnswer struct {
	riddleNumb int
	riddle     string
	answers    []string
}

func startNewGame(userAnswer usersAnswer) usersAnswer {
	userAnswer.riddleNumb++
	if userAnswer.riddleNumb < len(answers) {
		userAnswer.riddle = answers[userAnswer.riddleNumb]
	}
	userAnswer.answers = make([]string, 0)
	return userAnswer
}

func main() {

	users := make(map[string]usersAnswer, 0)
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	//	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	endOfDay := ((timeOfStart / 86400) + 1) * 86400
	days := 0
	currentWord := 0

	var ownerChatId int64
	for update := range updates {

		if ownerChatId == 0 {
			if update.Message.From.UserName == botOwner {
				ownerChatId = update.Message.Chat.ID
			}
		}
		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, " ")
			msg.ReplyToMessageID = update.Message.MessageID

			wordFromUser := strings.ToLower(update.Message.Text)
			userName := update.Message.From.UserName
			log.Printf(".Message.From.UserName[%s] Message %s", update.Message.From.UserName, update.Message.Text)

			if len([]rune(wordFromUser)) != 5 {
				msg.Text = onlyFiveLetter
				bot.Send(msg)
				continue
			}

			wordFromUser = vocab.DoRequest(wordFromUser)
			if len(wordFromUser) == 0 {
				msg.Text = wordNotExist
				bot.Send(msg)
				continue
			}

			if update.Message.Date > endOfDay {
				msgToOwnew := tgbotapi.NewMessage(ownerChatId, "Number of users for a last day is "+strconv.Itoa(len(users)))
				bot.Send(msgToOwnew)
				users = make(map[string]usersAnswer, 0)
				endOfDay = ((update.Message.Date / 86400) + 1) * 86400
				days = (update.Message.Date - timeOfStart) / 86400
				currentWord = days * wordsPerDay
				log.Println(days)
			}

			var userAnswer usersAnswer
			var ok bool
			if userAnswer, ok = users[userName]; !ok {
				userAnswer.riddle = answers[currentWord]
			}

			if userAnswer.riddleNumb >= wordsPerDay {
				msg.Text = wordsEnded
				bot.Send(msg)
				continue
			}

			userAnswer.answers = append(userAnswer.answers, wordFromUser)

			imgToSend := imgwordle.CreateImage(userAnswer.riddle, userAnswer.answers)

			photoFileBytes := tgbotapi.FileBytes{
				Name:  "picture",
				Bytes: imgToSend,
			}
			bot.Send(tgbotapi.NewPhoto(int64(update.Message.Chat.ID), photoFileBytes))

			log.Println(userAnswer)

			if userAnswer.riddle == wordFromUser {
				msg.Text = youWin + startingNewGame
				userAnswer = startNewGame(userAnswer)
			} else if len(userAnswer.answers) == 6 {
				msg.Text = correctWordIs + userAnswer.riddle + ".\n" + startingNewGame
				userAnswer = startNewGame(userAnswer)
			} else {
				msg.Text = numberOfTry + strconv.Itoa(6-len(userAnswer.answers))
			}

			bot.Send(msg)

			users[userName] = userAnswer
			log.Println("Number of users", len(users))
		}
	}

}

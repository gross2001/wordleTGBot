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

	for update := range updates {

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

			var userAnswer usersAnswer
			var ok bool
			if userAnswer, ok = users[userName]; !ok {
				userAnswer.riddle = answers[0]
			}

			if userAnswer.riddleNumb >= len(answers) {
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
			log.Println(userAnswer)
		}
	}

}

package main

import (
	"log"
	"strconv"
	"strings"
	imgwordle "wordle/imagewordle"
	"wordle/vocab"
	"net/http"


	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func giveNextString(strSlice []string, value string) string {
	for p, v := range strSlice[:len(strSlice)-1] {
		if v == value {
			return strSlice[p+1]
		}
	}
	return ""
}

func main() {

	var answers = []string{"мамык", "капка", "көрәк", "бакча", "дәрес", "китап"}
	var users map[string][]string

	bot, err := tgbotapi.NewBotAPI("5347435152:AAFtYhWbwK19LdIXypTJD6sZ-qoiGIKyHPU")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	
	go http.ListenAndServe("0.0.0.0:443", nil)

	for update := range updates {
		
		if update.Message != nil {
			wordFromUser := update.Message.Text
			userName := update.Message.From.UserName
			log.Printf(".Message.From.UserName[%s] Message %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, " ")
			msg.ReplyToMessageID = update.Message.MessageID

			if len([]rune(wordFromUser)) != 5 {
				msg.Text = "Сүз 5 хәрефтән торорга тиеш"
				bot.Send(msg)
				continue
			}
			//	isWordExist := vocab.DoRequest(strings.ToLower(wordFromUser))
			wordFromUser = vocab.DoRequest(strings.ToLower(wordFromUser))
			if len(wordFromUser) == 0 {
				msg.Text = "Сүзлектә мондый сүз табылмады. ә/э, ү/у, җ/ж, ң/н, һ/х, ө/о"
				bot.Send(msg)
				continue
			}

			if _, ok := users[userName]; ok {
				users[userName] = append(users[userName], strings.ToLower(wordFromUser))
			} else {
				users = make(map[string][]string, 0)
				users[userName] = append(users[userName], answers[0])
				users[userName] = append(users[userName], strings.ToLower(wordFromUser))
			}
			imgToSend := imgwordle.CreateImage(users[userName][0], users[userName][1:])

			photoFileBytes := tgbotapi.FileBytes{
				Name:  "picture",
				Bytes: imgToSend,
			}
			bot.Send(tgbotapi.NewPhoto(int64(update.Message.Chat.ID), photoFileBytes))

			var needNewWord int
			if users[userName][0] == users[userName][len(users[userName])-1] {
				msg.Text = "Дөрес! Яңа сүз уйлап куйдым!"
				needNewWord = 1
			} else if len(users[userName]) == 7 {
				msg.Text = "Дөрес сүз булды: " + users[userName][0] + ". Яңа сүз уйлап куйдым!"
				needNewWord = 1
			} else {
				msg.Text = "Син тагын " + strconv.Itoa(7-len(users[userName])) + " кабат җавап бираласен"
			}
			bot.Send(msg)

			if needNewWord == 1 {
				users[userName][0] = giveNextString(answers, users[userName][0])
				users[userName] = users[userName][:1]
			}

		}
	}

}

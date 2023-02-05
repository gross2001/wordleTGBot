package bot

import (
	"log"
	"time"
	image "wordle/image"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	gameInProcess int8 = iota
	userWins
	userLoses
)

type usersAnswer struct {
	IsgameEnd int8
	answers   []string
	lastMsgID int
}

type dayInfo struct {
	dayNumb         int
	currentWordNumb int
	currentWord     string
}

func startGameByNewDay() dayInfo {
	currentDay := dayInfo{
		dayNumb:         timeOfStart / 86400,
		currentWordNumb: 0,
		currentWord:     riddles[0],
	}
	log.Println("Current day is ", currentDay.dayNumb)
	log.Println("Current word is ", currentDay.currentWord)
	return currentDay
}

func (currentDay dayInfo) startNewDay(updateDate int) {
	currentDay.dayNumb = updateDate / 86400
	currentDay.currentWordNumb = currentDay.dayNumb - timeOfStart/86400
	currentDay.currentWord = riddles[currentDay.currentWordNumb]
	log.Println("New day is started")
	log.Println("Current day is ", currentDay.dayNumb)
	log.Println("Current word is ", currentDay.currentWord)
}

func checkWordIsOk(word string) (bool, string) {
	if len([]rune(word)) != 5 {
		return false, dialogs.onlyFiveLetter
	}
	if ok := vocabs.DoRequest(word); !ok {
		return false, dialogs.wordNotExist
	}
	return true, ""
}

func handleCommand(command string, chatId int64) {
	var text string
	switch command {
	case "start":
		text = dialogs.linksToRules
	case "help":
		text = dialogs.linksToRules
	default:
	}
	if len(text) > 0 {
		msg := tgbotapi.NewMessage(chatId, text)
		bot.Send(msg)
	}
}

func checkIsGameEnd(userWord string, currentWord string, userAnswers usersAnswer) int8 {
	if userWord == currentWord {
		return userWins
	} else if len(userAnswers.answers) == 6 {
		return userLoses
	}
	return gameInProcess
}

func sendImageAfterWord(chatId int64, currentWord string, userAnswers usersAnswer) tgbotapi.Message {
	log.Println(userAnswers)

	imgToSend := image.FullImage(currentWord, userAnswers.answers)
	msg, _ := sendPhotoByChatID(chatId, imgToSend)
	log.Println(userAnswers.answers)
	return msg
}

func sendFinalMessage(chatId int64, initialRequest int, currentWord string, userAnswers usersAnswer) {
	var finalMessage string

	if userAnswers.IsgameEnd == userWins {
		finalMessage = dialogs.youWin
	}
	if userAnswers.IsgameEnd == userLoses {
		finalMessage = dialogs.correctWordIs + currentWord
	}
	sendMessageToReply(chatId, initialRequest, finalMessage)
	imgToSend := image.RectOnly(currentWord, userAnswers.answers)

	time.Sleep(time.Millisecond * 500)
	result, err := sendPhotoByChatID(chatId, imgToSend)
	if err == nil {
		sendMessageToReply(chatId, result.MessageID, dialogs.pictureToFriends)
	}
}

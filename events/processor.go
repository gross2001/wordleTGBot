package events

import (
	"log/slog"
	"net/url"
	"strconv"
	"strings"
	"time"
	"wordle/clients/telegram"
	"wordle/image"
	"wordle/vocab/sqlite"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const secPerDay = 86400

const (
	GameInProcess int8 = iota
	UserWins
	UserLoses
)

type Processor struct {
	vocab      sqlite.TatSQLVocab
	client     telegram.Client
	currentDay DayInfo
	users      map[int64]UsersAnswer
	dialogs    Dialogs
	painter    image.Painter
}

func New(Client telegram.Client, Vocab sqlite.TatSQLVocab, dialogs Dialogs, painter image.Painter) *Processor {
	users := make(map[int64]UsersAnswer, 0)
	currentDay := DayInfo{}
	return &Processor{
		client:     Client,
		vocab:      Vocab,
		currentDay: currentDay,
		dialogs:    dialogs,
		users:      users,
		painter:    painter,
	}
}

func (p *Processor) HandleMessage(domain MessageInfo) {
	logger := slog.Default()

	if (domain.UserMessageDate / secPerDay) > p.currentDay.DayNumb {
		p.updateDay(domain)
	}

	userAnswers := p.users[domain.ChatID]
	userWord := strings.ToLower(domain.UserMessage)
	logger.Info("New message", "chatID", domain.ChatID, "Message", userWord)

	// check game status
	if userAnswers.IsgameEnd != 0 {
		p.client.SendMessageToReply(domain.ChatID, domain.UserMessageID, p.dialogs.WordsEnded)
		return
	}

	// check word
	if ok, reason := p.checkWordIsOk(userWord); !ok {
		p.client.SendMessageToReply(domain.ChatID, domain.UserMessageID, reason)
		return
	}
	userAnswers.Answers = append(userAnswers.Answers, userWord)

	// generate image
	imgToSend := p.painter.FullImage(p.currentDay.CurrentWord, userAnswers.Answers)
	msg, err := p.client.SendPhotoByChatID(domain.ChatID, imgToSend)
	if err != nil {
		// TODO
	}

	// end game
	userAnswers.IsgameEnd = p.checkIsGameEnd(userWord, userAnswers)
	if userAnswers.IsgameEnd != 0 {
		p.endGame(domain, userAnswers)
	}

	// delete unnecessary messages
	if userAnswers.LastMsgID != 0 {
		msgToDelete := tgbotapi.NewDeleteMessage(domain.ChatID, userAnswers.LastMsgID)
		p.client.Send(msgToDelete)
	}
	userAnswers.LastMsgID = msg.MessageID

	p.users[domain.ChatID] = userAnswers
}

func (p *Processor) HandleCommand(domain MessageInfo) {
	var text string
	switch domain.UserCommand {
	case "start":
		text = p.dialogs.LinksToRules
	case "help":
		text = p.dialogs.LinksToRules
	default:
	}
	if len(text) > 0 {
		msg := tgbotapi.NewMessage(domain.ChatID, text)
		p.client.Send(msg)
	}
}

func (p *Processor) checkWordIsOk(word string) (bool, string) {
	if len([]rune(word)) != 5 {
		return false, p.dialogs.OnlyFiveLetter
	}
	ok, err := p.vocab.DoRequest(word)
	if err != nil || !ok {
		return false, p.dialogs.WordNotExist
	}
	return true, ""
}

func (p *Processor) checkIsGameEnd(userWord string, userAnswers UsersAnswer) int8 {
	if userWord == p.currentDay.CurrentWord {
		return UserWins
	} else if len(userAnswers.Answers) == 6 {
		return UserLoses
	}
	return GameInProcess
}

func (p *Processor) updateDay(domain MessageInfo) {
	logger := slog.Default()

	p.currentDay.DayNumb = domain.UserMessageDate / secPerDay
	var err error
	p.currentDay.CurrentWord, err = p.vocab.ChooseWord(p.currentDay.DayNumb)
	logger.Info("New riddle is " + p.currentDay.CurrentWord)
	if err != nil {
		logger.Error("Can't take new words: ", err)
	}
	p.client.SendMessageToOwner("number of users " + strconv.Itoa(len(p.users)))
	p.users = make(map[int64]UsersAnswer, 0)
}

func (p *Processor) endGame(domain MessageInfo, userAnswers UsersAnswer) {
	var finalMessage string

	switch userAnswers.IsgameEnd {
	case UserWins:
		finalMessage = p.dialogs.YouWin
	case UserLoses:
		finalMessage = p.dialogs.CorrectWordIs + p.currentDay.CurrentWord

		if p.dialogs.DictURL != "" && p.dialogs.Title == "tatar" {
			finalMessage += p.parseTatarDictURL()
		}
	}
	p.client.SendMessageToReply(domain.ChatID, domain.UserMessageID, finalMessage)
	imgToSend := p.painter.RectOnly(p.currentDay.CurrentWord, userAnswers.Answers)
	time.Sleep(time.Millisecond * 1000)
	result, err := p.client.SendPhotoByChatID(domain.ChatID, imgToSend)
	if err == nil {
		p.client.SendMessageToReply(domain.ChatID, result.MessageID, p.dialogs.PictureToFriends)

	}
}

func (p *Processor) parseTatarDictURL() string {
	logger := slog.Default()

	baseUrl, err := url.Parse(p.dialogs.DictURL)
	if err != nil {
		logger.Error("Can't parse a url to dict: ", err)
		return ""
	}

	params := url.Values{}
	params.Add("txtW", p.currentDay.CurrentWord)
	params.Add("source[]", "9")
	baseUrl.RawQuery = params.Encode()

	return ("\n" + p.dialogs.LinkToDict + baseUrl.String())
}

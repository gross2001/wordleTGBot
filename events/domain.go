package events

type MessageInfo struct {
	ChatID          int64
	UserMessage     string
	UserCommand     string
	UserMessageID   int
	UserMessageDate int
}

type DayInfo struct {
	DayNumb     int
	CurrentWord string
}

type UsersAnswer struct {
	IsgameEnd int8
	Answers   []string
	LastMsgID int
}

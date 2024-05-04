package events

type Dialogs struct {
	OnlyFiveLetter   string
	WordNotExist     string
	NumberOfTry      string
	CorrectWordIs    string
	YouWin           string
	StartingNewGame  string
	WordsEnded       string
	PictureToFriends string
	LinksToRules     string
}

func NewDialogsTat() Dialogs {
	return Dialogs{
		OnlyFiveLetter:   "Исем 5 хәрефтән торырга тиеш ",
		WordNotExist:     "Сүзлектә мондый исем табылмады ",
		NumberOfTry:      "Җавап бирү мөмкинлеге: ",
		CorrectWordIs:    "Дөрес сүз булды: ",
		YouWin:           "Дөрес! ",
		StartingNewGame:  "Яңа сүз уйлап куйдым! ",
		WordsEnded:       "Бүгенгә сүзләр бетте. Иртәгә кил! ",
		PictureToFriends: "Дусларына җибәрер өчен сурәт. Кем сүзне тизрәк табала;)\n @TatarWordle_bot ",
		LinksToRules:     "Кагыйдәләргә сылтама - https://telegra.ph/Rulet-of-a-wordle-12-13",
	}
}

func NewDialogsUdm() Dialogs {
	return Dialogs{
		OnlyFiveLetter:   "Кыл 5 буквалэсь гинэ луыны быгатэ ",
		WordNotExist:     "Кыллюкамын сыӵе кыл ӧз сюры ",
		CorrectWordIs:    "Шонер кыл вал таӵе: ",
		YouWin:           "Шонер! ",
		WordsEnded:       "Туннэлы кылъёс быризы. Лыкты ӵуказе! ",
		PictureToFriends: "Та суредэз эшъёсыдлы келя. Ӵошатске, кин кылэз ӝоггес шедьтоз! \n @UdmurtWordleBot ",
		LinksToRules:     "Кызьы та шудонэн шудоно: https://telegra.ph/Kyzy-UdmurtWordle-shudonehn-shudono-01-17-2",
	}
}

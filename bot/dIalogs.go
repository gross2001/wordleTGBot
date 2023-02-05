package bot

type Dialogs struct {
	onlyFiveLetter   string
	wordNotExist     string
	numberOfTry      string
	correctWordIs    string
	youWin           string
	startingNewGame  string
	wordsEnded       string
	pictureToFriends string
	linksToRules     string
}

func GetDialogsTat() Dialogs {
	return Dialogs{
		onlyFiveLetter:   "Исем 5 хәрефтән торырга тиеш ",
		wordNotExist:     "Сүзлектә мондый исем табылмады ",
		numberOfTry:      "Җавап бирү мөмкинлеге: ",
		correctWordIs:    "Дөрес сүз булды: ",
		youWin:           "Дөрес! ",
		startingNewGame:  "Яңа сүз уйлап куйдым! ",
		wordsEnded:       "Бүгенгә сүзләр бетте. Иртәгә кил! ",
		pictureToFriends: "Дусларына җибәрер өчен сурәт. Кем сүзне тизрәк табала;)\n @TatarWordle_bot ",
		linksToRules:     "Кагыйдәләргә сылтама - https://telegra.ph/Rulet-of-a-wordle-12-13",
	}
}

func GetDialogsUdm() Dialogs {
	return Dialogs{
		onlyFiveLetter:   "Кыл 5 буквалэсь гинэ луыны быгатэ ",
		wordNotExist:     "Кыллюкамын сыӵе кыл ӧз сюры ",
		correctWordIs:    "Шонер кыл вал таӵе: ",
		youWin:           "Шонер! ",
		wordsEnded:       "Туннэлы кылъёс быризы. Лыкты ӵуказе! ",
		pictureToFriends: "Та суредэз эшъёсыдлы келя. Ӵошатске, кин кылэз ӝоггес шедьтоз! \n @UdmurtWordleBot ",
		linksToRules:     "Кызьы та шудонэн шудоно: https://telegra.ph/Kyzy-UdmurtWordle-shudonehn-shudono-01-17-2",
	}
}

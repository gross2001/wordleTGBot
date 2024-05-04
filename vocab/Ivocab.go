package vocab

type Vocab interface {
	DoRequest(string) (bool, error)
	ChooseWord(int) (string, error)
}

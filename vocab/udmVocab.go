package vocab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	korpus = "http://udmcorpus.udman.ru/api/public/dictionary/search"
)

type UdmVocab struct{}

func (udmVocab UdmVocab) DoRequest(word string) bool {

	if jsonBytes, err := udmGetJSON(word); err != nil {
		log.Printf("Failed to get XML: %v", err)
	} else {
		if len(jsonBytes) > 2 {
			return true
		}
	}
	return false
}

type udmRequest struct {
	Langid int     `json:"langid"`
	Word   string  `json:"word"`
	Lang   udmLang `json:"lang"`
}

type udmLang struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Prefix string `json:"prefix"`
}

func udmcontainsNoun() bool {
	return false
}

func udmGetJSON(word string) ([]byte, error) {
	url := korpus

	initLang := udmLang{Id: 1, Title: "Удмуртский", Prefix: "udm"}
	initial := udmRequest{Langid: 1, Lang: initLang, Word: word}
	postBody, _ := json.Marshal(initial)

	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(url, "application/json", responseBody)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("read body: %v", err)
	}

	return data, nil
}

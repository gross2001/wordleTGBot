package vocab

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	yandex    = "dict.1.1.20221203T124750Z.44df4b4d5ab923fc.21ee7e9333917cd13a50c8fc552bc7590a1d7e43"
	lang      = "tt-ru"
	urlToDict = "https://dictionary.yandex.net/api/v1/dicservice.json/lookup?key="
)

type myJSONstruct struct {
	Def []def `json:"def"`
}

type def struct {
	PartOfSpeech string `json:"pos"`
}

func containsNoun(s []def) bool {
	for _, a := range s {
		if a.PartOfSpeech == "noun" {
			return true
		}
	}
	return false
}

func getJSON(word string) ([]byte, error) {
	//	url := url.QueryEscape(urlToDict + yandex + "&lang=" + lang + "&text=" + word)
	url := urlToDict + yandex + "&lang=" + lang + "&text=" + word
	log.Println(url)

	resp, err := http.Get(url)
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

func DoRequest(word string) string {

	var result myJSONstruct

	if jsonBytes, err := getJSON(word); err != nil {
		log.Printf("Failed to get XML: %v", err)
	} else {
		json.Unmarshal(jsonBytes, &result)
		if containsNoun(result.Def) == true {
			return word
		}
	}

	return ""

}

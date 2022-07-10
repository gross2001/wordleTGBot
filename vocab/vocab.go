package vocab

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type myXMLstruct struct {
	XMLName      xml.Name `xml:"res"`
	Text         string   `xml:",chardata"`
	ResponseType string   `xml:"responseType"`
	Word         string   `xml:"word"`
	POS          string   `xml:"POS"`
	Translation  string   `xml:"translation"`
	Examples     string   `xml:"examples"`
	Mt           string   `xml:"mt"`
}

func getXML(url string) ([]byte, error) {
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

func giveAllWords(word string) []string {
	var dict map[rune]string = make(map[rune]string, 0)
	dict['э'] = "ә"
	dict['у'] = "ү"
	dict['ж'] = "җ"
	dict['н'] = "ң"
	dict['х'] = "һ"
	dict['о'] = "ө"

	var result []string
	var anotherStr string
	result = append(result, word)
	for _, char := range word {
		for key := range dict {
			if char == key {
				anotherStr = strings.Replace(word, string(key), dict[key], 1)
				result = append(result, giveAllWords(anotherStr)...)
			}
		}
	}
	return result
}

func DoRequest(word string) string {

	var result myXMLstruct
	allWords := giveAllWords(word)
	log.Println(allWords)

	for _, everyWord := range allWords {
		if xmlBytes, err := getXML("https://translate.tatar/translate?lang=1&text=" + url.QueryEscape(everyWord)); err != nil {
			log.Printf("Failed to get XML: %v", err)
		} else {
			xml.Unmarshal(xmlBytes, &result)
			if len(result.Word) != 0 {
				return everyWord
			}
		}
	}

	return ""

}

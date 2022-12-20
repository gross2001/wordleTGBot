package main

import (
	"log"
	"wordle/vocab"
)

func checkWordIsOk(word string) (bool, string) {
	if len([]rune(word)) != 5 {
		return false, onlyFiveLetter
	}
	if ok := vocab.DoRequest(word); !ok {
		return false, wordNotExist
	}
	return true, ""
}

func startNewDay(currentDay *dayInfo, updateDate int) {
	currentDay.dayNumb = updateDate / 86400
	currentDay.currentWordNumb = currentDay.dayNumb - timeOfStart/86400
	currentDay.currentWord = answers[currentDay.currentWordNumb]
	log.Println("New day is started")
	log.Println("Current day is ", currentDay.dayNumb)
	log.Println("Current word is ", currentDay.currentWord)
}

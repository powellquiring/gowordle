package main

import (
	"github.com/powellquiring/gowordle/gowordle"
)

func FirstWithInitialGuesses() {
	// wordList := WordleDictionary[0:800]
	// wordList := gowordle.WordleDictionary[0:50]
	wordList := gowordle.WordleDictionary[0:]
	score, words := gowordle.FirstGuessProvideInitialGuesses1(wordList, wordList)
	print(score)
	gowordle.PrintWords(words)
}

func main() {
	FirstWithInitialGuesses()
}

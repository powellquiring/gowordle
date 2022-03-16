package gowordle

import (
	"sort"
	"testing"

	"github.com/bits-and-blooms/bitset"
	"github.com/stretchr/testify/assert"
)

func TestMany(t *testing.T) {
	set := bitset.New(5)
	set.Set(0)
	set.Set(3)
	indices := make([]uint, set.Count())
	set.NextSetMany(0, indices)
}

func TestBest(t *testing.T) {
	words := StringsToWordleWords([]string{"aaaaa", "abbbb"})
	score := NextGuess(words, words)
	assert.NotZero(t, score)
	print("testbest")
}

func WW(ins string) WordleWord {
	return []rune(ins)
}
func TestMatching1(t *testing.T) {
	words := StringsToWordleWords([]string{"aaaaa", "abbbb"})
	wds := NewWordleMatcher(words)
	assert := assert.New(t)
	matching := wds.matching(WW("aazzz"), WW("ggrrr"))
	assert.Equal(matching, StringsToWordleWords([]string{"aaaaa"}))

	matching = wds.matching(WW("bzzzz"), WW("yrrrr"))
	assert.Equal(matching, StringsToWordleWords([]string{"abbbb"}))
}

func TestMatching2(t *testing.T) {
	words := StringsToWordleWords([]string{"aaaaa", "abbbb"})
	wds := NewWordleMatcher(words)
	assert := assert.New(t)
	matching := wds.matching(WW("bzzzz"), WW("yrrrr"))
	assert.Equal(matching, StringsToWordleWords([]string{"abbbb"}))
}

func WordSort(ws []WordleWord) []string {
	s := make([]string, len(ws))
	for i, w := range ws {
		s[i] = string(w)
	}
	sort.Strings(s)
	return s
}

func testMatching(t *testing.T, words []string, guess string, answer string, expected []string) {
	wwords := StringsToWordleWords(words)
	wds := NewWordleMatcher(wwords)
	matching := wds.matching(WW(guess), WW(answer))
	assert := assert.New(t)
	sort.Strings(expected)
	matching_s := WordSort(matching)
	assert.Equal(expected, matching_s)
}

func TestMatching3(t *testing.T) {
	testMatching(t,
		[]string{"aaazz", "abbbb", "bcazz"},
		"bxxac", "yrryr", // answer abbbb
		[]string{"abbbb"},
	)
}

func TestMatching4(t *testing.T) {
	testMatching(t,
		[]string{"aaazz", "abbzz", "abczz", "abazz", "bbazz"},
		"xabxx", "ryyrr", // answer abazz
		[]string{"abczz", "abazz", "bbazz"},
	)
}

func TestGreenYellow(t *testing.T) {
	testMatching(t,
		[]string{"aaazz", "abbzz", "abczz", "abazz", "bbazz", "azzza", "azzzz"},
		"axxxa", "grrry", // answer abazz
		[]string{"aaazz", "abazz"},
	)
}

func TestYellowRed(t *testing.T) {
	testMatching(t,
		[]string{"aaazz", "abbzz", "abczz", "abazz", "bbazz", "azzza", "azzzz", "aazzz", "aaazz"},
		"axxaa", "grryr", // answer abazz, two a's, but not 3
		[]string{"abazz", "aazzz"},
	)
}

func TestW1(t *testing.T) {
	testMatching(t,
		WordleDictionary,
		"aaxxd", "yyrry",
		[]string{"drama"},
	)
}

func TestFirst20(t *testing.T) {
	wordList := []string{"cigar", "rebut", "sissy", "humph", "awake", "blush", "focal", "evade", "naval", "serve", "heath", "dwarf", "model", "karma", "stink", "grade", "quiet", "bench", "abate", "feign"}
	simulateWords := Simulate(wordList, "karma", "cigar")
	assert := assert.New(t)
	assert.Greater(float32(5.0), float32(len(simulateWords)))
}

func TestFirst(t *testing.T) {
	wordList := WordleDictionary[0:100]
	// wordList := wordleDictionary
	FirstGuess(wordList)
}

func TestFirst1(t *testing.T) {
	wordList := WordleDictionary[0:800]
	score, words := FirstGuess1(wordList)
	print(score)
	PrintWords(words)
}

func TestFirstWithInitialGuesses(t *testing.T) {
	// wordList := wordleDictionary[0:800]
	wordList := WordleDictionary[0:1100]
	score, words := FirstGuessProvideInitialGuesses1(wordList, wordList)
	print(score)
	PrintWords(words)
}

/*
func BenchmarkFirstN(t *testing.B) {
	wordList := WordleDictionary[0:400]
	FirstGuess(wordList)
}
*/

func TestAgainstHeron400(t *testing.T) {
	// tested and got cigar/4 for 0:400
	wordList := WordleDictionary[0:400]
	// bug: cigar, rebut, serve, ferry, heron
	guess := "cigar"
	worst := 4
	solution := "heron"
	simulateWords := Simulate(wordList, solution, guess)
	if len(simulateWords) > worst {
		println(len(simulateWords), string(solution))
		worst = len(simulateWords)
	}
}
func TestAgainstHeronArise(t *testing.T) {
	wordList := WordleDictionary[0:]
	guess := "arise"
	worst := 4
	solution := "heron"
	simulateWords := Simulate(wordList, solution, guess)
	if len(simulateWords) > worst {
		println(len(simulateWords), string(solution))
		worst = len(simulateWords)
	}
}
func TestAgainstServeArise(t *testing.T) {
	wordList := WordleDictionary[0:600] // contains arise
	guess := "arise"
	worst := 4
	solution := "serve"
	assertStringInSlice(t, guess, wordList)
	assertStringInSlice(t, solution, wordList)
	simulateWords := Simulate(wordList, solution, guess)
	if len(simulateWords) > worst {
		println(len(simulateWords), string(solution))
		worst = len(simulateWords)
	}
}
func assertStringInSlice(t *testing.T, tst string, wordList []string) {
	for _, word := range wordList {
		if word == tst {
			return
		}
	}
	t.Error("word not in list: " + tst)
}
func TestAriseAgainstAll(t *testing.T) {
	// tested and got cigar/4 for 0:200
	// tested and got cigar/4 for 0:400
	wordList := WordleDictionary[0:600] // contains arise
	// bug: cigar, rebut, serve, ferry, heron
	guess := "arise"
	assertStringInSlice(t, guess, wordList)
	worst := 0
	result := []string{}
	for i, solution := range wordList {
		simulateWords := Simulate(wordList, solution, guess)
		assert.Equal(t, solution, string(simulateWords[len(simulateWords)-1]))
		if i%10 == 0 {
			println(i)
		}
		if len(simulateWords) == worst {
			result = append(result, string(solution))
		} else if len(simulateWords) > worst {
			println(len(simulateWords), string(solution))
			worst = len(simulateWords)
			result = []string{string(solution)}
		}
	}
	println(worst, result)
}

func TestPlayArise(t *testing.T) {
	wordList := StringsToWordleWords(WordleDictionary[0:])
	guessAnswers := []GuessAnswer{
		{[]rune("arise"), []rune("yrrgg")},
		{[]rune("cigar"), []rune("rrryr")},
		{[]rune("panel"), []rune("ryryr")},
	}
	print(string(playWordle(wordList, guessAnswers)))
	print("done")
}

/*
func BenchmarkProof(t *testing.B) {
	wordList := wordleDictionary[0:400]
	_, ret := FirstGuess(wordList)
	assert := assert.New(t)
	for _, solution := range wordList {
		simulateWords := Simulate(wordList, solution, ret)
		assert.Greater(float32(7.0), float32(len(simulateWords)))
	}
}

*/

func BenchmarkFirst1(t *testing.B) {
	wordList := WordleDictionary[0:400]
	score, words := FirstGuess1(wordList)
	print(score)
	PrintWords(words)
}

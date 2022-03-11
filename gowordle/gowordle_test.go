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
	words := NewWordleWords([]string{"aaaaa", "abbbb"})
	score := NextGuess(words, words)
	assert.NotZero(t, score)
	print("testbest")
}

func WW(ins string) WordleWord {
	return []rune(ins)
}
func TestMatching1(t *testing.T) {
	words := NewWordleWords([]string{"aaaaa", "abbbb"})
	wds := NewWordleMatcher(words)
	assert := assert.New(t)
	matching := wds.matching(WW("aazzz"), WW("ggrrr"))
	assert.Equal(matching, NewWordleWords([]string{"aaaaa"}))

	matching = wds.matching(WW("bzzzz"), WW("yrrrr"))
	assert.Equal(matching, NewWordleWords([]string{"abbbb"}))
}

func TestMatching2(t *testing.T) {
	words := NewWordleWords([]string{"aaaaa", "abbbb"})
	wds := NewWordleMatcher(words)
	assert := assert.New(t)
	matching := wds.matching(WW("bzzzz"), WW("yrrrr"))
	assert.Equal(matching, NewWordleWords([]string{"abbbb"}))
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
	wwords := NewWordleWords(words)
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
		wordleDictionary,
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

func BenchmarkGuessN(t *testing.B) {
	wordList := wordleDictionary[0:400]
	_, ret := FirstGuess(wordList)
	assert := assert.New(t)
	for _, solution := range wordList {
		simulateWords := Simulate(wordList, solution, ret)
		assert.Greater(float32(7.0), float32(len(simulateWords)))
	}
}

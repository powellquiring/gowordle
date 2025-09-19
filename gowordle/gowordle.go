package gowordle

import (
	"fmt"
	"sort"

	"github.com/bits-and-blooms/bitset"
	mapset "github.com/deckarep/golang-set"
)

var s mapset.Set = nil

type WordleWord []rune

// WordleWordMap represents a map with ordered integer keys and values
// that are a slice of WordleWord.
type WordleWordMap struct {
	keys   []int
	values map[int][]WordleWord
}

// NewWordleWordMap creates and initializes a new WordleWordMap.
func NewWordleWordMap() *WordleWordMap {
	return &WordleWordMap{
		keys:   make([]int, 0),
		values: make(map[int][]WordleWord),
	}
}

// Set adds or updates a key-value pair.
// It appends a WordleWord to the slice associated with the given key.
func (om *WordleWordMap) Set(key int, word WordleWord) {
	if _, ok := om.values[key]; !ok {
		// Key does not exist, so add it to the keys slice.
		om.keys = append(om.keys, key)
		om.values[key] = make([]WordleWord, 0)
	}
	// Append the new word to the existing slice.
	om.values[key] = append(om.values[key], word)
}

// Get retrieves the slice of WordleWord for a given key.
func (om *WordleWordMap) Get(key int) ([]WordleWord, bool) {
	value, ok := om.values[key]
	return value, ok
}

// Len returns the number of unique keys in the map.
func (om *WordleWordMap) Len() int {
	return len(om.keys)
}

// Keys returns the ordered slice of keys.
func (om *WordleWordMap) Keys() []int {
	return om.keys
}

// Iterate provides a way to loop through the key-value pairs in order.
func (om *WordleWordMap) Iterate(f func(key int, values []WordleWord)) {
	for _, key := range om.keys {
		if values, ok := om.values[key]; ok {
			f(key, values)
		}
	}
}

// SortKeys sorts the keys in ascending order.
func (om *WordleWordMap) SortKeys() {
	sort.Ints(om.keys)
}

func wwsToString(ww []WordleWord) string {
	ret := ""
	sep := ""
	for _, w := range ww {
		ret = ret + sep + string(w)
		sep = ","
	}
	return ret
}

func StringsToWordleWords(words []string) []WordleWord {
	ret := make([]WordleWord, 0, len(words))
	for _, word := range words {
		rune_word := []rune(word)
		if len(rune_word) != 5 {
			panic("not 5 letter word:" + word)
		}
		ret = append(ret, rune_word)
	}
	return ret
}

func WordleWordsToStrings(words []WordleWord) []string {
	ret := make([]string, 0, len(words))
	for _, word := range words {
		ret = append(ret, string(word))
	}
	return ret
}

func PrintWords(words []WordleWord) {
	for _, word := range words {
		fmt.Println(string(word))
	}
}

/*
letters['a'][0] all words whose first letter is an a, [1] second letter is an a, ...

a word is represented by it's index into words
*/
type WordleMatcher struct {
	words   []WordleWord
	letters [5]map[rune]*bitset.BitSet // letters[0]['a'] set of words with first letter 'a'
	count   map[rune][]*bitset.BitSet  // count['a'][0] set of words with 1 or more a, count['b'][1] words with 2 or more b
}

// take a slice of strings and make wordle words
func NewWordleMatcher(words []WordleWord) *WordleMatcher {
	ret := WordleMatcher{}
	ret.words = words
	ret.count = make(map[rune][]*bitset.BitSet, 26)
	for w, word := range words {
		word_letters := make(map[rune]int, 5)
		for l, letter := range word {
			// letters
			if ret.letters[l] == nil {
				ret.letters[l] = make(map[rune]*bitset.BitSet)
			}
			if _, ok := ret.letters[l][letter]; !ok {
				ret.letters[l][letter] = bitset.New(uint(len(words)))
			}
			ret.letters[l][letter].Set(uint(w))
			word_letters[letter] = word_letters[letter] + 1
		}
		// count
		for letter, count := range word_letters {
			for c := 0; c < count; c++ {
				if ret.count[letter] == nil {
					ret.count[letter] = make([]*bitset.BitSet, 1) // [0]
					ret.count[letter][0] = bitset.New(uint(len(words)))
				} else if len(ret.count[letter]) <= count {
					ret.count[letter] = append(ret.count[letter], bitset.New(uint(len(words))))
				}
				ret.count[letter][c].Set(uint(w))
			}
		}
	}
	return &ret
}

type LetterCount struct {
	letter rune
	count  int
}

type LetterMatch struct {
	must     map[rune]int // only consider words with this many (or more) of the letter, 0 means 1 or more
	must_not map[rune]int // eliminate all words with this many (or more) of the letter, 0 means 1 or more
}

func MakeLetterMatch(guess, answer WordleWord) LetterMatch {
	ret := LetterMatch{}
	yellow_green := make(map[rune]int, 5)
	ret.must_not = make(map[rune]int, 5)
	ret.must = make(map[rune]int, 5)
	for index, letter := range guess {
		if answer[index] == 'g' {
			yellow_green[letter] = yellow_green[letter] + 1
		} else if answer[index] == 'y' {
			yellow_green[letter] = yellow_green[letter] + 1
			ret.must[letter] = ret.must[letter] + 1
		} else { // r
			ret.must_not[letter] = 0
		}
	}
	// The number of red letters not found in the word depends on how many green/yellow
	// aaabb/ryggg means that that all words with 3 or more a's can be eliminated
	for red, _ := range ret.must_not {
		ret.must_not[red] = yellow_green[red]
	}

	// The number of yellow letters that must be in the word, more is good
	// aaabb/ygggg menas that there must be 3 a's in the word
	for yellow, _ := range ret.must {
		ret.must[yellow] = yellow_green[yellow] - 1 // 0 is 1 or more letter, 1 is 2 or more, ....
	}
	return ret
}
func MakeLetterMatch2(guess, answer WordleWord) (must, mustNot []LetterCount) {
	ret := MakeLetterMatch(guess, answer)
	retMust := []LetterCount{}
	retMustNot := []LetterCount{}
	for letter, count := range ret.must {
		retMust = append(retMust, LetterCount{letter, count})
	}
	for letter, count := range ret.must_not {
		retMustNot = append(retMustNot, LetterCount{letter, count})
	}
	return retMust, retMustNot
}

type Answer struct {
	guess   WordleWord
	Colors  WordleWord
	must    []LetterCount
	mustNot []LetterCount
}

var Hitmiss map[string]Answer = make(map[string]Answer, 10000)
var HitCount int
var MissCount int

func WordleAnswer2(solution, guess WordleWord) Answer {
	key := string(solution) + string(guess)
	if ret, ok := Hitmiss[key]; ok {
		HitCount++
		return ret
	}

	ret := Answer{
		guess: guess,
		// must:    []LetterCount{},
		// mustNot: []LetterCount{},
		must:    make([]LetterCount, 0, 5),
		mustNot: make([]LetterCount, 0, 5),
		Colors:  []rune{'r', 'r', 'r', 'r', 'r'},
	}
	solutionNotGreenCount := [26]int{}
	guessYellowGreenCount := [26]int{}
	must := [26]bool{}
	mustNot := [26]bool{}
	for i, solutionLetter := range solution {
		guessLetter := guess[i]
		if solutionLetter == guessLetter {
			ret.Colors[i] = 'g'
			guessYellowGreenCount[guessLetter-'a'] = guessYellowGreenCount[guessLetter-'a'] + 1
		} else {
			// answer[i] = 'r'
			solutionNotGreenCount[solutionLetter-'a'] = solutionNotGreenCount[solutionLetter-'a'] + 1
		}
	}
	// turn the red to yellow if in the word but not green
	for i, guessLetter := range guess {
		if ret.Colors[i] == 'r' {
			if solutionNotGreenCount[guessLetter-'a'] > 0 {
				ret.Colors[i] = 'y'
				solutionNotGreenCount[guessLetter-'a'] = solutionNotGreenCount[guessLetter-'a'] - 1
				guessYellowGreenCount[guessLetter-'a'] = guessYellowGreenCount[guessLetter-'a'] + 1
			}
		}
	}
	for i, guessLetter := range guess {
		if ret.Colors[i] == 'r' {
			if !mustNot[guessLetter-'a'] {
				ret.mustNot = append(ret.mustNot, LetterCount{guessLetter, guessYellowGreenCount[guessLetter-'a']})
				mustNot[guessLetter-'a'] = true
			}
		} else if ret.Colors[i] == 'y' {
			if !must[guessLetter-'a'] {
				// add one for each red letter
				ret.must = append(ret.must, LetterCount{guessLetter, guessYellowGreenCount[guessLetter-'a'] - 1})
				must[guessLetter-'a'] = true
			}
		}
	}
	MissCount++
	Hitmiss[key] = ret
	return ret
}

// new try
// Matching returns the set of Matching words from the game's dictionary
func (wd *WordleMatcher) Matching(guess, answer WordleWord) []WordleWord {
	must, must_not := MakeLetterMatch2(guess, answer)
	return wd.matchingWorker(guess, answer, must, must_not)
}

func (wd *WordleMatcher) Matching2(answer Answer) []WordleWord {
	return wd.matchingWorker(answer.guess, answer.Colors, answer.must, answer.mustNot)
}

//var Compliment []uint64

// [0] is BitSet for 1 bit.  Index off by 1
var bitsetAllSetPreAllocated []*bitset.BitSet = make([]*bitset.BitSet, 0)

// Length is 1..N
func NewBitsetAllSet(length int) *bitset.BitSet {
	if length < 1 {
		panic("bad length")
	}
	for i := len(bitsetAllSetPreAllocated); i < length; i++ {
		bitsetAllSetPreAllocated = append(bitsetAllSetPreAllocated, bitset.New(uint(i+1)).Complement())
	}
	set := make([]uint64, ((length-1)/64)+1)

	copy(set, bitsetAllSetPreAllocated[length-1].Bytes())
	ret := bitset.FromWithLength(uint(length), set)
	return ret
}

func (wd *WordleMatcher) matchingWorker(guess, answer WordleWord, must, must_not []LetterCount) []WordleWord {
	if len(guess) != 5 {
		panic("not 5 letter word:" + string(guess))
	}
	if len(answer) != 5 {
		panic("not 5 letter word:" + string(answer))
	}
	ret := NewBitsetAllSet(len(wd.words))
	// if there are greens then the starting point only contains words with matching letter
	for i, color := range answer {
		if color == 'g' {
			set := wd.letters[i][guess[i]]
			ret.InPlaceIntersection(set)
		}
	}

	// must letter is for yellow letters.  It indicates how many of these letters
	// must be in the word
	for _, letterCount := range must {
		yellow := letterCount.letter
		count := letterCount.count
		if counts, ok := wd.count[yellow]; ok {
			if len(counts) > count {
				set := counts[count]
				ret.InPlaceIntersection(set)
			}
		}
	}

	// red letters removes words that do not contain the required count of matching letters
	for _, letterCount := range must_not {
		red := letterCount.letter
		count := letterCount.count
		if counts, ok := wd.count[red]; ok {
			if len(counts) > count {
				set := counts[count]
				ret.InPlaceDifference(set)
			}
		}
	}

	// if there are yellow remove the words with matching letters - those would have been green
	// also remove any words that have the red letter in the same index
	for l, color := range answer {
		if color == 'y' {
			// words may not exist with this letter
			if wd.letters[l] != nil {
				if set, ok := wd.letters[l][guess[l]]; ok {
					ret.InPlaceDifference(set)
				}
			}
		}
		if color == 'r' {
			// words may not exist with this letter
			if wd.letters[l] != nil {
				if set, ok := wd.letters[l][guess[l]]; ok {
					ret.InPlaceDifference(set)
				}
			}
		}
	}
	indices := make([]uint, ret.Count())
	ret.NextSetMany(0, indices)
	retSlice := make([]WordleWord, len(indices))
	for i, index := range indices {
		retSlice[i] = wd.words[index]
	}
	return retSlice
}

func (wd *WordleMatcher) matchingWords(guess, answer string) []string {
	return []string{"todo", "todo2"}
}

// return the wordle answer for the quess given the solution
func WordleAnswer(solution, guess WordleWord) WordleWord {
	answer := WordleAnswer2(solution, guess)
	return answer.Colors
}

func WordleAnswerOrig(solution, guess WordleWord) WordleWord {
	answer := make([]rune, 5)
	solution_not_green := make(map[rune]int, 5)
	for i, letter := range solution {
		if letter == guess[i] {
			answer[i] = 'g'
		} else {
			answer[i] = 'r'
			solution_not_green[letter] = solution_not_green[letter] + 1
		}
	}
	// turn the red to yellow if in the word but not green
	for i, letter := range guess {
		if answer[i] == 'r' {
			if solution_not_green[letter] > 0 {
				answer[i] = 'y'
				solution_not_green[letter] = solution_not_green[letter] - 1
			}
		}
	}
	return answer
}

type GuessAnswer struct {
	Guess  WordleWord
	Answer WordleWord
}

// play wordle against the computer providing the current board state
// return the next best answer
func PlayWorldReturnPossible(allWordleWords []WordleWord, guessAnswers []GuessAnswer) (WordleWord, []WordleWord) {
	possibleAnswers := allWordleWords

	for _, guessAnswer := range guessAnswers {
		game := NewWordleMatcher(possibleAnswers)
		possibleAnswers = game.Matching([]rune(guessAnswer.Guess), []rune(guessAnswer.Answer))
	}
	//ret := NextGuess(allWordleWords, possibleAnswers)
	ret := NextGuess1(allWordleWords, possibleAnswers)
	return ret, possibleAnswers
}

func PlayWordle(allWordleWords []WordleWord, guessAnswers []GuessAnswer) WordleWord {
	ret, _ := PlayWorldReturnPossible(allWordleWords, guessAnswers)
	return ret
}

func NextGuess1(allWords, possibleAnswers []WordleWord) WordleWord {
	_, wordsPossible := BestGuess1(allWords, possibleAnswers, possibleAnswers, 1, len(possibleAnswers)+1)
	return wordsPossible[0]
}

func FirstGuess1(allWords []string) (float32, []WordleWord) {
	wws := StringsToWordleWords(allWords)
	score, ret := BestGuess1(wws, wws, wws, 1, len(allWords))
	return float32(score), ret
}

func FirstGuessProvideInitialGuesses1(initialGuesses_s, allWords_s []string) (float32, []WordleWord) {
	allWords := StringsToWordleWords(allWords_s)
	initialGuesses := StringsToWordleWords(initialGuesses_s)
	// score, ret := BestGuess1(allWords, allWords, initialGuesses, 1, len(allwords))
	score, ret := BestGuess1(allWords, allWords, initialGuesses, 1, 10)
	return float32(score), ret
}

type WordFloat struct {
	word WordleWord
	flt  float32
}
type WordFloatByFloat []WordFloat

func (wf WordFloatByFloat) Len() int           { return len(wf) }
func (wf WordFloatByFloat) Swap(i, j int)      { wf[i], wf[j] = wf[j], wf[i] }
func (wf WordFloatByFloat) Less(i, j int) bool { return wf[i].flt < wf[j].flt }

// configurable from command line
var RECURSIVE bool = true

var matching2 bool = true
var Logging bool = false
var BetterGuesses map[string]int = make(map[string]int)

// find best next guess, return the low score and the slice of words that have that score
// The score will be the average number of guesses it will take to solve if one the best guesses is used
func ScoreAlgorithmRecursive(allWords, possibleWords, initialGuesses []WordleWord, depth int, bestScoreSoFar int) (int, []WordleWord) {
	const INIFINITY_SCORE = 1000000
	if depth > 6 {
		// this is the path for a bad guess, eject
		return INIFINITY_SCORE, []WordleWord{WordleWord("BADWD")}
	}
	if len(possibleWords) == 0 {
		panic("possibleWords is empty")
	}
	if len(possibleWords) == 1 {
		return 100, possibleWords // just guess it
	}
	if len(possibleWords) == 2 {
		// if there are two words choose either of the words and the guesses will be 1 if the right guess and 2 if the wrong guess
		return 150, possibleWords
	}
	// Using the possible words the best guess is the matching one for one solution and 2 for the rest of the solutions.
	game := NewWordleMatcher(possibleWords)
	bestScore := 1000
	var bestGuess []WordleWord
	possibleWordsSet := make(map[string]bool)
	for _, guess := range possibleWords {
		score := 0 // running average
		for count, solution := range possibleWords {
			matching := game.Matching2(WordleAnswer2(solution, guess))
			subscore, _ := ScoreAlgorithmRecursive(allWords, matching, matching, depth+1, bestScoreSoFar)
			if !((len(matching) == 1) && (string(matching[0]) == string(guess))) {
				subscore += 100 // if the guess is the solution then the subscore will be 100, otherwise it will be the current guess pluss the recursive guess
			}
			score = score + ((subscore - score) / (count + 1))
		}
		// all but one of the guesses is incorrect (hence the -1) then each of the second guesses takes only 1 guess
		if score <= (((len(possibleWords)-1)+len(possibleWords))*100)/len(possibleWords) {
			// found this guess compared to all the possible words returns the correct answer (1) or narrows down to a single answer for the correct guess
			return score, []WordleWord{guess}
		}
		possibleWordsSet[string(guess)] = true
		// although not best possible it could still be the best so far
		if score < bestScore {
			bestScore = score
			bestGuess = []WordleWord{guess}
		} else if score == bestScore {
			bestGuess = append(bestGuess, guess)
		}
	}

	// the best that can be done by using a non possible word is finding a guess that narrows it down to 1 in all cases, thus taking 2 guesses total
	// so if that has already been achieved use the best guess
	if bestScore <= 200 {
		return bestScore, bestGuess
	}

	// did not find an optimal answer
	for _, guess := range allWords {
		if _, ok := possibleWordsSet[string(guess)]; ok {
			// already tried this guess
			continue
		}
		score := 0 // running average
		for count, solution := range possibleWords {
			matching := game.Matching2(WordleAnswer2(solution, guess))
			if len(matching) == len(possibleWords) {
				// not narrowing it down any this solution so it is a bad guess, go to next guess
				score = INIFINITY_SCORE
				break
			}
			subscore, _ := ScoreAlgorithmRecursive(allWords, matching, matching, depth+1, bestScoreSoFar)
			subscore += 100 // current score plus the recursive score
			score = score + ((subscore - score) / (count + 1))

			// 200 is the best for the remaining words, if the current average plus best possible result for the remaining words
			// is alread over that bail
			bestPossibleScore := ((score * (count + 1)) + (200 * (len(possibleWords) - (count + 1)))) / len(possibleWords)
			if bestPossibleScore > bestScore {
				score = bestPossibleScore
				break
			}
		}

		if score < bestScore {
			bestScore = score
			bestGuess = []WordleWord{guess}
		} else if score == bestScore {
			bestGuess = append(bestGuess, guess)
		}
	}
	return bestScore, bestGuess
}

/******
	// Did not get the optimal score.  The best possible to 2 * len(possibleWords)
	if bestScore == 2 * len(possibleWords) {
		return bestScore, bestGuess
	}

	// Painfull - need to go through all possible wors

		if RECURSIVE && depth < 2 {
			// get a more accurage score of this guess by looking at all the possible next guesses
			for _, solutionInner := range matching {
				// for each possible solution look at the produced subset of words
				// and guess again.
				gameInner := NewWordleMatcher(matching)
				matchingInner := gameInner.matching2(WordleAnswer2(solutionInner, guess))
				scores := ScoreAlgorithmTotalMatches1LevelAll(allWords, matchingInner, matchingInner, depth+1, 0)
				score += scores.Keys()[0]
			}
		} else {
			score += len(matching)
		}
	}
	return score
	for _, guess := range possibleWords {
		score := 0
		for _, solution := range possibleWords {
			answer := WordleAnswer2(solution, guess)
			matching := game.matching2(answer)
			if len(matching) == 1 {
				score += 1
			} else {
				score += 2
			}
		}

	}

				gameInner := NewWordleMatcher(matching)
				matchingInner := gameInner.matching2(WordleAnswer2(solutionInner, guess))
				scores := ScoreAlgorithmTotalMatches1LevelAll(allWords, matchingInner, matchingInner, depth+1, 0)
				score += scores.Keys()[0]



	// use possible words first - then rest of the words
	initialGuessMap := make(map[string]bool, len(initialGuesses))
	orderedGuesses := make([]WordleWord, len(initialGuesses))
	copy(orderedGuesses, initialGuesses)

	for _, guess := range initialGuesses {
		initialGuessMap[string(guess)] = true
	}
	for _, guess := range allWords {
		if _, ok := initialGuessMap[string(guess)]; !ok {
			orderedGuesses = append(orderedGuesses, guess)
		}
	}
	ret := NewWordleWordMap()
	for _, guess := range orderedGuesses {
		score := GuessScore(guess, possibleWords, allWords, depth)
		ret.Set(score, guess)
	}
	ret.SortKeys()
	return ret
	keys := wwm.Keys()
	values, _ := wwm.Get(keys[0])
	return keys[0], values
}
// first try
***/

func ScoreAlgorithmTotalMatches1Level(allWords, possibleWords, initialGuesses []WordleWord, depth int, bestScoreSoFar int) (int, []WordleWord) {
	wwm := ScoreAlgorithmTotalMatches1LevelAll(allWords, possibleWords, initialGuesses, depth, bestScoreSoFar)
	keys := wwm.Keys()
	values, _ := wwm.Get(keys[0])
	return keys[0], values
}

// total number of words
func GuessScore(guess WordleWord, possibleWords []WordleWord, allWords []WordleWord, depth int) int {
	game := NewWordleMatcher(possibleWords)
	score := 0
	for _, solution := range possibleWords {
		matching := game.Matching2(WordleAnswer2(solution, guess))
		if RECURSIVE && depth < 2 {
			// get a more accurage score of this guess by looking at all the possible next guesses
			for _, solutionInner := range matching {
				// for each possible solution look at the produced subset of words
				// and guess again.
				gameInner := NewWordleMatcher(matching)
				matchingInner := gameInner.Matching2(WordleAnswer2(solutionInner, guess))
				scores := ScoreAlgorithmTotalMatches1LevelAll(allWords, matchingInner, matchingInner, depth+1, 0)
				score += scores.Keys()[0]
			}
		} else {
			score += len(matching)
		}
	}
	return score
}

// try all the guesses and return a map of score to guess.
func ScoreAlgorithmTotalMatches1LevelAll(allWords, possibleWords, initialGuesses []WordleWord, depth int, bestScoreSoFar int) *WordleWordMap {
	if len(possibleWords) == 0 {
		panic("possibleWords is empty")
	}
	if len(possibleWords) == 1 {
		ret := NewWordleWordMap()
		ret.Set(1, possibleWords[0])
		return ret
	}

	// use possible words first - then rest of the words
	initialGuessMap := make(map[string]bool, len(initialGuesses))
	orderedGuesses := make([]WordleWord, len(initialGuesses))
	copy(orderedGuesses, initialGuesses)

	for _, guess := range initialGuesses {
		initialGuessMap[string(guess)] = true
	}
	for _, guess := range allWords {
		if _, ok := initialGuessMap[string(guess)]; !ok {
			orderedGuesses = append(orderedGuesses, guess)
		}
	}
	ret := NewWordleWordMap()
	for _, guess := range orderedGuesses {
		score := GuessScore(guess, possibleWords, allWords, depth)
		ret.Set(score, guess)
	}
	ret.SortKeys()
	return ret
}

var BestGuess1 func([]WordleWord, []WordleWord, []WordleWord, int, int) (int, []WordleWord) = ScoreAlgorithmTotalMatches1Level

// Simulate a game of wordle.
// words_s - dictionary of words
// solution - answer
// first_guess - first guess
func Simulate(words_s []string, solution_s string, first_guess_s string) []string {
	words := StringsToWordleWords(words_s)
	solution := []rune(solution_s)
	guess := []rune(first_guess_s)
	guesses := []string{}
	gas := make([]GuessAnswer, 0)
	for guessCount := 0; guessCount < 6; guessCount++ {
		guesses = append(guesses, string(guess))
		answer := WordleAnswer(solution, guess)
		if string(answer) == "ggggg" {
			return guesses
		}
		gas = append(gas, GuessAnswer{guess, WordleAnswer(solution, guess)})
		guess = PlayWordle(words, gas)
	}
	panic("unexpected Simulate end")
}

type SolutionsAnswers struct {
	Solutions    []string
	AnswerColors []string
}

// for the given guess return a map key = matching possible solutions, value = slice of Solutions and answerColors
func UniqueGuessResults(wordListStrings []string, guess string) map[string]SolutionsAnswers {
	guessWW := []rune(guess)
	wordList := StringsToWordleWords(wordListStrings)
	game := NewWordleMatcher(wordList)
	sortedSolutions := make(map[string]SolutionsAnswers)
	for _, solution := range wordList {
		answer := WordleAnswer2(solution, guessWW)
		possibleSolutions := game.Matching(guessWW, answer.Colors)
		sort.Slice(possibleSolutions, func(i, j int) bool {
			return string(possibleSolutions[i]) < string(possibleSolutions[j])
		})
		allSolutions := ""
		for _, solution := range possibleSolutions {
			allSolutions += string(solution) + " "
		}
		if solutionAnswers, ok := sortedSolutions[allSolutions]; ok {
			newSolutionAnswers := sortedSolutions[allSolutions]
			newSolutionAnswers.Solutions = append(solutionAnswers.Solutions, string(solution))
			if solutionAnswers.AnswerColors[0] != string(answer.Colors) {
				newSolutionAnswers.AnswerColors = append(solutionAnswers.AnswerColors, "BUG-"+string(solution)+"-"+string(answer.Colors))
			}
			sortedSolutions[allSolutions] = newSolutionAnswers
		} else {
			sortedSolutions[allSolutions] = SolutionsAnswers{Solutions: []string{string(solution)}, AnswerColors: []string{string(answer.Colors)}}
		}
	}
	return sortedSolutions
}

// for the given guess return a map key = answer colors value = matching solutions
func UniqueAnswerResults(wordListStrings []string, guess string) map[string][]string {
	wordList := StringsToWordleWords(wordListStrings)
	guessWW := []rune(guess)
	answerSolutions := make(map[string][]string)
	for _, solutionWW := range wordList {
		answer := WordleAnswer2(solutionWW, guessWW)
		answerColors := string(answer.Colors)
		if answers, ok := answerSolutions[string(answerColors)]; ok {
			answerSolutions[answerColors] = append(answers, string(solutionWW))
		} else {
			answerSolutions[answerColors] = []string{string(solutionWW)}
		}
	}
	return answerSolutions
}

package gowordle

import (
	"fmt"
	"sort"

	"github.com/bits-and-blooms/bitset"
	mapset "github.com/deckarep/golang-set"
)

var s mapset.Set = nil

type WordleWord []rune

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
		println(string(word))
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

////////////////////////////////////////////////////
type RuneCount struct {
	l     rune
	count int
}
type LetterMatch2 struct {
	must     []RuneCount // only consider words with this many (or more) of the letter, 0 means 1 or more
	must_not []RuneCount // eliminate all words with this many (or more) of the letter, 0 means 1 or more
}

//
func MakeLetterMatch2(guess, answer WordleWord) LetterMatch {
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

// matching returns the set of matching words from the game's dictionary
func (wd *WordleMatcher) matching(guess, answer WordleWord) []WordleWord {
	if len(guess) != 5 {
		panic("not 5 letter word:" + string(guess))
	}
	if len(answer) != 5 {
		panic("not 5 letter word:" + string(answer))
	}
	ret := bitset.New(uint(len(wd.words))).Complement()
	// if there are greens then the starting point only contains words with matching letter
	for i, color := range answer {
		if color == 'g' {
			set := wd.letters[i][guess[i]]
			ret = ret.Intersection(set)
		}
	}

	letterMatch := MakeLetterMatch(guess, answer)

	// must letter is for yellow letters.  It indicates how many of these letters
	// must be in the word
	for yellow, count := range letterMatch.must {
		if counts, ok := wd.count[yellow]; ok {
			if len(counts) > count {
				set := counts[count]
				ret = ret.Intersection(set)
			}
		}
	}

	// red letters removes words that do not contain the required count of matching letters
	for red, count := range letterMatch.must_not {
		if counts, ok := wd.count[red]; ok {
			if len(counts) > count {
				set := counts[count]
				ret = ret.Difference(set)
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
					ret = ret.Difference(set)
				}
			}
		}
		if color == 'r' {
			// words may not exist with this letter
			if wd.letters[l] != nil {
				if set, ok := wd.letters[l][guess[l]]; ok {
					ret = ret.Difference(set)
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

const maxDepth = 4
const spaces = "         "

type GuessAnswer struct {
	Guess  WordleWord
	Answer WordleWord
}

type NG func(allWords, possibleAnswers []WordleWord) WordleWord

var nextGuess NG = NextGuess
var nextGuess1 NG = NextGuess1

// play wordle against the computer providing the current board state
// return the next best answer
func playWordle(allWordleWords []WordleWord, guessAnswers []GuessAnswer) WordleWord {
	possibleAnswers := allWordleWords

	for _, guessAnswer := range guessAnswers {
		game := NewWordleMatcher(possibleAnswers)
		possibleAnswers = game.matching([]rune(guessAnswer.Guess), []rune(guessAnswer.Answer))
	}
	//ret := NextGuess(allWordleWords, possibleAnswers)
	ret := nextGuess1(allWordleWords, possibleAnswers)
	return ret
}

// return the next guess given all words as well as the subset of all words
// which are the possible answers
func NextGuess(allWords, possibleAnswers []WordleWord) WordleWord {
	scoreAll, wordsAll := BestGuess(allWords, possibleAnswers, 1)
	scorePossible, wordsPossible := BestGuess(possibleAnswers, possibleAnswers, 1)
	if scorePossible <= scoreAll {
		return wordsPossible[0]
	} else {
		return wordsAll[0]
	}
}
func NextGuess1(allWords, possibleAnswers []WordleWord) WordleWord {
	_, wordsPossible := BestGuess1(allWords, possibleAnswers, allWords, 1, len(possibleAnswers)+1)
	return wordsPossible[0]
}

func FirstGuess(allWords []string) (float32, string) {
	wws := StringsToWordleWords(allWords)
	score, ret := BestGuess(wws, wws, 1)
	return score, string(ret[0])
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
	score, ret := BestGuess1(allWords, allWords, initialGuesses, 1, 6)
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

// return the next guess given all words as well as the subset of all words
// which are the possible answers
func BestGuess(allWords, possibleWords []WordleWord, depth int) (float32, []WordleWord) {
	game := NewWordleMatcher(possibleWords)
	guessAverageMatching := make(WordFloatByFloat, 0, len(possibleWords))
	for _, guess := range allWords {
		totalMatches := 0
		for _, solution := range possibleWords {
			answer := WordleAnswer(solution, guess)
			matching := game.matching(guess, answer)
			totalMatches = totalMatches + len(matching)
		}
		guessAverageMatching = append(guessAverageMatching, WordFloat{guess, float32(totalMatches) / float32(len(possibleWords))})
	}
	return AnalyzeGuesses(guessAverageMatching)
}

func AnalyzeGuesses(guessAverageMatching WordFloatByFloat) (float32, []WordleWord) {
	sort.Sort(guessAverageMatching)
	/*
		for i, guessAvg := range guessAverageMatching {
			if i > 10 {
				break
			}
			fmt.Printf("%f, %s\n", guessAvg.flt, string(guessAvg.word))
		}
	*/
	return guessAverageMatching[0].flt, []WordleWord{guessAverageMatching[0].word}
}

func BestGuess1(allWords, possibleWords, initialGuesses []WordleWord, depth int, bestScoreSoFar int) (int, []WordleWord) {
	if depth > 10 {
		panic("BestGuess1 too deep")
	}
	if len(possibleWords) == 0 {
		panic("possibleWords is empty")
	}
	if len(possibleWords) == 1 {
		return depth, possibleWords
	}
	guessesWithBestScore := []WordleWord{}
	game := NewWordleMatcher(possibleWords)
	scores := make(map[int][]WordleWord)
	for guessNumber, guess := range initialGuesses {
		worstScoreForGuess := 0
		for solutionNumber, solution := range possibleWords {
			if depth == 1 && ((solutionNumber % 10) == 0) {
				println("guess/solution: ", guessNumber, "/", solutionNumber)
			}
			score := 0
			if string(guess) == string(solution) {
				score = depth
			} else {
				answer := WordleAnswer(solution, guess)
				matching := game.matching(guess, answer)
				if len(matching) == len(possibleWords) {
					score = len(matching) + depth
				} else if len(matching) == 2 {
					score = 2 + depth
				} else if len(matching) == 1 {
					score = 1 + depth
				} else if len(matching) == 0 {
					fmt.Println("***bug")
				} else {
					if depth+2 >= bestScoreSoFar {
						worstScoreForGuess = depth + 3 // could be equal, but not sure, add an extra (3 instead of 2)
						break                          // short circuit, no need to look for even worse scores
					}
					score, _ = BestGuess1(allWords, matching, allWords, depth+1, bestScoreSoFar)
				}
			}
			if score > worstScoreForGuess {
				worstScoreForGuess = score
			}
			if worstScoreForGuess >= bestScoreSoFar {
				// short circuit, no need to look for even worse scores
				worstScoreForGuess++ // it may be the same, not worse, oh well
				break
			}
		}
		if _, ok := scores[worstScoreForGuess]; !ok {
			scores[worstScoreForGuess] = make([]WordleWord, 0)
		}
		if worstScoreForGuess < bestScoreSoFar {
			bestScoreSoFar = worstScoreForGuess
			guessesWithBestScore = []WordleWord{guess}
			if depth == 1 {
				println("Better guess/score:", string(guess), "/", bestScoreSoFar)
			}
		} else if worstScoreForGuess == bestScoreSoFar {
			guessesWithBestScore = append(guessesWithBestScore, guess)
		} else {
			if depth == 1 {
				println("no better guess/score", string(guess), "/", worstScoreForGuess)
			}
		}
		scores[worstScoreForGuess] = append(scores[worstScoreForGuess], guess)
	}
	return bestScoreSoFar, guessesWithBestScore
	/*
		// Sort the keys
		scoreKeys := make([]int, len(scores))
		i := 0
		for k, _ := range scores {
			scoreKeys[i] = k
			i++
		}
		sort.Ints(scoreKeys)
		return scoreKeys[0], scores[scoreKeys[0]]
	*/
}

//Simulate a game of wordle.
//words_s - dictionary of words
//solution - answer
//first_guess - first guess
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
		guess = playWordle(words, gas)
	}
	panic("unexpected Simulate end")
}

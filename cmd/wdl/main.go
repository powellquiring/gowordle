package main

import (
	"context"
	"log"
	"os"
	"sort"

	"github.com/powellquiring/gowordle/gowordle"
	"github.com/urfave/cli/v3" // imports as package "cli"
)

func FirstWithInitialGuesses(wordCount int) {
	if wordCount == 0 {
		wordCount = len(gowordle.WordleDictionary)
	}
	wordList := gowordle.WordleDictionary[0:wordCount]
	score, words := gowordle.FirstGuessProvideInitialGuesses1(wordList, wordList)
	print(score)
	gowordle.PrintWords(words)
}

func FirstWords(wordCount int) {
	if wordCount == 0 {
		wordCount = len(gowordle.WordleDictionary)
	}
	wordList := gowordle.WordleDictionary[0:wordCount]
	wws := gowordle.StringsToWordleWords(wordList)
	ret := gowordle.ScoreAlgorithmTotalMatches1LevelAll(wws, wws, wws, 0, len(wordList))
	for keyCount, key := range ret.Keys() {
		values, _ := ret.Get(key)
		print(key, " ")
		for _, value := range values {
			print(string(value), " ")
		}
		println()
		if keyCount > 10 {
			break
		}
	}
}

const DefaultFirstWord = "raise"

func simulate(wordCount int, answers []string) {
	if wordCount == 0 {
		wordCount = len(gowordle.WordleDictionary)
	}
	wordList := gowordle.WordleDictionary[0:wordCount]
	if len(answers) == 0 {
		answers = wordList
	}
	type Game struct {
		Answer  string
		Guesses []string
	}
	sortedGames := make(map[int][]Game)
	for _, answer := range answers {
		guesses := gowordle.Simulate(wordList, answer, DefaultFirstWord)
		if _, ok := sortedGames[len(guesses)]; !ok {
			sortedGames[len(guesses)] = make([]Game, 0)
		}
		sortedGames[len(guesses)] = append(sortedGames[len(guesses)], Game{answer, guesses})
	}

	// create slice of number of guesses
	keys := make([]int, 0, len(sortedGames))
	for k := range sortedGames {
		keys = append(keys, k)
	}
	// Sort the slice of keys
	sort.Ints(keys)

	for _, numGuesses := range keys {
		games := sortedGames[numGuesses]
		println(numGuesses, len(games), " ---------------------")
		for _, game := range games {
			print(game.Answer, ":")
			for _, guess := range game.Guesses {
				print(" ", guess)
			}
			println()
		}
	}
}

// playWordle with guess/answer pairs provided
func playWordle(wordCount int, answers []string) {
	if wordCount == 0 {
		wordCount = len(gowordle.WordleDictionary)
	}
	wordList := gowordle.StringsToWordleWords(gowordle.WordleDictionary[0:wordCount])

	gas := make([]gowordle.GuessAnswer, 0)
	for i := 0; i < len(answers); i += 2 {
		guess := answers[i]
		answer := answers[i+1]
		gas = append(gas, gowordle.GuessAnswer{Guess: gowordle.WordleWord(guess), Answer: gowordle.WordleWord(answer)})
	}
	nextGuess, possible := gowordle.PlayWorldReturnPossible(wordList, gas)
	print(string(nextGuess), ":")
	for _, word := range possible {
		print(" ", string(word))
	}
	println()
}

func main() {
	count := 0
	cmd := &cli.Command{
		Name:  "wdl",
		Usage: "wordle",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "count",
				Value:       0,
				Aliases:     []string{"c"},
				Usage:       "number of words, 0 is all words",
				Destination: &count,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "first",
				Usage: "first guess",
				Action: func(context.Context, *cli.Command) error {
					FirstWords(count)
					return nil
				},
			},
			{
				Name: "sim",
				Usage: `sim [answer] ...
				Simulate multiple games by specifying a list of answers for each game.  If no answers are provided,
				simulate all words.  All words can be cut back by using the -count flag.
				`,
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.NArg() == 0 {
						simulate(count, []string{})
					} else {
						simulate(count, cmd.Args().Slice())
					}
					return nil
				},
			},
			{
				Name:  "play",
				Usage: "play a game of wordle by entering pairs of [guess answer]...",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.NArg()%2 != 0 {
						return cli.Exit("must have pairs of guess answer", 1)
					} else if cmd.NArg() < 2 {
						return cli.Exit("must have at least one guess answer", 2)
					} else {
						playWordle(count, cmd.Args().Slice())
					}
					return nil
				},
			},
			{
				Name:  "measure",
				Usage: "measure the performance of an algorithm by playing against a set of answers",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					playWordle(count, cmd.Args().Slice())
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

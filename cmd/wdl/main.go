package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/powellquiring/gowordle/gowordle"
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v3" // imports as package "cli"
)

func server(globalConfig GlobalConfiguration, solution string, guesses []string) {
	wordList := globalConfig.AllWords
	wws := gowordle.StringsToWordleWords(wordList)
	fmt.Print(solution, " ")
	solutionWW := []rune(solution)
	for _, guess := range guesses {
		game := gowordle.NewWordleMatcher(wws)
		guessWW := []rune(guess)
		answer := gowordle.WordleAnswer2(solutionWW, guessWW)
		wws = game.Matching2(answer)
		fmt.Println(guess, string(answer.Colors), gowordle.WordleWordsToStrings(wws))
	}
}

func FirstWords(globalConfig GlobalConfiguration) {
	wordList := globalConfig.AllWords
	results := gowordle.UniqueGuessResults(wordList, globalConfig.FirstWord)
	for unique, solutionAnswers := range results {
		fmt.Println(len(unique)/6, unique, solutionAnswers.AnswerColors, solutionAnswers.Solutions)
	}
}

func FirstWordsByAnswerColor(globalConfig GlobalConfiguration) {
	wordList := globalConfig.AllWords
	results := gowordle.UniqueAnswerResults(wordList, globalConfig.FirstWord)
	sortedAnswerColors := []string{}
	for answerColors, _ := range results {
		sortedAnswerColors = append(sortedAnswerColors, answerColors)
	}
	sort.Strings(sortedAnswerColors)
	fmt.Println(globalConfig.FirstWord, "--- distribution of solutions when using this guess")
	for _, answerColors := range sortedAnswerColors {
		solutions := results[answerColors]
		fmt.Println(len(solutions), answerColors, solutions)
	}
}

func FirstWords1(globalConfig GlobalConfiguration) {
	wordList := globalConfig.AllWords
	wws := gowordle.StringsToWordleWords(wordList)
	ret := gowordle.ScoreAlgorithmTotalMatches1LevelAll(wws, wws, wws, 0, len(wordList))
	for keyCount, key := range ret.Keys() {
		values, _ := ret.Get(key)
		fmt.Print(key, " ")
		for _, value := range values {
			fmt.Print(string(value), " ")
		}
		fmt.Println()
		if keyCount > 10 {
			break
		}
	}
}

func simulate(globalConfig GlobalConfiguration, answers []string) {
	wordList := globalConfig.AllWords
	if len(answers) == 0 {
		answers = wordList
	}
	type Game struct {
		Answer  string
		Guesses []string
	}
	sortedGames := make(map[int][]Game)

	var bar *progressbar.ProgressBar
	if globalConfig.progress {
		bar = progressbar.Default(int64(len(answers)))
	} else {
		bar = progressbar.DefaultSilent(int64(len(answers)))
	}

	for _, answer := range answers {
		bar.Add(1)
		guesses := gowordle.Simulate(wordList, answer, globalConfig.FirstWord)
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
		fmt.Println(numGuesses, len(games), " ---------------------")
		for _, game := range games {
			fmt.Print(game.Answer, ":")
			for _, guess := range game.Guesses {
				fmt.Print(" ", guess)
			}
			fmt.Println()
		}
	}
}

// playWordle with guess/answer pairs provided
func playWordle(globalConfig GlobalConfiguration, answers []string) {
	wordList := gowordle.StringsToWordleWords(globalConfig.AllWords)

	gas := make([]gowordle.GuessAnswer, 0)
	for i := 0; i < len(answers); i += 2 {
		guess := answers[i]
		answer := answers[i+1]
		gas = append(gas, gowordle.GuessAnswer{Guess: gowordle.WordleWord(guess), Answer: gowordle.WordleWord(answer)})
	}
	nextGuess, possible := gowordle.PlayWorldReturnPossible(wordList, gas)
	fmt.Print(string(nextGuess), ":")
	for _, word := range possible {
		fmt.Print(" ", string(word))
	}
	fmt.Println()
}

type GlobalConfiguration struct {
	AllWords  []string
	Recursive bool
	progress  bool
	FirstWord string
}

func globalCofiguration(count int, recursive bool, progress bool, firstWord string) GlobalConfiguration {
	if count == 0 {
		count = len(gowordle.WordleDictionary)
	}
	gowordle.RECURSIVE = recursive
	if recursive {
		gowordle.BestGuess1 = gowordle.ScoreAlgorithmRecursive
	}
	if firstWord == "" {
		firstWord = "raise"
	}
	return GlobalConfiguration{
		AllWords:  gowordle.WordleDictionary[0:count],
		Recursive: recursive,
		progress:  progress,
		FirstWord: firstWord,
	}

}

func main() {
	count := 0
	recursive := false
	progress := false
	firstWord := ""
	// going raise blunt
	//server(globalCofiguration(count, recursive, progress, firstWord), "going", []string{"raise", "blunt"})
	// FirstWords(globalCofiguration(count, recursive, progress, firstWord))
	// playWordle(globalCofiguration(count, true, progress, firstWord), []string{"raise", "ryyry"})
	simulate(globalCofiguration(count, true, progress, firstWord), []string{})
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
			&cli.BoolFlag{
				Name:        "recursive",
				Value:       false,
				Aliases:     []string{"r"},
				Usage:       "turn on recursive flag slower but better",
				Destination: &recursive,
			},
			&cli.BoolFlag{
				Name:        "progress",
				Value:       false,
				Aliases:     []string{"p"},
				Usage:       "show progress bar",
				Destination: &progress,
			},
			&cli.BoolFlag{
				Name:        "progress",
				Value:       false,
				Aliases:     []string{"p"},
				Usage:       "show progress bar",
				Destination: &progress,
			},
			&cli.StringFlag{
				Name:        "first",
				Value:       "",
				Aliases:     []string{"f"},
				Usage:       "first word to guess, default is 'raise', only used with sim command",
				Destination: &firstWord,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "first",
				Usage: "first guess",
				Action: func(context.Context, *cli.Command) error {
					FirstWords(globalCofiguration(count, recursive, progress, firstWord))
					return nil
				},
			},
			{
				Name:    "firstbycolor",
				Aliases: []string{"fc"},
				Usage:   "firstbycolor",
				Action: func(context.Context, *cli.Command) error {
					FirstWordsByAnswerColor(globalCofiguration(count, recursive, progress, firstWord))
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
					globalCofiguration(count, recursive, progress, firstWord)
					if cmd.NArg() == 0 {
						simulate(globalCofiguration(count, recursive, progress, firstWord), []string{})
					} else {
						simulate(globalCofiguration(count, recursive, progress, firstWord), cmd.Args().Slice())
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
						playWordle(globalCofiguration(count, recursive, progress, firstWord), cmd.Args().Slice())
					}
					return nil
				},
			},
			{
				Name: "server",
				Usage: `server solution guess...
				be a wordle server return the ryg for each guess along with the remaining words`,
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.NArg() < 2 {
						return cli.Exit("must supply both an answer and one or more guesses", 3)
					}
					args := cmd.Args().Slice()
					server(globalCofiguration(count, recursive, progress, firstWord), args[0], args[1:])
					return nil
				},
			},
			{
				Name:  "measure",
				Usage: "measure the performance of an algorithm by playing against a set of answers",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					playWordle(globalCofiguration(count, recursive, progress, firstWord), cmd.Args().Slice())
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

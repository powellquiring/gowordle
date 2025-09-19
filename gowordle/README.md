# gowordle
A few wordle related programs to come up with good initial guesses and then play today's wordle and simulate all wordle games.

Determining the algorithm for the best guess give then current information is challenging.

## Algorithm - total matches one level
A guess is rated by the possible matches. Fewer matches better guess. The score for a guess is the sum of matches for all possible answers.

## Algorithm - total matches one level
A better guess could consider guesses that could be applied to the next level.  This is the sum of the matches for all possible answers for all possible guesses.


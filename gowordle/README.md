# gowordle
A few wordle related programs to come up with good initial guesses and then play today's wordle and simulate all wordle games.

Determining the algorithm for the best guess give then current information is challenging.

## Algorithm - total matches one level
A guess is rated by the possible matches. Fewer matches better guess. The number of matches is the sum of the number of matches for each possible answer. 

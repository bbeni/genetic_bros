package main

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/bbeni/genetic_bros/game"
)

const (
	NUMBER_OF_BOTS         = 100 // has to be even too!
	NUMBER_OF_MOVES        = 1200
	NUMBER_SLICE_POSITIONS = 6 // has to be even!
	NUMBER_GENERATIONS     = 5000
)

func main() {

	for i := range NUMBER_OF_BOTS {
		bots[i] = MakeBot(69)
	}

	var kid1 Bot
	var kid2 Bot

	for generation_nr := range NUMBER_GENERATIONS {
		var scores_sum float64
		var scores [NUMBER_OF_BOTS]float64

		for i := range NUMBER_OF_BOTS {
			for _, move := range bots[i].moves {
				bots[i].Gs.Update(move)
				if bots[i].Gs.GameOver() {
					break
				}
			}
			scores[i] = evaluate(&bots[i])
			scores_sum += scores[i]

		}

		// normalize
		for i := range NUMBER_OF_BOTS {
			scores[i] /= scores_sum
		}

		cdf := get_cdf(scores)

		var children [NUMBER_OF_BOTS]Bot

		for pair := range NUMBER_OF_BOTS / 2 {
			// chose two random bots with probability proportional to scores
			mother_idx, father_idx := choose2(cdf[:])

			// the pair makes 2 chilren with crossover and random mutation
			kid1 = MakeBot(69)
			kid2 = MakeBot(69)

			// copy
			for i := range NUMBER_OF_MOVES {
				kid1.moves[i] = bots[mother_idx].moves[i]
				kid2.moves[i] = bots[father_idx].moves[i]
			}

			// crossover slices
			var slice_indices [NUMBER_SLICE_POSITIONS]int
			for i := range NUMBER_SLICE_POSITIONS {
				slice_idx := rand.Int() % NUMBER_OF_MOVES
				slice_indices[i] = slice_idx
			}
			sort.Ints(slice_indices[:])

			var t []game.Direction
			for i := range NUMBER_SLICE_POSITIONS / 2 {
				idx1 := slice_indices[i*2]
				idx2 := slice_indices[i*2+1]
				copy(t, kid1.moves[idx1:idx2])
				copy(kid1.moves[idx1:idx2], kid2.moves[idx1:idx2])
				copy(kid2.moves[idx1:idx2], t)
			}

			// do mutations on both kids
			mutation_rate := 0.0001 // 10**-4 to 10**-6 https://www.sciencedirect.com/topics/biochemistry-genetics-and-molecular-biology/mutation-rate
			for i := range NUMBER_OF_MOVES {
				r := rand.Float64()
				if r <= mutation_rate {
					kid1.moves[i] = game.Direction(rand.Int() % 4)
				}
			}

			for i := range NUMBER_OF_MOVES {
				r := rand.Float64()
				if r <= mutation_rate {
					kid2.moves[i] = game.Direction(rand.Int() % 4)
				}
			}

			children[pair*2] = kid1
			children[pair*2+1] = kid2
		}

		best_bot := find_best_bot(bots[:])
		fmt.Printf("Generation %v %v\n", generation_nr, best_bot.Gs.Step)

		/*
			for i := range NUMBER_OF_BOTS {
				fmt.Printf("Steps: %v, Board: %v\n", bots[i].Gs.Step, bots[i].Gs.Board)
			}*/

		if generation_nr != NUMBER_GENERATIONS-1 {
			bots = children
		}
	}

	for i := range NUMBER_OF_BOTS {
		fmt.Printf("Steps: %v, Board: %v\n", bots[i].Gs.Step, bots[i].Gs.Board)
	}

}

var bots [NUMBER_OF_BOTS]Bot

type Bot struct {
	moves []game.Direction
	Gs    game.GameState
}

func (bot *Bot) generate_moves() {
	for _ = range NUMBER_OF_MOVES {
		bot.moves = append(bot.moves, game.Direction(rand.Int()%4))
	}
}

func MakeBot(game_seed int64) Bot {
	bot := Bot{}
	bot.Gs = game.MakeSeedGame(game_seed)
	bot.generate_moves()
	return bot
}

func evaluate(bot *Bot) float64 {
	linear_weight, exp_weight := 0.9, 0.1

	m := 0
	for _, row := range bot.Gs.Board {
		for _, elem := range row {
			if elem > m {
				m = elem
			}
		}
	}

	return float64(bot.Gs.Step)*linear_weight + float64(m)*exp_weight
}

func find_best_bot(bots []Bot) *Bot {
	var max_x float64
	var index int
	for i := range bots {
		x := evaluate(&bots[i])
		if x > max_x {
			max_x = x
			index = i
		}
	}
	return &bots[index]
}

// utility probability

func get_cdf(pdf [NUMBER_OF_BOTS]float64) [NUMBER_OF_BOTS]float64 {
	var out_cdf [NUMBER_OF_BOTS]float64

	out_cdf[0] = pdf[0]
	for i := 1; i < NUMBER_OF_BOTS; i++ {
		out_cdf[i] = out_cdf[i-1] + pdf[i]
	}
	return out_cdf
}

// returns chosen index
func choose(cdf []float64) int {
	r := rand.Float64()
	bucket := 0
	for r > cdf[bucket] {
		bucket++
	}
	return bucket
}

func choose2(cdf []float64) (int, int) {
	first_idx := choose(cdf)
	second_idx := choose(cdf)
	for first_idx == second_idx {
		second_idx = choose(cdf)
	}
	return first_idx, second_idx
}

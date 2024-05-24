package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"

	"github.com/bbeni/genetic_bros/game"
	"github.com/bbeni/genetic_bros/visualizer"
)

const (
	NUMBER_OF_BOTS         = 400 // has to be even too!
	NUMBER_OF_MOVES        = 1500
	NUMBER_SLICE_POSITIONS = 6 // has to be even!
	NUMBER_GENERATIONS     = 1750
	MUTATION_RATE          = 0.0002 // 10**-4 to 10**-6 https://www.sciencedirect.com/topics/biochemistry-genetics-and-molecular-biology/mutation-rate

	HAMMER_RANDOM_REPLACEMENT   = 0.03 // complete random bot
	KILL_THRESHOLD_GENENERATION = 250  // after this generation it's hammer time
	KILL_THRESHOLD              = 3000 // score below get's killed with probability:
	HAMMER_PROBABILITY          = 0.05
)

var random_gen *rand.Rand

func main() {

	// make the graph visualizer window 1
	var graph_viz1 visualizer.Graph_Viz
	graph_viz1.UserData = &visualizer.Data_Info{
		Title:  "Best Bot Evaluation Each Generation",
		XLabel: "Number of Generations",
		YLabel: "Steps Reached",
		XY:     make([]visualizer.XYData, 2),
	}
	graph_viz1.UserData.XY[0].Label = "Best Bot"
	graph_viz1.UserData.XY[1].Label = "Worst Bot"

	graph_viz1.Update_And_Draw()

	random_gen = rand.New(rand.NewSource(42))

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

		// kill some bots
		if generation_nr >= KILL_THRESHOLD_GENENERATION {

			if random_gen.Float32() > HAMMER_PROBABILITY {
				goto hammer_end
			}

			if generation_nr == KILL_THRESHOLD_GENENERATION {
				fmt.Println("It's hammer time! (by Urs B.)")
			}

			best_idx := 0
			best_score := 0.0
			for i := range NUMBER_OF_BOTS {
				if scores[i] > float64(best_score) {
					best_score = scores[i]
					best_idx = i
				}
			}

			for i := range NUMBER_OF_BOTS {
				if scores[i] < KILL_THRESHOLD {

					if random_gen.Float32() < HAMMER_RANDOM_REPLACEMENT {
						// just another pleb
						bots[i] = MakeBot(random_gen.Int63())
					} else {
						// copy the best
						alpha_bot := MakeBot(69)
						for i := range NUMBER_OF_MOVES {
							alpha_bot.moves[i] = bots[best_idx].moves[i]
						}
						//replace it
						bots[i] = alpha_bot
					}
				}
			}
		}
	hammer_end:

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
				slice_idx := random_gen.Int() % NUMBER_OF_MOVES
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
			for i := range NUMBER_OF_MOVES {
				r := random_gen.Float64()
				if r <= MUTATION_RATE {
					kid1.moves[i] = game.Direction(random_gen.Int() % 4)
				}
			}

			for i := range NUMBER_OF_MOVES {
				r := random_gen.Float64()
				if r <= MUTATION_RATE {
					kid2.moves[i] = game.Direction(random_gen.Int() % 4)
				}
			}

			children[pair*2] = kid1
			children[pair*2+1] = kid2
		}

		best_bot := find_best_bot(bots[:])
		best_bot_score := evaluate(best_bot)
		worst_bot := find_worst_bot(bots[:])
		//worst_bot_score := evaluate(best_bot)
		fmt.Printf("Generation: %v steps: %v score: %v\n", generation_nr, best_bot.Gs.Step, best_bot_score)

		//TOTO Add Data to viz 1
		graph_viz1.UserData.XY[0].XYs = append(graph_viz1.UserData.XY[0].XYs,
			visualizer.XY{
				Y: float64(best_bot.Gs.Step),
				X: float64(generation_nr),
			})
		graph_viz1.UserData.XY[1].XYs = append(graph_viz1.UserData.XY[1].XYs,
			visualizer.XY{
				Y: float64(worst_bot.Gs.Step),
				X: float64(generation_nr),
			})
		graph_viz1.Update_And_Draw()

		if generation_nr != NUMBER_GENERATIONS-1 {
			bots = children
		}
	}

	for i := range NUMBER_OF_BOTS {
		fmt.Printf("Steps: %v, Board: %v\n", bots[i].Gs.Step, bots[i].Gs.Board)
	}

	// visualize the best bot we found on the sepecific we optimized it for
	best_bot := find_best_bot(bots[:])
	gs := game.MakeSeedGame(69)
	fmt.Println("Visualizing the game! It should have the following state in the end:")
	fmt.Println(best_bot.Gs)
	visualizer.Visualize_Game(&gs, best_bot.moves, 0.01, 70)
}

var bots [NUMBER_OF_BOTS]Bot

type Bot struct {
	moves []game.Direction
	Gs    game.GameState
}

func (bot *Bot) generate_moves() {
	for _ = range NUMBER_OF_MOVES {
		bot.moves = append(bot.moves, game.Direction(random_gen.Int()%4))
	}
}

func MakeBot(game_seed int64) Bot {
	bot := Bot{}
	bot.Gs = game.MakeSeedGame(game_seed)
	bot.generate_moves()
	return bot
}

// calculate the score that is proportional to the probability of getting chosen
func evaluate(bot *Bot) float64 {

	m := 0
	for _, row := range bot.Gs.Board {
		for _, elem := range row {
			if elem > m {
				m = elem
			}
		}
	}

	exp_factor := 0.005
	return math.Exp(float64(bot.Gs.Step+m) * exp_factor)
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

func find_worst_bot(bots []Bot) *Bot {
	var min_x float64
	var index int
	for i := range bots {
		x := evaluate(&bots[i])
		if x < min_x {
			min_x = x
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
	r := random_gen.Float64()
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

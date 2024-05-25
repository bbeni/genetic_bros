package main

import (
	"fmt"
	"image/png"
	"math"
	"math/rand"
	"os"

	"github.com/bbeni/genetic_bros/game"
	"github.com/bbeni/genetic_bros/visualizer"
)

const (
	NUMBER_OF_BOTS      = 1000 // has to be even!
	CROSSOVER_PROB      = 0.13 // prob to exchange a whole layer of nn - Disabled for now!
	NUMBER_GENERATIONS  = 60
	MUTATUON_RATE       = 0.001
	NN_MAX_START_WEIGHT = 0.5
)

/*  Neural Network Architecture 1

4x4 inputs
36  hidden  -> w1, b1
16  hidden  -> w2, b2
4   outputs -> w3, b3

*/

const N_WEIGHTS = 36*4*4 + 36 + 16*36 + 16 + 4*16 + 4

type Neural_Net struct {
	w1 [36][4][4]float32
	b1 [36]float32

	w2 [16][36]float32
	b2 [16]float32

	w3 [4][16]float32
	b3 [4]float32

	// temporary store activations
	activations1 [36]float32
	activations2 [16]float32
	activations3 [4]float32
}

func exp(input float64) float32 {
	return float32(math.Exp(input))
}

func sigmoid(input float64) float32 {
	return float32(1 / (1 + math.Exp(-input)))
}

func relu(input float64) float32 {
	return max(0, float32(input))
}

var random_gen *rand.Rand

func main() {

	viz := visualizer.Graph_Viz{}

	viz.UserData = &visualizer.Data_Info{
		XY:     make([]visualizer.XYData, 2),
		XLabel: "Generation Nr.",
		YLabel: "Score",
		Title:  fmt.Sprintf("Neural Net Genetic Algorithm Evolution (%v parameters, 3 layers)", N_WEIGHTS),
	}

	viz.UserData.XY[0].Label = "Best Bot"
	viz.UserData.XY[1].Label = "Worst Bot"

	viz.Update_And_Draw()

	random_gen = rand.New(rand.NewSource(42))

	for i := range NUMBER_OF_BOTS {
		//bots[i] = MakeBot(time.Now().UnixMilli())
		//bots[i] = MakeBot(random_gen.Int63())
		bots[i] = MakeBot(69)
	}

	var kid1 Bot
	var kid2 Bot

	for generation_nr := range NUMBER_GENERATIONS {
		var scores_sum float64
		var scores [NUMBER_OF_BOTS]float64

		for i := range NUMBER_OF_BOTS {
			for {

				move := bots[i].Nn.FeedForward(&bots[i].Gs)

				for j := 1; !bots[i].Gs.MovePossible(move) && j < 4; j++ {
					//fmt.Println("skipped", bots[i].Gs.Board, move)
					move = game.Direction((int(move) + 1) % 4)
				}

				if !bots[i].Gs.MovePossible(move) {
					fmt.Println(bots[i].Gs.Board)
					fmt.Println(move)
					panic("move impossible!")
				}

				if bots[i].Gs.Update(move) {
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

			// the pair makes 2 chilren that are clones
			//kid1 = bots[mother_idx].Clone(random_gen.Int63())
			//kid2 = bots[father_idx].Clone(random_gen.Int63())
			kid1 = *bots[mother_idx].Clone(69)
			kid2 = *bots[father_idx].Clone(69)

			//Crossover(&kid1, &kid2)

			kid1.Nn.Mutate()
			kid2.Nn.Mutate()

			children[pair*2] = kid1
			children[pair*2+1] = kid2
		}

		best_bot := find_best_bot(bots[:])
		best_score := evaluate(best_bot)
		worst_bot := find_worst_bot(bots[:])
		worst_score := evaluate(worst_bot)

		fmt.Printf("Generation %v %v %v\n", generation_nr, best_bot.Gs.Step, best_score)

		viz.UserData.XY[0].XYs = append(viz.UserData.XY[0].XYs, visualizer.XY{float64(generation_nr), best_score})
		viz.UserData.XY[1].XYs = append(viz.UserData.XY[1].XYs, visualizer.XY{float64(generation_nr), worst_score})
		viz.Update_And_Draw()

		if generation_nr != NUMBER_GENERATIONS-1 {
			bots = children
		}
	}

	for i := range NUMBER_OF_BOTS {
		fmt.Printf("Steps: %v, Board: %v\n", bots[i].Gs.Step, bots[i].Gs.Board)
	}

	// save the end plot
	plot_img, _ := visualizer.Make_Plot(viz.W, viz.H, viz.UserData)

	f, err := os.Create("plot_neural_net.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err = png.Encode(f, plot_img); err != nil {
		fmt.Printf("failed to encode: %v", err)
	}

	//  take the best bot and emulate
	best := *find_best_bot(bots[:])
	best.Gs = game.MakeSeedGame(69)
	moves := make([]game.Direction, 0)

	for {
		move := best.Nn.FeedForward(&best.Gs)

		for j := 1; !best.Gs.MovePossible(move) && j < 4; j++ {
			//fmt.Println("skipped", bots[i].Gs.Board, move)
			move = game.Direction((int(move) + 1) % 4)
		}

		if !best.Gs.MovePossible(move) {
			fmt.Println(best.Gs.Board)
			fmt.Println(move)
			panic("move impossible!")
		}

		moves = append(moves, move)

		if best.Gs.Update(move) {
			break
		}

	}

	vs_game := game.MakeSeedGame(69)
	visualizer.Visualize_Game(&vs_game, moves, 0.0001, 2000)
}

var bots [NUMBER_OF_BOTS]Bot

type Bot struct {
	Nn *Neural_Net
	Gs game.GameState
}

func MakeBot(game_seed int64) Bot {
	bot := Bot{}
	bot.Gs = game.MakeSeedGame(game_seed)
	bot.Nn = MakeNeuralNet()
	return bot
}

func evaluate(bot *Bot) float64 {

	m := 0
	for _, row := range bot.Gs.Board {
		for _, elem := range row {
			if elem > m {
				m = elem
			}
		}
	}

	return math.Exp(0.01*float64(bot.Gs.Step)) + float64(m)/10
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

func MakeNeuralNet() *Neural_Net {
	nn := Neural_Net{}

	// layer 1 initialize
	for i := range 36 {
		for j := range 4 {
			for k := range 4 {
				nn.w1[i][j][k] = NN_MAX_START_WEIGHT * (2*random_gen.Float32() - 1) / 36
			}
		}
	}

	for i := range 36 {
		nn.b1[i] = NN_MAX_START_WEIGHT * (2*random_gen.Float32() - 1) / 36
	}

	// layer 2 initialize
	for i := range 16 {
		for j := range 36 {
			nn.w2[i][j] = NN_MAX_START_WEIGHT * (2*random_gen.Float32() - 1) / 16
		}
	}

	for i := range 16 {
		nn.b2[i] = NN_MAX_START_WEIGHT * (2*random_gen.Float32() - 1) / 16
	}

	// layer 3 initialize
	for i := range 4 {
		for j := range 16 {
			nn.w3[i][j] = NN_MAX_START_WEIGHT * (2*random_gen.Float32() - 1) / 4
		}
	}

	for i := range 4 {
		nn.b3[i] = NN_MAX_START_WEIGHT * (2*random_gen.Float32() - 1) / 4
	}

	return &nn
}

func (nn *Neural_Net) FeedForward(gs *game.GameState) game.Direction {

	// first layer
	for i := range 36 {
		nn.activations1[i] = nn.b1[i]
	}

	for i := range 36 {
		for j := range 4 {
			for k := range 4 {
				nn.activations1[i] += nn.w1[i][j][k] * float32(gs.Board[j][k]) // input here
			}
		}
	}

	for i := range 36 {
		nn.activations1[i] = nn.activations1[i]
	}

	// second layer
	for i := range 16 {
		nn.activations2[i] = nn.b2[i]
	}

	for i := range 16 {
		for j := range 36 {
			nn.activations2[i] += nn.w2[i][j] * nn.activations1[j]
		}
	}

	for i := range 16 {
		nn.activations2[i] = nn.activations2[i]
	}

	// third layer
	for i := range 4 {
		nn.activations3[i] = nn.b2[i]
	}

	for i := range 4 {
		for j := range 16 {
			nn.activations3[i] += nn.w3[i][j] * nn.activations2[j]
		}
	}

	for i := range 4 {
		nn.activations3[i] = nn.activations3[i]
	}

	// find argmax and cast to direction
	var m float32 = -math.MaxFloat32
	var mi int = 0

	for i := range 4 {
		if nn.activations3[i] > m {
			m = nn.activations3[i]
			mi = i
		}
	}

	dir := game.Direction(mi)
	return dir
}

func (nn *Neural_Net) Mutate() {

	for i := range 36 {
		for j := range 4 * 4 {
			r := random_gen.Float64()
			if r <= MUTATUON_RATE {
				nn.w1[i][j/4][j%4] = NN_MAX_START_WEIGHT * (2*random_gen.Float32() - 1) / 36
			}
		}
	}

	for i := range 36 {
		r := random_gen.Float64()
		if r <= MUTATUON_RATE {
			nn.b1[i] = NN_MAX_START_WEIGHT * (2*random_gen.Float32() - 1) / 36
		}
	}

	for i := range 16 * 36 {
		r := random_gen.Float64()
		if r <= MUTATUON_RATE {
			nn.w2[i/36][i%36] = NN_MAX_START_WEIGHT * (2*random_gen.Float32() - 1) / 16
		}
	}

	for i := range 16 {
		r := random_gen.Float64()
		if r <= MUTATUON_RATE {
			nn.b1[i] = NN_MAX_START_WEIGHT * (2*random_gen.Float32() - 1) / 16
		}
	}

	for i := range 4 * 16 {
		r := random_gen.Float64()
		if r <= MUTATUON_RATE {
			nn.w2[i/16][i%16] = NN_MAX_START_WEIGHT * (2*random_gen.Float32() - 1) / 4
		}
	}

	for i := range 4 {
		r := random_gen.Float64()
		if r <= MUTATUON_RATE {
			nn.b1[i] = NN_MAX_START_WEIGHT * (2*random_gen.Float32() - 1) / 4
		}
	}
}

func Crossover(bot1 *Bot, bot2 *Bot) {
	nn1 := bot1.Nn
	nn2 := bot2.Nn

	if random_gen.Float32() < CROSSOVER_PROB {
		for i := range 36 {
			for j := range 4 {
				for k := range 4 {
					nn1.w1[i][j][k], nn2.w1[i][j][k] = nn2.w1[i][j][k], nn1.w1[i][j][k]
				}
			}
		}
		for i := range 36 {
			nn1.b1[i], nn2.b1[i] = nn2.b1[i], nn1.b1[i]
		}
	}

	if random_gen.Float32() < CROSSOVER_PROB {

		for i := range 16 {
			for j := range 36 {
				nn1.w2[i][j], nn2.w2[i][j] = nn2.w2[i][j], nn1.w2[i][j]
			}
		}

		for i := range 16 {
			nn1.b2[i], nn2.b2[i] = nn2.b2[i], nn1.b2[i]
		}
	}

	if random_gen.Float32() < CROSSOVER_PROB {

		for i := range 4 {
			for j := range 16 {
				nn1.w3[i][j], nn2.w3[i][j] = nn2.w3[i][j], nn1.w3[i][j]
			}
		}

		for i := range 4 {
			nn1.b3[i], nn2.b3[i] = nn2.b3[i], nn1.b3[i]
		}
	}

}

func (bot *Bot) Clone(seed int64) *Bot {
	newb := MakeBot(seed)
	*newb.Nn = *bot.Nn
	return &newb
}

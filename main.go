package main

import (
	"github.com/bbeni/genetic_bros/game"
	"github.com/bbeni/genetic_bros/visualizer"
)

func main() {
	gs := game.MakeSeedGame(69)

	move_list := []game.Direction{
		game.East, game.West, game.East, game.South, game.North,
		game.West, game.South, game.North, game.West, game.East,
		game.West, game.East, game.South, game.North, game.West,
		game.South, game.North, game.West,
	}

	visualizer.Visualize_Game(&gs, move_list, 0.1, 1)
}

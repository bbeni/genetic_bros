package main

import (
	"time"

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

	driver := visualizer.Game_Driver{
		DriverMoves:     move_list,
		MoveTime:        0.4,
		DelayOnGameover: 1.0,
	}

	vis_game := visualizer.New_Game_Visual(&gs, &driver)

	for !vis_game.Destroyed {
		vis_game.Update_And_Draw()

		time.Sleep(time.Millisecond * 10)
		vis_game.GameTime += 0.015
	}
}

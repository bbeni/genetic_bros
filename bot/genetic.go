package main

import (
	"fmt"
	"math/rand"

	"github.com/bbeni/genetic_bros/game"
)

func main() {
	rand.Seed(69)
	gamestate := game.MakeGame()
	bot := Bot{}
	bot.generate_moves()
	fmt.Println(bot)
	for _, move := range bot.moves {
		gamestate.Update(move)
		fmt.Println(gamestate.Board, gamestate.Step)
		if gamestate.GameOver() {
			break
		}
	}
}

type Bot struct {
	moves []game.Direction
}

func (bot *Bot) generate_moves() {
	for _ = range 500 {
		bot.moves = append(bot.moves, game.Direction(rand.Int()%4))
	}
}

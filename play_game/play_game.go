package main

import (
	"github.com/bbeni/genetic_bros/game"
	"github.com/bbeni/genetic_bros/visualizer"
)

func main() {
	g := game.MakeGame()
	visualizer.Play_Game(&g)
}

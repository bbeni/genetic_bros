package main

import (
	"fmt"

	"github.com/bbeni/genetic_bros/game"
)

func main() {
	g := game.MakeGame()
	fmt.Println(g)

	fmt.Println("West")
	g.Move(game.West)
	fmt.Println(g)

	fmt.Println("East")
	g.Move(game.East)
	fmt.Println(g)

	fmt.Println("North")
	g.Move(game.North)
	fmt.Println(g)

	fmt.Println("South")
	g.Move(game.South)
	fmt.Println(g)
}

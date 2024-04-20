package main

import "fmt"

type GameState struct {
	Board [4][4]int
	Step  int
}

func MakeGame() GameState {
	return GameState{} // just empty for now
}

func main() {
	game := MakeGame()
	fmt.Println("Hello Martin!")
	fmt.Println(game)
}

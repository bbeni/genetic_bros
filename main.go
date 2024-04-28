package main

import (
	"fmt"
	"math/rand"
)

type GameState struct {
	Board [4][4]int
	Step  int
}

type Direction uint8

const (
	East Direction = iota
	South
	West
	North
)

func (game_state *GameState) Move(direction Direction) {
	if direction == West {
		for j := range 4 {
			c_i := 0
			for i := range 4 {
				if game_state.Board[j][i] != 0 {
					if game_state.Board[j][c_i] == 0 {
						game_state.Board[j][c_i] = game_state.Board[j][i]
						game_state.Board[j][i] = 0

					}else if game_state.Board[j][i] == game_state.Board[j][c_i]{
						game_state.Board[j][c_i] += game_state.Board[j][i]
						game_state.Board[j][i] = 0
					
					}else {
						game_state.Board[j][c_i+1] = game_state.Board[j][i]
						game_state.Board[j][i] = 0
					}
					c_i++
					
				}else {
					panic("Not implemented yet")
				}
				
			}
		}
		
}

func MakeGame() GameState {
	game := GameState{}
	random_index := rand.Int() % 16
	game.Board[random_index/4][random_index%4] = 2
	return game
}

func (game GameState) String() string {
	str := ""
	for i := range 4 {
		str = str + fmt.Sprintf("%4v %4v %4v %4v\n\n", game.Board[i][0], game.Board[i][1], game.Board[i][2], game.Board[i][3])
	}
	return str
}

func main() {
	game := MakeGame()
	fmt.Println(game)
	game.Move(West)
	fmt.Println(game)
	game.Move(West)
	fmt.Println(game)
}

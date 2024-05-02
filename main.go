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
	switch direction {
	case West:
		for j := range 4 {
			c_i := 0
			for i := 1; i < 4; i++ {
				if game_state.Board[j][i] != 0 {
					if game_state.Board[j][c_i] == 0 {
						game_state.Board[j][c_i] = game_state.Board[j][i]
						game_state.Board[j][i] = 0

					} else {
						if game_state.Board[j][i] == game_state.Board[j][c_i] {
							game_state.Board[j][c_i] += game_state.Board[j][i]
						} else {
							game_state.Board[j][c_i+1] = game_state.Board[j][i]
						}
						game_state.Board[j][i] = 0
						c_i++
					}
				}
			}
		}
	case East:
		for j := range 4 {
			c_i := 3
			for i := 2; i >= 0; i-- {
				if game_state.Board[j][i] != 0 {
					if game_state.Board[j][c_i] == 0 {
						game_state.Board[j][c_i] = game_state.Board[j][i]
						game_state.Board[j][i] = 0

					} else {
						if game_state.Board[j][i] == game_state.Board[j][c_i] {
							game_state.Board[j][c_i] += game_state.Board[j][i]
						} else {
							game_state.Board[j][c_i-1] = game_state.Board[j][i]
						}
						game_state.Board[j][i] = 0
						c_i--
					}
				}
			}
		}
	case North:
		for i := range 4 {
			c_j := 0
			for j := 1; j < 4; j++ {
				if game_state.Board[j][i] != 0 {
					if game_state.Board[c_j][i] == 0 {
						game_state.Board[c_j][i] = game_state.Board[j][i]
						game_state.Board[j][i] = 0

					} else {
						if game_state.Board[j][i] == game_state.Board[c_j][i] {
							game_state.Board[c_j][i] += game_state.Board[j][i]
						} else {
							game_state.Board[c_j+1][i] = game_state.Board[j][i]
						}
						game_state.Board[j][i] = 0
						c_j++
					}
				}
			}
		}
	case South:
		for i := range 4 {
			c_j := 3
			for j := 2; j >= 0; j-- {
				if game_state.Board[j][i] != 0 {
					if game_state.Board[c_j][i] == 0 {
						game_state.Board[c_j][i] = game_state.Board[j][i]
						game_state.Board[j][i] = 0

					} else {
						if game_state.Board[j][i] == game_state.Board[c_j][i] {
							game_state.Board[c_j][i] += game_state.Board[j][i]
						} else {
							game_state.Board[c_j-1][i] = game_state.Board[j][i]
						}
						game_state.Board[j][i] = 0
						c_j--
					}
				}
			}
		}
	}

	for {
		random_index := rand.Int() % 16
		if game_state.Board[random_index/4][random_index%4] == 0 {
			game_state.Board[random_index/4][random_index%4] = 2
			break
		}
	}
}

func MakeGame() GameState {
	game := GameState{}
	random_index := rand.Int() % 16
	game.Board[random_index/4][random_index%4] = 2
	random_index = rand.Int() % 16
	game.Board[random_index/4][random_index%4] = 2
	random_index = rand.Int() % 16
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

	fmt.Println("West")
	game.Move(West)
	fmt.Println(game)

	fmt.Println("East")
	game.Move(East)
	fmt.Println(game)

	fmt.Println("North")
	game.Move(North)
	fmt.Println(game)

	fmt.Println("South")
	game.Move(South)
	fmt.Println(game)
}

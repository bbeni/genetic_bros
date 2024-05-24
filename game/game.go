package game

import (
	"fmt"
	"math/rand"
	"time"
)

type GameState struct {
	Board           [4][4]int
	Step            int
	RandomGenerator *rand.Rand
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
			for i := 0; i < 4; i++ {

				// when origin is 0 -> do nothing
				// when target is 0 and origin is non-zero -> move
				// when target is non-zero and origin is equal -> combine, inc
				// when target is non-zero and origin is different -> move to neighbor, inc

				origin := game_state.Board[j][i]
				target := game_state.Board[j][c_i]

				if origin == 0 || c_i == i {
					continue
				}

				// move or combine or move to neighbour
				game_state.Board[j][i] = 0

				if target == 0 {
					game_state.Board[j][c_i] = origin
				} else {
					if origin == target {
						// combine the tiles
						game_state.Board[j][c_i] = origin + origin
						c_i++
					} else {
						c_i++
						game_state.Board[j][c_i] = origin
					}
				}
			}
		}

	case East:
		for j := range 4 {
			c_i := 3
			for i := 3; i >= 0; i-- {

				origin := game_state.Board[j][i]
				target := game_state.Board[j][c_i]

				if origin == 0 || c_i == i {
					continue
				}

				// move or combine or move to neighbour
				game_state.Board[j][i] = 0

				if target == 0 {
					game_state.Board[j][c_i] = origin
				} else {
					if origin == target {
						// combine the tiles
						game_state.Board[j][c_i] = origin + origin
						c_i--
					} else {
						c_i--
						game_state.Board[j][c_i] = origin
					}
				}
			}
		}
	case North:
		for i := range 4 {
			c_j := 0
			for j := 0; j < 4; j++ {
				origin := game_state.Board[j][i]
				target := game_state.Board[c_j][i]

				if origin == 0 || c_j == j {
					continue
				}

				// move or combine or move to neighbour
				game_state.Board[j][i] = 0

				if target == 0 {
					game_state.Board[c_j][i] = origin
				} else {
					if origin == target {
						// combine the tiles
						game_state.Board[c_j][i] = origin + origin
						c_j++
					} else {
						c_j++
						game_state.Board[c_j][i] = origin
					}
				}
			}
		}
	case South:
		for i := range 4 {
			c_j := 3
			for j := 3; j >= 0; j-- {
				origin := game_state.Board[j][i]
				target := game_state.Board[c_j][i]

				if origin == 0 || c_j == j {
					continue
				}

				// move or combine or move to neighbour
				game_state.Board[j][i] = 0

				if target == 0 {
					game_state.Board[c_j][i] = origin
				} else {
					if origin == target {
						// combine the tiles
						game_state.Board[c_j][i] = origin + origin
						c_j--
					} else {
						c_j--
						game_state.Board[c_j][i] = origin
					}
				}
			}
		}
	}
}

// MaxValue returns the highest integer in the Board
func (game_state *GameState) MaxValue() int {
	max := game_state.Board[0][0]
	for j := 0; j < 4; j++ {
		for i := 0; i < 4; i++ {
			if game_state.Board[j][i] > max {
				max = game_state.Board[j][i]
			}
		}
	}
	return max
}

// returns true if it is over
func (game_state *GameState) Update(dir Direction) bool {

	if !game_state.MovePossible(dir) {
		return false
	}

	game_state.Move(dir)

	for {
		random_index := game_state.RandomGenerator.Int() % 16
		if game_state.Board[random_index/4][random_index%4] == 0 {
			game_state.Board[random_index/4][random_index%4] = 2
			break
		}
	}

	game_state.Step++
	return game_state.GameOver()
}

func (game_state *GameState) MovePossible(dir Direction) bool {

	// copy the board over
	var test_game GameState
	for j := range 4 {
		for i := range 4 {
			test_game.Board[j][i] = game_state.Board[j][i]
		}
	}

	test_game.Move(dir)

	for j := range 4 {
		for i := range 4 {
			if test_game.Board[j][i] != game_state.Board[j][i] {
				return true
			}
		}
	}
	return false
}

func (game_state *GameState) GameOver() bool {

	for _, d := range [4]Direction{North, South, East, West} {
		if game_state.MovePossible(d) {
			return false
		}
	}

	return true
}
func MakeSeedGame(seed int64) GameState {
	game := GameState{}
	game.RandomGenerator = rand.New(rand.NewSource(seed))

	random_index := game.RandomGenerator.Int() % 16
	game.Board[random_index/4][random_index%4] = 2
	random_index = game.RandomGenerator.Int() % 16
	game.Board[random_index/4][random_index%4] = 2
	random_index = game.RandomGenerator.Int() % 16
	game.Board[random_index/4][random_index%4] = 2

	return game
}
func MakeGame() GameState {
	game := MakeSeedGame(time.Now().UnixNano())
	return game
}

func (game GameState) String() string {
	str := ""
	for i := range 4 {
		str = str + fmt.Sprintf("%4v %4v %4v %4v\n\n", game.Board[i][0], game.Board[i][1], game.Board[i][2], game.Board[i][3])
	}
	return str
}

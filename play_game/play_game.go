package main

import (
	"github.com/bbeni/genetic_bros/game"
	"github.com/bbeni/genetic_bros/visualizer"

	"runtime"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	g := game.MakeSeedGame(69)
	visualizer.Play_Game(&g)
}

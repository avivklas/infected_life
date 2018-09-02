package main

import (
	"flag"
	"bitbucket.org/avivklas/infected-life"
	"strings"
	"os"
	"fmt"
)

// parseSeed parses a string represents a grid of the game, along with width and height
// and returns Grid filled with cell values
func parseSeed (w, h int, seed string) *infected_life.Grid {
	grid := infected_life.NewGrid(w, h)
	for i, v := range strings.Split(seed, " ") {
		y := i / w
		x := i % w
		var alive bool
		if v == "1" {
			alive = true
		}
		grid.Set(x, y, alive)
	}
	return grid
}

func main () {

	width := flag.Int("width", 3, " The width of the world")
	height := flag.Int("height", 3, " The height of the world")
	infectAfter := flag.Int("infect-after", 1, " The number of generations after which the infection stage will start")
	maxGenerations := flag.Int("max-generations", 1, " The maximum number of generations that can be created. Including all phases of the game")
	seed := flag.String("seed", "", " The initial state of the world")

	flag.Parse()

	if *width < 1 || *height < 1 {
		fmt.Println("width and height must be higher that 1")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *infectAfter < 1{
		fmt.Println("infect-after must be higher that 1")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *maxGenerations < 1{
		fmt.Println("max-generations must be higher that 1")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Initiate new game with the specified parameters
	game := infected_life.NewGame(*width, *height, *infectAfter, *maxGenerations)

	if *seed == "" {
		fmt.Println("seed must be provided")
		flag.PrintDefaults()
		os.Exit(1)
	}

	initialGrid := parseSeed(*width, *height, *seed)

	// Apply the initial green
	game.Seed(*initialGrid)

	// Print the initial state of the game
	game.Grid.Print()

	// Advance game steps until game is over
	for !game.Over() {
		game.TimeStep()
		game.Grid.Print()
	}

}
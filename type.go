// An extensible implementation of Conway's "Game of Life"
// with Panoply's own "Infection" rules
package infected_life

import "fmt"

// Grid represents a two-dimensional grid of cells.
type Grid [][]bool

// NewGrid returns an empty grid of the specified width and height.
func NewGrid (w, h int) *Grid {
	g := make(Grid, h)
	for i := range g {
		g[i] = make([]bool, w)
	}
	return &g
}

// Set changes the state of a cell at the specified coordinates
// to the specified state
func (g *Grid) Set(x, y int, alive bool )  {
	(*g)[y][x] = alive
}

// Width returns the width of the grid
func (g *Grid) Width () int {
	return len((*g)[0])
}

// Height returns the height of the grid
func (g *Grid) Height () int {
	return len(*g)
}

// Alive returns whether the cell at the specified position is alive.
func (g *Grid) Alive (x, y int) bool {
	if x < 0 || y < 0 {
		return false
	}
	if x >= g.Width() || y >= g.Height() {
		return false
	}
	return (*g)[y][x]
}

// Print prints the current state of the grid in a flat format
func (g *Grid) Print () {
	for row := 0; row < g.Height(); row++ {
		for col := 0; col < g.Width(); col++ {
			if g.Alive(col, row) {
				fmt.Print("1")
			} else {
				fmt.Print("0")
			}
		}
	}
	fmt.Print("\n")
}

// DiagonalLiveNeighbours returns the amount of diagonally adjacent live cells
func (g *Grid) DiagonalLiveNeighbours (x, y int) int {
	alive := 0
	if g.Alive(x-1,y-1) {
		alive++
	}
	if g.Alive(x+1,y+1) {
		alive++
	}
	if g.Alive(x+1,y-1) {
		alive++
	}
	if g.Alive(x-1,y+1) {
		alive++
	}
	return alive
}

// HorizontalAndVerticalLiveNeighbours returns the amount of horizontally and vertically adjacent live cells
func (g *Grid) HorizontalAndVerticalLiveNeighbours (x, y int) int {
	alive := 0
	if g.Alive(x-1,y) {
		alive++
	}
	if g.Alive(x+1,y) {
		alive++
	}
	if g.Alive(x,y-1) {
		alive++
	}
	if g.Alive(x,y+1) {
		alive++
	}
	return alive
}

// LiveNeighbours returns the amount of all adjacent live cells
func (g *Grid) LiveNeighbours (x, y int) int {
	return g.DiagonalLiveNeighbours(x,y) + g.HorizontalAndVerticalLiveNeighbours(x,y)
}

// Rules is the interface that wraps a set of rules that should be applied
// to each cell through the Apply(x,y)
type Rules interface {
	// Apply returns the next state of cell at the specified coordinates
	Apply(x, y int) bool
}

// Life implements Rules by the rules of Conway's "Game of Life"
type Life struct {
	g	*Grid
}

// Apply implements Rules' Apply function
func (l *Life) Apply (x, y int) bool {
	liveNeighbours := l.g.LiveNeighbours(x,y)
	// According to Conway's "Game of Life", simplified:
	// If exactly 3 neighbours: alive (true)
	// If exactly 2 neighbours: state does not change
	// Otherwise: cell is dead (false)
	return liveNeighbours == 3 || liveNeighbours == 2 && l.g.Alive(x, y)
}

// Life implements Rules by the rules of Panoply's "Game of Life: Infection"
type Infection struct {
	g	*Grid
}

// Apply implements Rules' Apply function
func (i *Infection) Apply (x, y int) bool {
	alive := i.g.Alive(x,y)
	horizontalAndVerticalNeighbours := i.g.HorizontalAndVerticalLiveNeighbours(x,y)
	diagonalNeighbours := i.g.DiagonalLiveNeighbours(x,y)
	// According to Panoply's "Game of Life: Infection":
	// If cell is dead and has 1 neighbour: alive (true)
	// If cell is alive and has 0 horizontal and 0 vertical neighbours: dead (false)
	if !alive && horizontalAndVerticalNeighbours+diagonalNeighbours == 1 {
		return true
	}
	if alive && horizontalAndVerticalNeighbours > 0 {
		return true
	}
	return false
}

// Game wraps the basic definitions of the game of life
// and stores the current state of the grid
type Game struct {
	Grid              *Grid
	InfectAfter       int
	MaxGenerations    int

	currentGeneration int
}

// NewGame returns an initial Game with the specified definitions
func NewGame (width, height, infectAfter, maxGenerations int) *Game {
	return &Game{
		Grid:			NewGrid(width, height),
		InfectAfter:	infectAfter,
		MaxGenerations:	maxGenerations,
	}
}

// Seed replaces the game's grid with the specified one
func (g *Game) Seed(grid Grid) {
	g.Grid = &grid
}

// Over returns true if the game is over
// false otherwise
func (g *Game) Over() bool {
	return g.currentGeneration >= g.MaxGenerations
}

// currentRules returns the rules that should be applied to cells
// according to the current game state
func (g *Game) currentRules() Rules {
	if g.currentGeneration <= g.InfectAfter {
		return &Life{g.Grid}
	} else {
		return &Infection{g.Grid}
	}
}

// TimeStep advances the game in 1 generation
// and applies the new state into the grid
func (g *Game) TimeStep () {
	// Initiates a new grid to store the new state of each cell
	// without affecting calculations of other cells
	nextGeneration := NewGrid(g.Grid.Width(), g.Grid.Height())
	rules := g.currentRules()
	for y := 0; y < g.Grid.Height(); y++ {
		for x := 0; x < g.Grid.Width(); x++ {
			nextGeneration.Set(x, y, rules.Apply(x, y))
		}
	}
	// Replace the current state of the game to the new generation
	g.Grid = nextGeneration
	// Updates the current generation of the game
	g.currentGeneration++
}
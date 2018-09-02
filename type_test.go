package infected_life

import (
	"testing"
	"strconv"
)

func TestGrid_Set_Alive(t *testing.T) {
	g := NewGrid(1,1)
	if g.Alive(0,0) {
		t.Fatal("unexpected initial state of Grid")
	}
	g.Set(0,0,true)
	if !g.Alive(0,0) {
		t.Error("expected cell state to equal: alive. actual: dead")
		t.Fail()
	}
	g.Set(0,0,false)
	if g.Alive(0,0) {
		t.Error("expected cell state to equal: dead. actual: alive")
		t.Fail()
	}
}

func TestGrid_Width(t *testing.T) {
	g := NewGrid(2,4)
	if g.Width() != 2 {
		t.Errorf("expected width to equal: 2. actual: %d", g.Width())
		t.Fail()
	}
}

func TestGrid_Height(t *testing.T) {
	g := NewGrid(2,4)
	if g.Height() != 4 {
		t.Errorf("expected height to equal: 4. actual: %d", g.Height())
		t.Fail()
	}
}

func TestGrid_LiveNeighbours(t *testing.T) {

	cases := []struct{
		name		string
		w,h			int
		liveCells	[][2]int
		testCell	[2]int
		expected 	int
	} {
		{
			name: 		"1 cell grid - alive",
			w:			1,
			h:			1,
			liveCells:	[][2]int{{0,0}},
			testCell: 	[2]int{0,0},
			expected:	0,
		},
		{
			name: 		"1 cell grid - dead",
			w:			1,
			h:			1,
			liveCells:	[][2]int{},
			testCell: 	[2]int{0,0},
			expected:	0,
		},
		{
			name: 		"9 cell square grid - middle - horizontal",
			w:			3,
			h:			3,
			liveCells:	[][2]int{{0,1},{2,1}},
			testCell: 	[2]int{1,1},
			expected:	2,
		},
		{
			name: 		"9 cell square grid - middle - vertical",
			w:			3,
			h:			3,
			liveCells:	[][2]int{{1,0},{1,2}},
			testCell: 	[2]int{1,1},
			expected:	2,
		},
		{
			name: 		"9 cell square grid - middle-top - horizontal",
			w:			3,
			h:			3,
			liveCells:	[][2]int{{0,1},{2,1}},
			testCell: 	[2]int{1,0},
			expected:	2,
		},
		{
			name: 		"9 cell square grid - middle-top - vertical",
			w:			3,
			h:			3,
			liveCells:	[][2]int{{1,0},{1,2}},
			testCell: 	[2]int{1,0},
			expected:	0,
		},
		{
			name: 		"9 cell square grid - middle-left - horizontal",
			w:			3,
			h:			3,
			liveCells:	[][2]int{{0,1},{2,1}},
			testCell: 	[2]int{0,0},
			expected:	1,
		},
		{
			name: 		"9 cell square grid - middle-left - vertical",
			w:			3,
			h:			3,
			liveCells:	[][2]int{{1,0},{1,2}},
			testCell: 	[2]int{0,0},
			expected:	1,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			g := NewGrid(c.w, c.h)
			for _, cell := range c.liveCells {
				g.Set(cell[0], cell[1], true)
			}
			liveNeighbours := g.LiveNeighbours(c.testCell[0],c.testCell[1])
			if liveNeighbours != c.expected {
				t.Errorf("expected neighbour count to equal: %d. actual: %d", c.expected, liveNeighbours)
				t.Fail()
			}
		})
	}

}

func TestLife_Apply(t *testing.T) {

	cases := []struct{
		name		string
		liveCells	[][2]int
		testCell	[2]int
		expected 	bool
	} {
		{
			name: 		"rule #1",
			liveCells:	[][2]int{{1,1},{0,0}},
			testCell: 	[2]int{1,1},
			expected:	false,
		},
		{
			name: 		"rule #2 - two",
			liveCells:	[][2]int{{1,1},{0,0},{1,0}},
			testCell: 	[2]int{1,1},
			expected:	true,
		},
		{
			name: 		"rule #2 - three",
			liveCells:	[][2]int{{1,1},{0,0},{1,0},{2,0}},
			testCell: 	[2]int{1,1},
			expected:	true,
		},
		{
			name: 		"rule #3",
			liveCells:	[][2]int{{1,1},{0,0},{1,0},{2,0},{0,1}},
			testCell: 	[2]int{1,1},
			expected:	false,
		},
		{
			name: 		"rule #4",
			liveCells:	[][2]int{{0,0},{1,0},{2,0}},
			testCell: 	[2]int{1,1},
			expected:	true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			g := NewGrid(3, 3)
			for _, cell := range c.liveCells {
				g.Set(cell[0], cell[1], true)
			}
			life := &Life{g}
			alive := life.Apply(c.testCell[0],c.testCell[1])
			if alive != c.expected {
				t.Errorf("test cell state to equal: %s. actual: %s", strconv.FormatBool(c.expected), strconv.FormatBool(alive))
				t.Fail()
			}
		})
	}

}

func TestInfection_Apply(t *testing.T) {

	cases := []struct{
		name		string
		liveCells	[][2]int
		testCell	[2]int
		expected 	bool
	} {
		{
			name: 		"rule #1",
			liveCells:	[][2]int{{0,0}},
			testCell: 	[2]int{1,1},
			expected:	true,
		},
		{
			name: 		"rule #2",
			liveCells:	[][2]int{{1,1}},
			testCell: 	[2]int{1,1},
			expected:	false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			g := NewGrid(3, 3)
			for _, cell := range c.liveCells {
				g.Set(cell[0], cell[1], true)
			}
			infection := &Infection{g}
			alive := infection.Apply(c.testCell[0],c.testCell[1])
			if alive != c.expected {
				t.Errorf("test cell state to equal: %s. actual: %s", strconv.FormatBool(c.expected), strconv.FormatBool(alive))
				t.Fail()
			}
		})
	}

}

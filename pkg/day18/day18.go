package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"

	"golang.org/x/exp/constraints"
)

// Lazy copy-paste functions
func toInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func min[T constraints.Ordered](args ...T) T {
	m := args[0]
	for _, arg := range args {
		if arg < m {
			m = arg
		}
	}
	return m
}

func max[T constraints.Ordered](args ...T) T {
	m := args[0]
	for _, arg := range args {
		if arg > m {
			m = arg
		}
	}
	return m
}

type Grid struct {
	lights [][]bool
}

func LightsArray(n int) [][]bool {
	lights := make([][]bool, n)
	for i := 0; i < n; i++ {
		lights[i] = make([]bool, n)
	}
	return lights
}

func NewGrid() *Grid {
	return &Grid{LightsArray(100)}
}

func (g *Grid) SetState(row int, state string) {
	for i, c := range state {
		g.lights[row][i] = c == '#'
	}
}

func (g *Grid) At(row, col int) bool {
	if row < 0 || row >= 100 || col < 0 || col >= 100 {
		return false
	}
	return g.lights[row][col]
}

func (g *Grid) Neighbors(row, col int) int {
	dirs := []int{-1, 0, 1}
	count := 0
	for _, dr := range dirs {
		for _, dc := range dirs {
			if dr == 0 && dc == 0 {
				continue
			}
			if g.At(row+dr, col+dc) {
				count++
			}
		}
	}
	return count
}

func (g *Grid) Tick(part2 ...bool) {
	nextState := LightsArray(100)
	for i, row := range g.lights {
		for j, light := range row {
			count := g.Neighbors(i, j)
			if light && (count == 2 || count == 3) {
				nextState[i][j] = true
			} else if !light && count == 3 {
				nextState[i][j] = true
			}
		}
	}
	if len(part2) > 0 && part2[0] {
		nextState[0][0] = true
		nextState[0][99] = true
		nextState[99][0] = true
		nextState[99][99] = true
	}
	g.lights = nextState
}

func (g *Grid) String() string {
	s := ""
	for _, row := range g.lights {
		for _, light := range row {
			if light {
				s += "#"
			} else {
				s += "."
			}
		}
		s += "\n"
	}
	return s
}

func (g *Grid) CountOn() int {
	count := 0
	for _, row := range g.lights {
		for _, light := range row {
			if light {
				count++
			}
		}
	}
	return count
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	g := NewGrid()
	g2 := NewGrid()
	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		g.SetState(i, line)
		g2.SetState(i, line)
		i++
	}
	fmt.Println(g)

	for i := 0; i < 100; i++ {
		g.Tick()
		// Clear screen
		fmt.Print("\033[H\033[2J")
		fmt.Println(g)
		time.Sleep(50 * time.Millisecond)
	}

	fmt.Printf("Part 1: %d\n", g.CountOn())

	// Part 2
	g2.lights[0][0] = true
	g2.lights[0][99] = true
	g2.lights[99][0] = true
	g2.lights[99][99] = true
	for i := 0; i < 100; i++ {
		g2.Tick(true)
		// Clear screen
		fmt.Print("\033[H\033[2J")
		fmt.Println(g2)
		time.Sleep(50 * time.Millisecond)
	}
	fmt.Printf("Part 2: %d\n", g2.CountOn())
}

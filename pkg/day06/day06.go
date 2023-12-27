package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type Direction string

const (
	TurnOn  Direction = "turn on"
	TurnOff Direction = "turn off"
	Toggle  Direction = "toggle"
)

var instructionRegex = regexp.MustCompile(`^(turn on|turn off|toggle) (\d+),(\d+) through (\d+),(\d+)$`)

type Point struct {
	x int
	y int
}

func min(args ...int) int {
	min := args[0]
	for _, v := range args {
		if v < min {
			min = v
		}
	}
	return min
}

func max(args ...int) int {
	max := args[0]
	for _, v := range args {
		if v > max {
			max = v
		}
	}
	return max
}

type Instruction struct {
	dir   Direction
	start Point
	end   Point
}

func (i *Instruction) IterPoints() <-chan Point {
	ch := make(chan Point)
	go func() {
		xlo := min(i.start.x, i.end.x)
		xhi := max(i.start.x, i.end.x)
		ylo := min(i.start.y, i.end.y)
		yhi := max(i.start.y, i.end.y)
		for x := xlo; x <= xhi; x++ {
			for y := ylo; y <= yhi; y++ {
				ch <- Point{x, y}
			}
		}
		close(ch)
	}()
	return ch
}

func NewInstruction(s string) Instruction {
	matches := instructionRegex.FindStringSubmatch(s)
	if matches == nil {
		panic("invalid instruction")
	}
	sx, err := strconv.Atoi(matches[2])
	if err != nil {
		panic("invalid instruction")
	}
	sy, err := strconv.Atoi(matches[3])
	if err != nil {
		panic("invalid instruction")
	}
	ex, err := strconv.Atoi(matches[4])
	if err != nil {
		panic("invalid instruction")
	}
	ey, err := strconv.Atoi(matches[5])
	if err != nil {
		panic("invalid instruction")
	}

	return Instruction{
		dir:   Direction(matches[1]),
		start: Point{sx, sy},
		end:   Point{ex, ey},
	}
}

type Grid struct {
	grid  [][]bool
	grid2 [][]int
}

func NewGrid(n int) *Grid {
	g := make([][]bool, n)
	g2 := make([][]int, n)
	for i := range g {
		g[i] = make([]bool, n)
		g2[i] = make([]int, n)
	}
	return &Grid{g, g2}
}

func (g *Grid) Process(i Instruction) {
	switch i.dir {
	case TurnOn:
		for p := range i.IterPoints() {
			g.grid[p.x][p.y] = true
			g.grid2[p.x][p.y] += 1
		}
	case TurnOff:
		for p := range i.IterPoints() {
			g.grid[p.x][p.y] = false
			g.grid2[p.x][p.y] = max(0, g.grid2[p.x][p.y]-1)
		}
	case Toggle:
		for p := range i.IterPoints() {
			g.grid[p.x][p.y] = !g.grid[p.x][p.y]
			g.grid2[p.x][p.y] += 2
		}
	}
}

func (g *Grid) GetCounts() (int, int) {
	count := 0
	count2 := 0
	for i := range g.grid {
		for j := range g.grid[i] {
			if g.grid[i][j] {
				count += 1
			}
			count2 += g.grid2[i][j]
		}
	}
	return count, count2
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	g := NewGrid(1000)
	for scanner.Scan() {
		line := scanner.Text()
		g.Process(NewInstruction(line))
	}

	part1, part2 := g.GetCounts()
	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", part2)
}

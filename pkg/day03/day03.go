package main

import (
	"bufio"
	"fmt"
	"os"
)

type point struct {
	row int
	col int
}

func (p *point) move(dir rune) point {
	switch dir {
	case '^':
		p.row -= 1
	case 'v':
		p.row += 1
	case '<':
		p.col -= 1
	case '>':
		p.col += 1
	}
	return *p
}

func countHouses(m map[point]int) int {
	res := 0
	for _, v := range m {
		if v > 0 {
			res++
		}
	}
	return res
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
	}

	m := make(map[point]int)
	p := point{0, 0}
	m[p]++
	for _, c := range line {
		m[p.move(c)]++
	}

	fmt.Printf("Part 1: %d\n", countHouses(m))

	m = make(map[point]int)
	santa := point{0, 0}
	robot := point{0, 0}
	m[santa]++
	for i, c := range line {
		if i%2 == 0 {
			m[santa.move(c)]++
		} else {
			m[robot.move(c)]++
		}
	}

	fmt.Printf("Part 2: %d\n", countHouses(m))
}

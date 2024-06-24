package main

import "fmt"

func nextIndex(r, c int) (int, int) {
	r -= 2
	c += 1
	if r <= 0 {
		r = c
		c = 1
	}
	return r, c
}

func nextCode(code int) int {
	return code * 252533 % 33554393
}

func main() {
	var targetRow, targetCol int
	fmt.Scanf("%d %d", &targetRow, &targetCol)

	row := 1
	col := 1
	code := 20151125
	for row != targetRow || col != targetCol {
		row, col = nextIndex(row, col)
		code = nextCode(code)
	}

	fmt.Printf("Part 1: %d\n", code)
}

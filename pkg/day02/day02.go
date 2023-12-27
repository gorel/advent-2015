package day02

import (
	"bufio"
	"fmt"
	"os"
)

func min(args ...int) int {
	m := args[0]
	for _, v := range args {
		if v < m {
			m = v
		}
	}
	return m
}

func surfaceArea(l, w, h int) int {
	side1 := l * w
	side2 := w * h
	side3 := h * l
	return 2*side1 + 2*side2 + 2*side3 + min(side1, side2, side3)
}

func volume(l, w, h int) int {
	return l * w * h
}

func smallestPerimeter(l, w, h int) int {
	return min(2*l+2*w, 2*w+2*h, 2*h+2*l)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	total := 0
	ribbon := 0
	for scanner.Scan() {
		var l, w, h int
		fmt.Sscanf(scanner.Text(), "%dx%dx%d", &l, &w, &h)
		total += surfaceArea(l, w, h)
		ribbon += volume(l, w, h) + smallestPerimeter(l, w, h)
	}

	fmt.Printf("Part 1: %d\n", total)
	fmt.Printf("Part 2: %d\n", ribbon)
}

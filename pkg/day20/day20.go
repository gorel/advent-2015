package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

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

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var n int
	for scanner.Scan() {
		n = toInt(scanner.Text())
	}

	houses := make([]int, 1000000)
	for elf := 1; elf < len(houses); elf++ {
		for house := elf; house < len(houses); house += elf {
			houses[house] += elf * 10
		}
	}

	for house, presents := range houses {
		if presents >= n {
			fmt.Printf("Part 1: %d\n", house)
			break
		}
	}

	// Part 2
	houses = make([]int, 1000000)
	for elf := 1; elf < len(houses); elf++ {
		for house := elf; house < len(houses) && house <= elf*50; house += elf {
			houses[house] += elf * 11
		}
	}

	for house, presents := range houses {
		if presents >= n {
			fmt.Printf("Part 2: %d\n", house)
			break
		}
	}
}

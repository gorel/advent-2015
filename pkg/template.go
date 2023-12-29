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

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}

	fmt.Printf("Part 1: %d\n", 0)
	fmt.Printf("Part 2: %d\n", 0)
}

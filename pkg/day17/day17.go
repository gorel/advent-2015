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

func ways(containers []int, idx int, target int) int {
	if target == 0 {
		return 1
	} else if idx >= len(containers) || target < 0 {
		return 0
	} else {
		return ways(containers, idx+1, target) + ways(containers, idx+1, target-containers[idx])
	}
}

func ways2(containers []int, idx int, target int) (int, int, error) {
	if target == 0 {
		return 1, 0, nil
	} else if idx >= len(containers) || target < 0 {
		return 0, 0, fmt.Errorf("no solution")
	} else {
		// Find the better way by minimizing the number of containers used
		ways1, used1, err1 := ways2(containers, idx+1, target)
		ways2, used2, err2 := ways2(containers, idx+1, target-containers[idx])
		if err1 != nil {
			return ways2, used2 + 1, err2
		} else if err2 != nil {
			return ways1, used1, err1
		}

		used2++
		if used1 < used2 {
			// way 1 is better
			return ways1, used1, nil
		} else if used2 < used1 {
			// way2 is better
			return ways2, used2, nil
		} else {
			// they are equal
			return ways1 + ways2, used1, nil
		}
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var containers []int
	for scanner.Scan() {
		containers = append(containers, toInt(scanner.Text()))
	}

	fmt.Printf("Part 1: %d\n", ways(containers, 0, 150))
	ways, count, err := ways2(containers, 0, 150)
	if err != nil {
		fmt.Printf("Part 2: error: %s\n", err.Error())
	} else {
		fmt.Printf("Part 2: %d, %d\n", ways, count)
	}
}

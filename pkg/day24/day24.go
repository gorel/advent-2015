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

func product[T constraints.Integer](vals []T) T {
	p := T(1)
	for _, val := range vals {
		p *= val
	}
	return p
}

func combinations(nums []int, target int) [][]int {
	var result [][]int
	var currentCombination []int

	var helper func(int, int)
	helper = func(start int, target int) {
		if target == 0 {
			combination := make([]int, len(currentCombination))
			copy(combination, currentCombination)
			result = append(result, combination)
			return
		}
		for i := start; i < len(nums); i++ {
			if nums[i] <= target {
				currentCombination = append(currentCombination, nums[i])
				helper(i+1, target-nums[i])
				currentCombination = currentCombination[:len(currentCombination)-1]
			}
		}
	}

	helper(0, target)
	return result
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	sum := 0
	var packages []int
	for scanner.Scan() {
		pkg := toInt(scanner.Text())
		packages = append(packages, pkg)
		sum += pkg
	}

	cs := combinations(packages, sum/3)
	bestLen := 1 << 31
	bestQE := 0
	for _, c := range cs {
		qe := product(c)
		if len(c) < bestLen || (len(c) == bestLen && qe < bestQE) {
			bestLen = len(c)
			bestQE = qe
		}
	}

	fmt.Printf("Part 1: %d\n", bestQE)

	cs2 := combinations(packages, sum/4)
	bestLen2 := 1 << 31
	bestQE2 := 0
	for _, c := range cs2 {
		qe := product(c)
		if len(c) < bestLen2 || (len(c) == bestLen2 && qe < bestQE2) {
			bestLen2 = len(c)
			bestQE2 = qe
		}
	}
	fmt.Printf("Part 2: %d\n", bestQE2)
}

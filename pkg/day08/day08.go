package main

import (
	"bufio"
	"fmt"
	"os"
)

func charCount(s string) int {
	return len(s)
}

func stringCount(s string) int {
	s = s[1 : len(s)-1]
	res := 0
	i := 0
	for i < len(s) {
		res++
		if s[i] == '\\' {
			i++
			if s[i] == 'x' {
				i += 2
			}
		}
		i++
	}
	return res
}

func escapeCount(s string) int {
	res := 2
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' || s[i] == '"' {
			res++
		}
	}
	return res
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	part1 := 0
	part2 := 0
	for scanner.Scan() {
		line := scanner.Text()
		part1 += charCount(line) - stringCount(line)
		part2 += escapeCount(line)
	}

	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", part2)
}

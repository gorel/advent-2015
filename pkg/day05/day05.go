package main

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/exp/slices"
)

func isNice(s string) bool {
	vowels := 0
	repeated := false
	badStrings := []string{"ab", "cd", "pq", "xy"}
	for i, c := range s {
		switch c {
		case 'a', 'e', 'i', 'o', 'u':
			vowels++
		}
		if i > 0 && c == rune(s[i-1]) {
			repeated = true
		}
		if i > 0 && slices.Contains(badStrings, string(s[i-1])+string(c)) {
			return false
		}
	}
	return vowels >= 3 && repeated
}

func isNicePart2(s string) bool {
	pairs := make(map[string]int)
	repeated := false
	doublePair := false

	lastLastChar := ""
	lastChar := ""
	for i, c := range s {
		c := string(c)
		thisPair := lastChar + c

		if c == lastLastChar {
			repeated = true
		}

		if _, ok := pairs[thisPair]; ok && pairs[thisPair] != i-1 {
			doublePair = true
		}

		pairs[thisPair] = i
		lastLastChar = lastChar
		lastChar = c
	}

	return repeated && doublePair
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	nice := 0
	nice2 := 0
	for scanner.Scan() {
		line := scanner.Text()
		if isNice(line) {
			nice++
		}
		if isNicePart2(line) {
			nice2++
		}
	}
	fmt.Printf("Part 1: %d\n", nice)
	fmt.Printf("Part 2: %d\n", nice2)
}

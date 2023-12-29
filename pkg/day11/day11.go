package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func increment(s string) string {
	chars := make([]byte, len(s))
	carry := byte(1)
	for i := len(s) - 1; i >= 0; i-- {
		c := s[i] + carry
		carry = 0
		if c > 'z' {
			carry = 1
			c = 'a'
		}

		// Illegal character - skip forward
		if c == 'i' || c == 'o' || c == 'l' {
			c++
			for j := i + 1; j < len(chars); j++ {
				chars[j] = 'a'
			}
		}

		chars[i] = c
	}

	sb := strings.Builder{}
	if carry == 1 {
		sb.WriteByte('a')
	}

	for _, c := range chars {
		sb.WriteByte(c)
	}
	return sb.String()
}

func passes(s string) bool {
	hasStraight := false
	pairs := make(map[byte]int)

	for i := 0; i < len(s)-2; i++ {
		c0, c1, c2 := s[i], s[i+1], s[i+2]

		// Run of three consecutive letters
		if c0 == c1-1 && c1 == c2-1 {
			hasStraight = true
		}

		// Need at least two different pairs
		if c0 == c1 {
			pairs[c0]++
		} else if i == len(s)-3 && c1 == c2 {
			pairs[c1]++
		}
	}

	return hasStraight && len(pairs) >= 2
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		password := scanner.Text()
		for !passes(password) {
			password = increment(password)
		}
		fmt.Printf("Part 1: %s\n", password)

		password = increment(password)
		for !passes(password) {
			password = increment(password)
		}
		fmt.Printf("Part 2: %s\n", password)
	}
}

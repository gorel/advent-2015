package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	floor := 0
	part2 := false
	for scanner.Scan() {
		line := scanner.Text()
		for i, c := range line {
			if c == '(' {
				floor++
			} else {
				floor--
				if floor < 0 && !part2 {
					part2 = true
					fmt.Printf("Part 2: %d\n", i+1)
				}
			}
		}
	}

	fmt.Printf("Part 1: %d\n", floor)
}

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

func shortest(molecules map[string]int) string {
	res := ""
	for molecule := range molecules {
		if res == "" || len(molecule) < len(res) {
			res = molecule
		}
	}
	return res
}

type Reaction struct {
	from string
	to   string
}

func minReplacementSteps(start, target string, replacements []Reaction) int {
	// We start with the target moluecule and try to find the shortest path to the start
	molecules := make(map[string]int)
	molecules[target] = 0
	for {
		cur := shortest(molecules)
		for _, replacement := range replacements {
			for i := 0; i < len(cur)-len(replacement.to)+1; i++ {
				if cur[i:i+len(replacement.to)] == replacement.to {
					newMolecule := cur[:i] + replacement.from + cur[i+len(replacement.to):]
					molecules[newMolecule] = molecules[cur] + 1
					if newMolecule == start {
						return molecules[newMolecule]
					}
				}
			}
		}
		delete(molecules, cur)
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var replacements []Reaction
	var target string
	targetNext := false
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			targetNext = true
			continue
		}
		if targetNext {
			target = line
			continue
		}
		var from, to string
		fmt.Sscanf(line, "%s => %s", &from, &to)
		replacements = append(replacements, Reaction{from, to})
	}

	possibilities := make(map[string]Reaction)
	for _, reaction := range replacements {
		for i := 0; i < len(target)-len(reaction.from)+1; i++ {
			if target[i:i+len(reaction.from)] == reaction.from {
				possibilities[target[:i]+reaction.to+target[i+len(reaction.from):]] = reaction
			}
		}
	}

	fmt.Printf("Part 1: %d\n", len(possibilities))

	// Part 2
	fmt.Printf("Part 2: %d\n", minReplacementSteps("e", target, replacements))
}

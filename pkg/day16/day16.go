package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

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

type Sue struct {
	name        string
	children    int
	cats        int
	samoyeds    int
	pomeranians int
	akitas      int
	vizslas     int
	goldfish    int
	trees       int
	cars        int
	perfumes    int
}

func ParseSue(line string) Sue {
	name, rest, _ := strings.Cut(line, ": ")
	sue := Sue{
		name:        name,
		children:    -1,
		cats:        -1,
		samoyeds:    -1,
		pomeranians: -1,
		akitas:      -1,
		vizslas:     -1,
		goldfish:    -1,
		trees:       -1,
		cars:        -1,
		perfumes:    -1,
	}

	parts := strings.Split(rest, ", ")
	for _, part := range parts {
		var typeOf string
		var count int
		fmt.Sscanf(part, "%s %d", &typeOf, &count)
		switch typeOf {
		case "children:":
			sue.children = count
		case "cats:":
			sue.cats = count
		case "samoyeds:":
			sue.samoyeds = count
		case "pomeranians:":
			sue.pomeranians = count
		case "akitas:":
			sue.akitas = count
		case "vizslas:":
			sue.vizslas = count
		case "goldfish:":
			sue.goldfish = count
		case "trees:":
			sue.trees = count
		case "cars:":
			sue.cars = count
		case "perfumes:":
			sue.perfumes = count
		default:
			panic(fmt.Sprintf("Unknown type: %s", typeOf))
		}
	}

	return sue
}

func (s *Sue) MatchesTarget(target Sue, part2 bool) bool {
	return (s.children == -1 || s.children == target.children) &&
		((part2 && (s.cats == -1 || s.cats > target.cats)) || (!part2 && (s.cats == -1 || s.cats == target.cats))) &&
		(s.samoyeds == -1 || s.samoyeds == target.samoyeds) &&
		((part2 && (s.pomeranians == -1 || s.pomeranians < target.pomeranians)) || (!part2 && (s.pomeranians == -1 || s.pomeranians == target.pomeranians))) &&
		(s.akitas == -1 || s.akitas == target.akitas) &&
		(s.vizslas == -1 || s.vizslas == target.vizslas) &&
		((part2 && (s.goldfish == -1 || s.goldfish < target.goldfish)) || (!part2 && (s.goldfish == -1 || s.goldfish == target.goldfish))) &&
		((part2 && (s.trees == -1 || s.trees > target.trees)) || (!part2 && (s.trees == -1 || s.trees == target.trees))) &&
		(s.cars == -1 || s.cars == target.cars) &&
		(s.perfumes == -1 || s.perfumes == target.perfumes)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	target := Sue{
		children:    3,
		cats:        7,
		samoyeds:    2,
		pomeranians: 3,
		akitas:      0,
		vizslas:     0,
		goldfish:    5,
		trees:       3,
		cars:        2,
		perfumes:    1,
	}
	for scanner.Scan() {
		line := scanner.Text()
		sue := ParseSue(line)
		if sue.MatchesTarget(target, false) {
			fmt.Printf("Part 1: %s\n", sue.name)
		}
		if sue.MatchesTarget(target, true) {
			fmt.Printf("Part 2: %s\n", sue.name)
		}
	}
}

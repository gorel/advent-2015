package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var happyRegex = regexp.MustCompile(`^(\w+) would (gain|lose) (\d+) happiness units by sitting next to (\w+)\.$`)

func permutations(table []string) [][]string {
	if len(table) == 1 {
		return [][]string{table}
	}
	perms := [][]string{}
	for i, v := range table {
		rest := make([]string, len(table)-1)
		copy(rest, table[:i])
		copy(rest[i:], table[i+1:])
		for _, perm := range permutations(rest) {
			perms = append(perms, append([]string{v}, perm...))
		}
	}
	return perms
}

func happiness(table []string, rules map[string]map[string]int) int {
	res := 0
	for i, p := range table {
		left := table[(i-1+len(table))%len(table)]
		right := table[(i+1)%len(table)]
		res += rules[p][left] + rules[p][right]
	}
	return res
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	rules := make(map[string]map[string]int)
	attendees := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		s := line
		matches := happyRegex.FindStringSubmatch(s)
		src := matches[1]
		dst := matches[4]
		delta, _ := strconv.Atoi(matches[3])
		if matches[2] == "lose" {
			delta = -delta
		}

		if _, ok := rules[src]; !ok {
			attendees = append(attendees, src)
			rules[src] = make(map[string]int)
		}
		if _, ok := rules[dst]; !ok {
			attendees = append(attendees, dst)
			rules[dst] = make(map[string]int)
		}
		rules[src][dst] = delta
	}

	// Optimization: since it's a circular table, it doesn't matter where we seat the first person.
	// We just need the permutations of the remaining people.
	bestScore := 0
	perms := permutations(attendees[1:])
	for _, perm := range perms {
		perm = append(perm, attendees[0])
		if score := happiness(perm, rules); score > bestScore {
			bestScore = score
		}
	}
	fmt.Printf("Part 1: %d\n", bestScore)

	// Part 2: add myself to the table
	attendees = append(attendees, "me")
	rules["me"] = make(map[string]int)
	for _, p := range attendees {
		rules[p]["me"] = 0
		rules["me"][p] = 0
	}
	bestScore = 0
	perms = permutations(attendees[1:])
	for _, perm := range perms {
		perm = append(perm, attendees[0])
		if score := happiness(perm, rules); score > bestScore {
			bestScore = score
		}
	}
	fmt.Printf("Part 2: %d\n", bestScore)
}

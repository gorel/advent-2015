package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/golang-collections/collections/set"
)

var locRegex = regexp.MustCompile(`^(\w+) to (\w+) = (\d+)$`)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type Location struct {
	neighbors map[string]int
	name      string
}

func (l *Location) GetMinDistance(
	locations map[string]*Location,
	visited *set.Set,
) int {
	return l.getDistance(locations, visited, min, int(^uint(0)>>1))
}

func (l *Location) GetMaxDistance(
	locations map[string]*Location,
	visited *set.Set,
) int {
	return l.getDistance(locations, visited, max, 0)
}

func (l *Location) getDistance(
	locations map[string]*Location,
	visited *set.Set,
	f func(int, int) int,
	initialValue int,
) int {
	if visited.Len()+1 == len(locations) {
		return 0
	}

	visited.Insert(l.name)
	bestDistance := initialValue
	for neighbor, distance := range l.neighbors {
		if visited.Has(neighbor) {
			continue
		}
		distance += locations[neighbor].getDistance(locations, visited, f, initialValue)
		bestDistance = f(bestDistance, distance)
	}

	visited.Remove(l.name)
	return bestDistance
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	locations := make(map[string]*Location)
	for scanner.Scan() {
		line := scanner.Text()
		matches := locRegex.FindStringSubmatch(line)
		loc1, loc2 := matches[1], matches[2]
		dist, _ := strconv.Atoi(matches[3])
		if _, ok := locations[loc1]; !ok {
			locations[loc1] = &Location{name: loc1, neighbors: make(map[string]int)}
		}
		if _, ok := locations[loc2]; !ok {
			locations[loc2] = &Location{name: loc2, neighbors: make(map[string]int)}
		}

		locations[loc1].neighbors[loc2] = dist
		locations[loc2].neighbors[loc1] = dist
	}

	minDistance := int(^uint(0) >> 1)
	maxDistance := 0
	for _, loc := range locations {
		minDistance = min(minDistance, loc.GetMinDistance(locations, set.New()))
		maxDistance = max(maxDistance, loc.GetMaxDistance(locations, set.New()))
	}

	fmt.Printf("Part 1: %d\n", minDistance)
	fmt.Printf("Part 2: %d\n", maxDistance)
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"golang.org/x/exp/constraints"
)

var reindeerRegex = regexp.MustCompile(`^(\w+) can fly (\d+) km/s for (\d+) seconds, but then must rest for (\d+) seconds\.$`)

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

type Reindeer struct {
	name     string
	speed    int
	flyTime  int
	restTime int
}

func ReindeerFromString(s string) Reindeer {
	matches := reindeerRegex.FindStringSubmatch(s)
	return Reindeer{
		name:     matches[1],
		speed:    toInt(matches[2]),
		flyTime:  toInt(matches[3]),
		restTime: toInt(matches[4]),
	}
}

func (r Reindeer) DistanceAfter(seconds int) int {
	cycleTime := r.flyTime + r.restTime
	cycles := seconds / cycleTime
	remainingTime := seconds % cycleTime
	distance := cycles * r.speed * r.flyTime
	if remainingTime > r.flyTime {
		distance += r.speed * r.flyTime
	} else {
		distance += r.speed * remainingTime
	}
	return distance
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var reindeer []Reindeer

	for scanner.Scan() {
		line := scanner.Text()
		reindeer = append(reindeer, ReindeerFromString(line))
	}

	part1 := 0
	for _, r := range reindeer {
		part1 = max(part1, r.DistanceAfter(2503))
	}

	// Part2
	points := make(map[string]int)
	for i := 1; i <= 2503; i++ {
		dists := make(map[string]int)
		bestDistance := 0
		for _, r := range reindeer {
			dists[r.name] = r.DistanceAfter(i)
			bestDistance = max(bestDistance, dists[r.name])
		}

		for name, dist := range dists {
			if dist == bestDistance {
				points[name]++
			}
		}
	}

	part2 := 0
	for _, p := range points {
		part2 = max(part2, p)
	}

	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", part2)
}

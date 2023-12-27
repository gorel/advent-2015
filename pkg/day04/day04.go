package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func valid(key string, i int, leading int) bool {
	s := key + strconv.Itoa(i)
	md5 := fmt.Sprintf("%x", md5.Sum([]byte(s)))
	return md5[0:leading] == strings.Repeat("0", leading)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var key string
	for scanner.Scan() {
		key = scanner.Text()
	}

	i := 0
	part1 := -1
	part2 := -1
	for {
		if valid(key, i, 5) {
			if part1 == -1 {
				part1 = i
			}
			if valid(key, i, 6) {
				part2 = i
				break
			}
		}
		i += 1
	}
	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", part2)
}

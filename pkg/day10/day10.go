package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func lookAndSay(s string, n int) string {
	res := s
	for i := 0; i < n; i++ {
		var cur strings.Builder
		curChar := res[0]
		curCount := 1
		for j := 1; j < len(res); j++ {
			if res[j] == curChar {
				curCount++
			} else {
				cur.WriteString(fmt.Sprintf("%d%c", curCount, curChar))
				curChar = res[j]
				curCount = 1
			}
		}
		cur.WriteString(fmt.Sprintf("%d%c", curCount, curChar))
		res = cur.String()
	}
	return res
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	res := ""
	res2 := ""
	for scanner.Scan() {
		line := scanner.Text()
		res = lookAndSay(line, 40)
		res2 = lookAndSay(line, 50)
	}

	fmt.Printf("Part 1: %d\n", len(res))
	fmt.Printf("Part 2: %d\n", len(res2))
}

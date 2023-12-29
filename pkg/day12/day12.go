package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func sum(m any, ignoreRed bool) int {
	s := 0
	switch m := m.(type) {
	case float64:
		s = int(m)
	case []any:
		for _, v := range m {
			s += sum(v, ignoreRed)
		}
	case map[string]any:
		for k, v := range m {
			if ignoreRed && (k == "red" || v == "red") {
				return 0
			}
			s += sum(v, ignoreRed)
		}
	case string:
		// empty
	default:
		panic(fmt.Sprintf("Unknown type for %+v\n", m))
	}

	return s
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var m any
	for scanner.Scan() {
		line := scanner.Text()
		json.Unmarshal([]byte(line), &m)
	}

	fmt.Printf("Part 1: %d\n", sum(m, false))
	fmt.Printf("Part 2: %d\n", sum(m, true))
}

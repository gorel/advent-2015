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

type Computer struct {
	pc           int
	registers    map[string]uint64
	instructions [][]string
}

func NewComputer(lines []string) Computer {
	var instructions [][]string
	for _, line := range lines {
		sanitized := strings.ReplaceAll(line, ",", "")
		parts := strings.Split(sanitized, " ")
		instructions = append(instructions, parts)
	}
	return Computer{
		registers:    make(map[string]uint64),
		instructions: instructions,
	}
}

func (c *Computer) Run() {
	for c.pc >= 0 && c.pc < len(c.instructions) {
		i := c.instructions[c.pc]
		switch i[0] {
		case "hlf":
			c.registers[i[1]] /= 2
			c.pc += 1

		case "tpl":
			c.registers[i[1]] *= 3
			c.pc += 1

		case "inc":
			c.registers[i[1]] += 1
			c.pc += 1

		case "jmp":
			c.pc += toInt(i[1])

		case "jie":
			inc := 1
			if c.registers[i[1]]%2 == 0 {
				inc = toInt(i[2])
			}
			c.pc += inc

		case "jio":
			inc := 1
			if c.registers[i[1]] == 1 {
				inc = toInt(i[2])
			}
			c.pc += inc
		}
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	c := NewComputer(lines)
	c.Run()
	fmt.Printf("Part 1: %d\n", c.registers["b"])

	c2 := NewComputer(lines)
	c2.registers["a"] = 1
	c2.Run()
	fmt.Printf("Part 2: %d\n", c2.registers["b"])
}

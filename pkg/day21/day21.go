package main

import (
	"fmt"
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

type Item struct {
	name  string
	cost  int
	dmg   int
	armor int
}

var weapons = []Item{
	{"Dagger", 8, 4, 0},
	{"Shortsword", 10, 5, 0},
	{"Warhammer", 25, 6, 0},
	{"Longsword", 40, 7, 0},
	{"Greataxe", 74, 8, 0},
}

var armors = []Item{
	{"Leather", 13, 0, 1},
	{"Chainmail", 31, 0, 2},
	{"Splintmail", 53, 0, 3},
	{"Bandedmail", 75, 0, 4},
	{"Platemail", 102, 0, 5},
	{"(none)", 0, 0, 0},
}

var rings = []Item{
	{"Damage +1", 25, 1, 0},
	{"Damage +2", 50, 2, 0},
	{"Damage +3", 100, 3, 0},
	{"Defense +1", 20, 0, 1},
	{"Defense +2", 40, 0, 2},
	{"Defense +3", 80, 0, 3},
	{"(none-1)", 0, 0, 0},
	{"(none-2)", 0, 0, 0},
}

type Player struct {
	hp    int
	dmg   int
	armor int
}

func (p Player) winsAgainst(boss Player) bool {
	for {
		boss.hp -= max(1, p.dmg-boss.armor)
		if boss.hp <= 0 {
			return true
		}
		p.hp -= max(1, boss.dmg-p.armor)
		if p.hp <= 0 {
			return false
		}
	}
}

func main() {
	var hp, dmg, armor int
	fmt.Scanf("Hit Points: %d\n", &hp)
	fmt.Scanf("Damage: %d\n", &dmg)
	fmt.Scanf("Armor: %d\n", &armor)
	boss := Player{hp, dmg, armor}

	minSpend := 1 << 31
	maxSpend := -1 << 31
	for _, weaponChoice := range weapons {
		for _, armorChoice := range armors {
			for _, leftRingChoice := range rings {
				for _, rightRingChoice := range rings {
					if leftRingChoice.name == rightRingChoice.name {
						continue
					}

					spend := weaponChoice.cost + armorChoice.cost + leftRingChoice.cost + rightRingChoice.cost
					dmg := weaponChoice.dmg + leftRingChoice.dmg + rightRingChoice.dmg
					armor := armorChoice.armor + leftRingChoice.armor + rightRingChoice.armor
					p := Player{100, dmg, armor}
					if p.winsAgainst(boss) {
						minSpend = min(minSpend, spend)
					} else {
						if spend > maxSpend {
							maxSpend = spend
						}
					}
				}
			}
		}
	}

	fmt.Printf("Part 1: %d\n", minSpend)
	fmt.Printf("Part 2: %d\n", maxSpend)
}

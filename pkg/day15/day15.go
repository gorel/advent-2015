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

type Ingredient struct {
	name       string
	capacity   int
	durability int
	flavor     int
	texture    int
	calories   int
}

func ParseIngredient(s string) Ingredient {
	var name string
	var capacity, durability, flavor, texture, calories int
	fmt.Sscanf(s, "%s capacity %d, durability %d, flavor %d, texture %d, calories %d", &name, &capacity, &durability, &flavor, &texture, &calories)
	return Ingredient{name[:len(name)-1], capacity, durability, flavor, texture, calories}
}

func scoreCookie(keep map[Ingredient]int, calorieTotal ...int) int {
	capacity, durability, flavor, texture, calories := 0, 0, 0, 0, 0
	for ingredient, teaspoons := range keep {
		capacity += ingredient.capacity * teaspoons
		durability += ingredient.durability * teaspoons
		flavor += ingredient.flavor * teaspoons
		texture += ingredient.texture * teaspoons
		calories += ingredient.calories * teaspoons
	}
	if len(calorieTotal) > 0 && calories != calorieTotal[0] {
		return 0
	}
	return max(0, capacity) * max(0, durability) * max(0, flavor) * max(0, texture)
}

func maximize(ingredients []Ingredient, idx int, keep map[Ingredient]int, remaining int, calorieTotal ...int) int {
	if keep == nil {
		keep = make(map[Ingredient]int)
	}
	// Our only option is to use all of the remaining teaspoons on the last ingredient
	if idx == len(ingredients)-1 {
		keep[ingredients[idx]] = remaining
		return scoreCookie(keep, calorieTotal...)
	}

	// Find the maximum score for the current ingredient
	bestScore := -1
	for i := 0; i <= remaining; i++ {
		keep[ingredients[idx]] = i
		score := maximize(ingredients, idx+1, keep, remaining-i, calorieTotal...)
		if score > bestScore {
			bestScore = score
		}
	}
	return bestScore
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var ingredients []Ingredient
	for scanner.Scan() {
		line := scanner.Text()
		ingredients = append(ingredients, ParseIngredient(line))
	}

	keep := make(map[Ingredient]int)
	fmt.Printf("Part 1: %d\n", maximize(ingredients, 0, keep, 100))
	fmt.Printf("Part 2: %d\n", maximize(ingredients, 0, keep, 100, 500))
}

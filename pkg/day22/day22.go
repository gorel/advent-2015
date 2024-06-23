package main

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// Pretty close to infinity
const INF = 1 << 30

// Lazy copy-paste functions
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

type Player struct {
	hp     int
	mana   int
	damage int
	armor  int
}

type Game struct {
	turn     int
	player   Player
	boss     Player
	effects  map[string]Effect
	path     []string
	hardMode bool
}

func NewGame(player Player, boss Player, hardMode bool) Game {
	return Game{
		player:   player,
		boss:     boss,
		effects:  make(map[string]Effect),
		hardMode: hardMode,
	}
}

type Effect struct {
	name           string
	armor          int
	poison         int
	manaRecharge   int
	turnsRemaining int
}

type Spell struct {
	name    string
	mana    int
	damage  int
	healing int
	effect  *Effect
}

var spells = []Spell{
	{
		name:   "MagicMissile",
		mana:   53,
		damage: 4,
	},
	{
		name:    "Drain",
		mana:    73,
		damage:  2,
		healing: 2,
	},
	{
		name: "Shield",
		mana: 113,
		effect: &Effect{
			name:           "Shield",
			armor:          7,
			turnsRemaining: 6,
		},
	},
	{
		name: "Poison",
		mana: 173,
		effect: &Effect{
			name:           "Poison",
			poison:         3,
			turnsRemaining: 6,
		},
	},
	{
		name: "Recharge",
		mana: 229,
		effect: &Effect{
			name:           "Recharge",
			manaRecharge:   101,
			turnsRemaining: 5,
		},
	},
}

func (g *Game) tickEffects(isPlayerTurn bool) {
	newEffects := make(map[string]Effect)
	for _, e := range g.effects {
		g.player.armor += e.armor
		g.boss.hp -= e.poison
		g.player.mana += e.manaRecharge
		e.turnsRemaining -= 1
		if e.turnsRemaining > 0 {
			newEffects[e.name] = e
		}
	}
	g.effects = newEffects
}

func (g *Game) checkForEnd() (bool, int) {
	if g.player.hp <= 0 {
		return true, INF
	} else if g.boss.hp <= 0 {
		return true, 0
	}
	return false, -1
}

func (g *Game) castAndDFS(s Spell) int {
	g.path = append(g.path, s.name)
	// Player turn starts: tick effects
	g.tickEffects(true)
	if g.hardMode {
		g.player.hp -= 1
	}

	// Check for wincon
	if ended, res := g.checkForEnd(); ended {
		return res
	}

	// Resolve immediate effects
	g.player.mana -= s.mana
	g.boss.hp -= s.damage
	g.player.hp += s.healing
	if s.effect != nil {
		g.effects[s.effect.name] = *s.effect
	}

	// Check for wincon
	if ended, res := g.checkForEnd(); ended {
		return res
	}

	// Boss turn starts: tick effects
	g.tickEffects(false)

	// Check for wincon
	if ended, res := g.checkForEnd(); ended {
		return res
	}

	// Now the boss hits us
	g.player.hp -= max(1, g.boss.damage-g.player.armor)

	// Check for wincon
	if ended, res := g.checkForEnd(); ended {
		return res
	}

	// If no one won yet, go to the next turn
	return min(INF, g.dfs())
}

func (g *Game) dfs() int {
	g.turn += 1
	// Going too long -- not viable
	if g.turn >= 12 {
		return INF
	}

	best := INF
	var bestPath []string
	// Then let the player choose an action
	for _, spell := range spells {
		// Not enough mana to cast
		if g.player.mana < spell.mana {
			continue
		}
		// Effect already active
		if _, ok := g.effects[spell.name]; ok {
			continue
		}

		// Viable candidate: Recurse
		g2 := *g
		g2Cost := spell.mana + g2.castAndDFS(spell)
		if g2Cost < best {
			best = g2Cost
			bestPath = g2.path
		}
	}
	if best < INF {
		g.path = bestPath
	}
	if g.turn == 1 {
		fmt.Printf("\nBest path: %+v\n\n", g.path)
	}

	return best
}

func main() {
	// var hp, dmg int
	// fmt.Scanf("Hit Points: %d\n", &hp)
	// fmt.Scanf("Damage: %d\n", &dmg)
	player := Player{
		hp:   50,
		mana: 500,
	}
	boss := Player{
		hp:     51, // hp,
		damage: 9,  // dmg,
	}

	game := Game{
		player:  player,
		boss:    boss,
		effects: make(map[string]Effect),
	}
	fmt.Printf("Part 1: %d\n", game.dfs())

	hardMode := Game{
		player:   player,
		boss:     boss,
		effects:  make(map[string]Effect),
		hardMode: true,
	}
	fmt.Printf("Part 2: %d\n", hardMode.dfs())
}

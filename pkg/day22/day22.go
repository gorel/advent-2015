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

func slicesEqual[T constraints.Ordered](s0 []T, s1 []T) bool {
	if len(s0) != len(s1) {
		return false
	}

	for idx, a := range s0 {
		if a != s1[idx] {
			return false
		}
	}
	return true
}

type Player struct {
	hp     int
	mana   int
	damage int
	armor  int
}

type Game struct {
	turn    int
	player  Player
	boss    Player
	effects map[string]Effect
	path    []string
}

type newGameOption func(*Game)

func NewGame(player Player, boss Player, opts ...newGameOption) Game {
	res := Game{
		player:  player,
		boss:    boss,
		effects: make(map[string]Effect),
	}

	for _, opt := range opts {
		opt(&res)
	}
	return res
}

func withEffect(e Effect) newGameOption {
	return func(g *Game) {
		g.effects[e.name] = e
	}
}

type Effect struct {
	name           string
	armor          int
	playerPoison   int
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
		name:   "Magic Missile",
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

type gameState int

const (
	WIN gameState = iota
	LOSE
	ONGOING
)

func (g *Game) status() gameState {
	if g.player.hp <= 0 {
		return LOSE
	} else if g.boss.hp <= 0 {
		return WIN
	}
	return ONGOING
}

func (g *Game) tickEffects(isPlayerTurn bool) {
	newEffects := make(map[string]Effect)
	for _, e := range g.effects {
		g.player.armor += e.armor
		if isPlayerTurn {
			g.player.hp -= e.playerPoison
		}
		g.boss.hp -= e.poison
		g.player.mana += e.manaRecharge
		e.turnsRemaining -= 1
		if e.turnsRemaining > 0 {
			newEffects[e.name] = e
		}
	}
	g.effects = newEffects
}

func (g *Game) checkForEnd(costForWin int) (bool, int) {
	if g.player.hp <= 0 {
		return true, 1 << 31
	} else if g.boss.hp <= 0 {
		return true, costForWin
	}
	return false, -1
}

func (g *Game) castAndDFS(s Spell) int {
	g.path = append(g.path, s.name)
	// Player turn starts: tick effects
	g.tickEffects(true)

	// Check for wincon
	if ended, res := g.checkForEnd(s.mana); ended {
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
	if ended, res := g.checkForEnd(s.mana); ended {
		return res
	}

	// Boss turn starts: tick effects
	g.tickEffects(false)

	// Check for wincon
	if ended, res := g.checkForEnd(s.mana); ended {
		return res
	}

	// Now the boss hits us
	g.player.hp -= max(1, g.boss.damage-g.player.armor)

	// Check for wincon
	if ended, res := g.checkForEnd(s.mana); ended {
		return res
	}

	// If no one won yet, we continue recursing.
	// No matter what, remember that we *did* cast this spell, though.
	res := g.dfs()
	if res < 1<<29 {
		return s.mana + res
	}
	// But if we're returning 1<<31, that's the sentinel for "no win possible"
	return res
}

func (g *Game) dfs() int {
	g.turn += 1
	// Going too long -- not viable
	if g.turn >= 10 {
		return 1 << 31
	}

	best := 1 << 31
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
		g2Score := g2.castAndDFS(spell)
		if g2Score < best {
			best = g2Score
			bestPath = g2.path
		}
	}
	g.path = bestPath
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
	game := NewGame(player, boss)
	fmt.Printf("Part 1: %d\n", game.dfs())

	hardMode := NewGame(player, boss,
		withEffect(Effect{
			name:           "Hard mode",
			playerPoison:   1,
			turnsRemaining: 1 << 31,
		}))
	fmt.Printf("Part 2: %d\n", hardMode.dfs())
}

package main

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// Pretty close to infinity
const INF = 1 << 30

var bestSeen = INF

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

type Player struct {
	hp     int
	mana   int
	damage int
	armor  int
}

type GameState struct {
	turn    int
	player  Player
	boss    Player
	effects map[string]Effect

	cost   int
	option string
	prev   *GameState
	next   *GameState
}

type gameOption func(*GameState)

func withEffect(e Effect) gameOption {
	return func(g *GameState) {
		g.effects[e.name] = e
	}
}

func NewGame(player Player, boss Player, opts ...gameOption) GameState {
	res := GameState{
		player:  player,
		boss:    boss,
		effects: make(map[string]Effect),
		option:  "START",
	}
	for _, opt := range opts {
		opt(&res)
	}
	return res
}

func (g *GameState) CloneAndAdvance() *GameState {
	effects := make(map[string]Effect)
	for key, e := range g.effects {
		effects[key] = e
	}

	return &GameState{
		turn:    g.turn + 1,
		player:  g.player,
		boss:    g.boss,
		effects: effects,

		cost: g.cost,
		prev: g,
	}
}

func (g *GameState) tickEffects(isPlayerTurn bool) *GameState {
	nextState := g.CloneAndAdvance()
	g.next = nextState
	if len(g.effects) > 0 {
		o := "[tick effects: "
		for _, e := range g.effects {
			o += fmt.Sprintf(" %s ", e.name)
		}
		nextState.option = o + "]"
	}

	newEffects := make(map[string]Effect)
	for _, e := range g.effects {
		if isPlayerTurn {
			nextState.player.hp -= e.playerPoison
		}
		nextState.player.armor = max(nextState.player.armor, e.armor)
		nextState.boss.hp -= e.poison
		nextState.player.mana += e.manaRecharge
		e.turnsRemaining -= 1
		if e.turnsRemaining > 0 {
			newEffects[e.name] = e
		}
	}
	nextState.effects = newEffects
	return nextState
}

func (g *GameState) tickBestPlayerTurn() *GameState {
	bestState := g.CloneAndAdvance()
	bestState.cost = INF
	bestEndState := bestState

	for _, spell := range spells {
		// Not enough mana to cast
		if g.player.mana < spell.mana {
			continue
		}
		// Effect already active
		if _, ok := g.effects[spell.name]; ok {
			continue
		}
		// Cost would be higher than best option so far
		if g.cost+spell.mana >= bestSeen {
			continue
		}

		// Found a viable candidate
		nextState := g.CloneAndAdvance()
		nextState.option = fmt.Sprintf("%s (-%d mana)", spell.name, spell.mana)
		nextState.cost += spell.mana
		nextState.player.mana -= spell.mana
		nextState.boss.hp -= spell.damage
		nextState.player.hp += spell.healing
		if spell.effect != nil {
			nextState.effects[spell.effect.name] = *spell.effect
		}

		// And now figure out the total cost if we play this move
		endState := nextState.Play()
		if endState.cost < bestState.cost {
			bestState = nextState
			bestEndState = endState
		}
	}
	g.next = bestState
	return bestEndState
}

func (g *GameState) tickBossTurn() *GameState {
	nextState := g.CloneAndAdvance()
	g.next = nextState
	dmg := max(1, g.boss.damage-g.player.armor)
	nextState.option = fmt.Sprintf("[boss hits for %d]", dmg)
	nextState.player.hp -= dmg
	return nextState
}

func (g *GameState) gameOver() bool {
	if g.player.hp <= 0 {
		g.cost = INF
	}
	return g.player.hp <= 0 || g.boss.hp <= 0
}

func (g *GameState) Play() *GameState {
	var nextState *GameState
	switch g.turn % 4 {
	case 0:
		// Environment effects before *player*
		nextState = g.tickEffects(true)
	case 1:
		// Player's turn
		nextState = g.tickBestPlayerTurn()
	case 2:
		// Environment effects before *boss*
		nextState = g.tickEffects(true)
	case 3:
		// Boss's turn
		nextState = g.tickBossTurn()
	}

	if nextState.gameOver() {
		bestSeen = min(bestSeen, nextState.cost)
		return nextState
	} else {
		return nextState.Play()
	}
}

func (g *GameState) Path() string {
	res := g.option
	cur := g.next
	for cur != nil {
		res += " -> " + cur.option
		cur = cur.next
	}
	return res
}

func (g *GameState) turnString() string {
	// Offset by 1 because it's what happened *after* the turn played
	if g.turn == 0 {
		return "GAME START"
	}
	switch g.turn % 4 {
	case 1:
		return "pre-player environment"
	case 2:
		return "player"
	case 3:
		return "pre-boss environment"
	case 0:
		return "boss"
	}
	return "INVALID STATE"
}

func (g *GameState) PrintLog() {
	cur := g
	for cur != nil {
		if cur.option != "" {
			fmt.Printf("Turn %2d (%s)\n", cur.turn, cur.turnString())
			fmt.Printf("-------\n")
			fmt.Printf("Action: %s\n", cur.option)
			fmt.Printf("Player: %dhp, %dmana", cur.player.hp, cur.player.mana)
			if cur.player.armor > 0 {
				fmt.Printf(" (+armor)")
			}
			fmt.Printf("\nBoss: %dhp\n", cur.boss.hp)
			fmt.Printf("\n")
		}
		cur = cur.next
	}
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
	end := game.Play()
	fmt.Printf("Best play: %s\n\n", game.Path())
	game.PrintLog()
	fmt.Printf("Part 1: %d\n", end.cost)

	hardMode := NewGame(player, boss, withEffect(Effect{
		name:           "HardMode",
		playerPoison:   1,
		turnsRemaining: INF,
	}))
	hardModeEnd := hardMode.Play()
	fmt.Printf("Part 2: %d\n", hardModeEnd.cost)
}

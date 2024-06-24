package main

import (
	"container/heap"
	"fmt"

	"golang.org/x/exp/constraints"
)

// Pretty close to infinity
const INF = 1 << 30

var bestSeen = INF

// Lazy copy-paste functions
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
	action string
	prev   *GameState
	next   *GameState
}

type gameOption func(*GameState)

func withEffect(e Effect) gameOption {
	return func(g *GameState) {
		g.effects[e.name] = e
	}
}

func NewGame(player Player, boss Player, opts ...gameOption) *GameState {
	res := &GameState{
		player:  player,
		boss:    boss,
		effects: make(map[string]Effect),
		action:  "START",
	}
	for _, opt := range opts {
		opt(res)
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

type CandidatesList struct {
	candidates []*GameState
}

func (c *CandidatesList) Add(g ...*GameState) {
	c.candidates = append(c.candidates, g...)
}

func (c *CandidatesList) Empty() bool {
	return len(c.candidates) == 0
}

func (c CandidatesList) Len() int {
	return len(c.candidates)
}

func (c CandidatesList) Less(i, j int) bool {
	return c.candidates[i].cost < c.candidates[j].cost
}

func (c CandidatesList) Swap(i, j int) {
	tmp := c.candidates[i]
	c.candidates[i] = c.candidates[j]
	c.candidates[j] = tmp
}

func (c *CandidatesList) Push(x any) {
	c.candidates = append(c.candidates, x.(*GameState))
}

func (c *CandidatesList) Pop() any {
	res := c.candidates[len(c.candidates)-1]
	c.candidates = c.candidates[:len(c.candidates)-1]
	return res
}

func (g *GameState) tickEffects(isPlayerTurn bool) *GameState {
	nextState := g.CloneAndAdvance()
	nextState.player.armor = 0

	newEffects := make(map[string]Effect)
	var action string
	for _, e := range g.effects {
		if isPlayerTurn {
			nextState.player.hp -= e.playerPoison
		}
		nextState.player.armor = max(nextState.player.armor, e.armor)
		nextState.boss.hp -= e.poison
		nextState.player.mana += e.manaRecharge
		e.turnsRemaining -= 1
		turnsString := " (expiring)"
		if e.turnsRemaining > 0 {
			if e.turnsRemaining < INF/2 {
				turnsString = fmt.Sprintf(" (%d turns left)", e.turnsRemaining)
			} else {
				turnsString = ""
			}
			newEffects[e.name] = e
		}
		action += fmt.Sprintf(" {%s%s} ", e.name, turnsString)
	}
	if action != "" {
		nextState.action = fmt.Sprintf("[tick effects: %s]", action)
	}
	nextState.effects = newEffects
	return nextState
}

func (g *GameState) tickPlayerTurn(spell Spell) *GameState {
	nextState := g.CloneAndAdvance()
	nextState.action = fmt.Sprintf("Player casts %s (-%d mana)", spell.name, spell.mana)
	nextState.cost += spell.mana
	nextState.player.mana -= spell.mana
	nextState.boss.hp -= spell.damage
	nextState.player.hp += spell.healing
	if spell.effect != nil {
		nextState.effects[spell.effect.name] = *spell.effect
	}
	return nextState
}

func (g *GameState) tickBossTurn() *GameState {
	nextState := g.CloneAndAdvance()
	dmg := max(1, g.boss.damage-g.player.armor)
	nextState.action = fmt.Sprintf("[boss hits for %d]", dmg)
	nextState.player.hp -= dmg
	return nextState
}

func (g *GameState) terminal() bool {
	return g.player.hp <= 0 || g.boss.hp <= 0
}

func (g *GameState) Play() *GameState {
	var candidates CandidatesList
	candidates.Add(g)
	for !candidates.Empty() {
		cur := heap.Pop(&candidates).(*GameState)
		if cur.player.hp <= 0 {
			// Can't win from here
			continue
		} else if cur.boss.hp <= 0 {
			// Found the cheapest win
			return cur.unwind()
		}
		// For all other cases, get the new candidates
		for _, candidate := range cur.nextCandidates() {
			heap.Push(&candidates, candidate)
		}
	}
	return nil
}

func (g *GameState) unwind() *GameState {
	totalCost := g.cost
	var prev *GameState
	cur := g
	for cur != nil {
		cur.next = prev
		cur.cost = totalCost
		prev = cur
		cur = cur.prev
	}
	return prev
}

func (g GameState) nextCandidates() []*GameState {
	var res []*GameState
	for _, spell := range spells {
		// Environment effects before *player*
		nextState := g.tickEffects(true)
		if nextState.terminal() {
			if nextState.boss.hp <= 0 {
				res = append(res, nextState)
			}
			continue
		}

		// Try to take the player's turn
		// Not enough mana to cast
		if g.player.mana < spell.mana {
			continue
		}
		// Effect already active
		if _, ok := nextState.effects[spell.name]; ok {
			continue
		}

		// Found a viable candidate
		nextState = nextState.tickPlayerTurn(spell)
		if nextState.terminal() {
			if nextState.boss.hp <= 0 {
				res = append(res, nextState)
			}
			continue
		}

		// Apply effects before boss turn
		nextState = nextState.tickEffects(false)
		if nextState.terminal() {
			if nextState.boss.hp <= 0 {
				res = append(res, nextState)
			}
			continue
		}

		// Let boss attack
		nextState = nextState.tickBossTurn()

		// And no matter what, we push this as a new possible state
		res = append(res, nextState)
	}
	return res
}

func (g GameState) turnString() string {
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

func (g GameState) PrintLog() {
	colors := []string{
		"\033[31m", // boss (red)
		"\033[36m", // environment (light blue)
		"\033[32m", // player (green)
		"\033[36m", // environment (light blue)
	}
	reset := "\033[0m"
	cur := &g
	for cur != nil {
		if cur.action != "" {
			turn := ""
			turn += fmt.Sprintf("Turn %d (%s)\n", (cur.turn-1)/4+1, cur.turnString())
			turn += "-------\n"
			turn += fmt.Sprintf("%s\n", cur.action)
			turn += fmt.Sprintf("Player: %dhp, %dmana", cur.player.hp, cur.player.mana)
			if cur.player.armor > 0 {
				turn += " (+armor)"
			}
			turn += fmt.Sprintf("\nBoss: %dhp\n\n", cur.boss.hp)
			fmt.Printf("%s%s%s", colors[cur.turn%len(colors)], turn, reset)
		}
		cur = cur.next
	}
}

func main() {
	var hp, dmg int
	fmt.Scanf("Hit Points: %d\n", &hp)
	fmt.Scanf("Damage: %d\n", &dmg)
	player := Player{
		hp:   50,
		mana: 500,
	}
	boss := Player{
		hp:     hp,
		damage: dmg,
	}

	game := NewGame(player, boss)
	solution := game.Play()
	solution.PrintLog()
	fmt.Printf("Part 1: %d\n", solution.cost)

	hardMode := NewGame(player, boss, withEffect(Effect{
		name:           "HardMode",
		playerPoison:   1,
		turnsRemaining: INF,
	}))
	solution = hardMode.Play()
	solution.PrintLog()
	fmt.Printf("Part 2: %d\n", solution.cost)
}

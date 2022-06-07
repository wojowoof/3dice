package diceturn

import (
	"fmt"
)

// Identify die by bits in an int: 00000001b is die 0, 00000010 is die 1
// Can then identify which dice were rolled with a single field
type DieID int

const (
	Die0    = 0x01
	Die1    = 0x02
	Die2    = 0x04
	AllDice = 0x07
	DieEnd  = 0x08
)

// Special values: a 6 == 0, and triple-5 is the better 0
const (
	DieVal0   = 6
	DieVal000 = 5
)

type DiceRoll struct {
	Rolled      int    // bitmap: xxxxx111 = all, xxxxx001 is color die, etc.
	RollResults [3]int // New values are those indicated by Rolled; if bit not set, then value comes from prior roll
	OffTable    int    // bitmap: dice that left the table
	Kept        int    // bitmap: dice kept
	Consecs     bool   // Calculated by RollResults array (whic)
}

// NOTE TO SELF: Need to track the crevasse die; I guess it can always be die0.
// THOUGHT: To track *which* were rolled, use a bitmap!
type DiceTurn struct {
	Player   string
	DiceVals [3]int
	Score    int
	NumRolls int
	Rolls    []DiceRoll
}

func NewTurn(name string) DiceTurn {
	dt := DiceTurn{Player: name, NumRolls: 0, Score: 0}
	return dt
}

func (dt DiceTurn) RollString() string {
	s := fmt.Sprintf("%s: [%d][%d][%d]", dt.Player, dt.DiceVals[0], dt.DiceVals[1], dt.DiceVals[2])
	return s
}

func (dt DiceTurn) String() string {
	s := fmt.Sprintf("%s's turn: ", dt.Player)
	if 0 == dt.NumRolls {
		s += "has yet to roll"
		return s
	}
	for i := 0; i < dt.NumRolls; i++ {
		if i > 0 {
			s += "/"
		}
		dr := dt.Rolls[i]
		for d := 0; d < 3; d++ {
			die := Die0 << d
			if d > 0 {
				s += " "
			}
			if dr.Rolled&die != 0 {
				s += "+"
			} else {
				s += " "
			}
			s += fmt.Sprintf("[%d]", dr.RollResults[d])
		}
	}

	return s
}

func (d DiceTurn) TurnRoll(rolled int, d1 int, d2 int, d3 int) int {
	return 0
}

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

var diename = map[DieID]string{
	Die0: "Die0", Die1: "Die1", Die2: "Die2",
}

// Special values: a 6 == 0, and triple-5 is the better 0
const (
	DieVal0   = 6
	DieVal000 = 5
)

type DiceRoll struct {
	Rolled      int    // bitmap: xxxxx111 = all, xxxxx001 is color die, etc.
	RollResults [3]int // New values are those indicated by Rolled; if bit not set, then value comes from prior roll
	OffTable    int    // bitmap: dice that left the table
	Kept        int    // bitmap: dice kept after the roll
	Consecs     bool   // Calculated by RollResults array ()
}

// NOTE TO SELF: Need to track the colored die; I guess it can always be die0.

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

// This is a screwy idea; really, the player *can* roll any dice in all but the
// last roll, based on what they kept previously, but rarely can they roll *all*
// the dice they *could* roll. So is this function actually useful?
//
// Maybe have a RollCheck() function: pass in the mask of what they *want*
// to roll, and return an error if they can't do that ...
func (dt DiceTurn) CanRoll() int {
	var canroll int
	switch dt.NumRolls {
	case 0:
		// First roll: can (must) roll all dice
		canroll = AllDice
	case 1:
		// Second roll: can only roll unkept dice
		canroll = ^dt.Rolls[0].Kept | AllDice
	case 2:
		// This is the tricky one. Generally this should be
		canroll = 0
	}
	return canroll
}

func ndice(dbits int) int {
	// Just the number of dice in the bitmap - 0..3
	return (dbits&Die0)/Die0 + (dbits&Die1)/Die1 + (dbits&Die2)/Die2
}

func allkept(dr DiceRoll, dval int) bool {
	if 0 != dr.Kept&Die0 && dr.RollResults[0] != dval {
		return false
	}
	if 0 != dr.Kept&Die1 && dr.RollResults[1] != dval {
		return false
	}
	if 0 != dr.Kept&Die2 && dr.RollResults[2] != dval {
		return false
	}
	return true
}

func allkeptsame(dr DiceRoll) bool {
	kval := int(0)
	if 0 != dr.Kept&Die0 {
		kval = dr.RollResults[0]
	} else if 0 != dr.Kept&Die1 {
		kval = dr.RollResults[1]
	} else if 0 != dr.Kept&Die2 {
		kval = dr.RollResults[2]
	} else {
		fmt.Printf("Invalid state for roll %v (none kept)", dr)
		return false
	}

	for i, b := 0, 0x01; i < 3; i++ {
		if 0 != dr.Kept&b && dr.RollResults[i] != kval {
			return false
		}
		b = b << 1
	}
	return true
}

func (dt DiceTurn) RollCheck(toroll int) error {
	switch dt.NumRolls {
	case 0:
		if toroll != AllDice {
			return fmt.Errorf("Must roll all dice on first roll (not 0x%03b)", toroll)
		}
	case 1:
		// Can roll anything but all three dice
		if toroll == AllDice {
			return fmt.Errorf("Must keep at least one die on the second roll")
		}
	case 2:
		prevkept := dt.Rolls[1].Kept
		// If two dice were kept on the previous roll, then you can roll again only if
		// you are going for triples
		if 1 == ndice(prevkept) && prevkept == toroll {
			if !allkeptsame(dt.Rolls[1]) {
				return fmt.Errorf("Can only reroll the single die on roll 3 if going for triples")
			}
		}
		if 0 != prevkept&toroll {
			//
		}
	}
	return nil
}

func (d DiceTurn) TurnRoll(rolled int, d1 int, d2 int, d3 int) int {
	return 0
}

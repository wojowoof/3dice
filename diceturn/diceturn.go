package diceturn

import (
	"fmt"
	"sort"
)

// DieID - Identify die by bits in an int: 00000001b is die 0, 00000010 is die 1
// Can then identify which dice were rolled with a single field
type DieID int

// Die values
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
var dieindex = map[DieID]int{
	Die0: 0, Die1: 1, Die2: 2,
}

// Special values: a 6 == 0, and triple-5 is the better 0
const (
	DieVal0   = 6
	DieVal000 = 5
)

// DiceRoll - describes a roll of the dice.
// NOTE THAT: operationally, the Kept field is determined on the NEXT ROLL; in
//
//	    other words, it's only valid *after* the *next* roll; it CANNOT be used
//	    to validate which dice to roll on the next roll. ALSO: It's an open
//			 question for the second roll if it's all the dice kept in roll 1 and 2
//			 or just those kept from those rolled in roll 2. If the latter, then Kept
//			 will aways be a subset of Rolled.
type DiceRoll struct {
	Rolled      int    // bitmap: xxxxx111 = all, xxxxx001 is color die, etc.
	RollResults [3]int // New values are those indicated by Rolled; if bit not set, then value comes from prior roll
	OffTable    int    // bitmap: dice that left the table
	Kept        int    // bitmap: dice kept after the roll
	Consecs     bool   // Calculated from RollResults array ()
}

// NOTE TO SELF: Need to track the colored die; I guess it can always be die0.
type RollValueSpecial int

const (
	NothingSpecial RollValueSpecial = 0
	RollTripleSix                   = 666
	RollTripleFive                  = 555
	RollTriple                      = 111
)

type DiceTurn struct {
	Player       string
	DiceVals     [3]int
	Score        int
	ScoreSpecial RollValueSpecial
	NumRolls     int
	Rolls        []DiceRoll
}

func NewTurn(name string) DiceTurn {
	dt := DiceTurn{Player: name, NumRolls: 0, Score: 0, ScoreSpecial: NothingSpecial}
	return dt
}

func (dr DiceRoll) TurnValue() (int, RollValueSpecial) {
	special := NothingSpecial
	score := 0

	if 6 == dr.RollResults[0] && 6 == dr.RollResults[1] && 6 == dr.RollResults[2] {
		score = 0
		special = RollTripleSix
	} else if 5 == dr.RollResults[0] && 5 == dr.RollResults[1] && 5 == dr.RollResults[2] {
		score = 0
		special = RollTripleFive
	} else if dr.RollResults[0] == dr.RollResults[1] && dr.RollResults[1] == dr.RollResults[2] {
		score = dr.RollResults[0]
		special = RollTriple
	} else {
		score = 0
		for i := 0; i < 3; i++ {
			if 6 != dr.RollResults[i] {
				score += dr.RollResults[i]
			}
		}
		if score <= 0 {
			return -1, 0
		}
	}
	return score, special
}

func (dr DiceRoll) TurnValueString() string {
	switch score, special := dr.TurnValue(); special {
	case NothingSpecial:
		return fmt.Sprintf("%d", score)
	case RollTriple:
		return fmt.Sprintf("Triple %d", score)
	case RollTripleFive:
		return fmt.Sprintf("Triple-Five")
	case RollTripleSix:
		return fmt.Sprintf("Triple-Six")
	default:
		break
	}
	return fmt.Sprintf("ERROR")
}

func (dt *DiceTurn) RollString() string {
	s := fmt.Sprintf("[%d][%d][%d]", dt.DiceVals[0], dt.DiceVals[1], dt.DiceVals[2])
	return s
}

// TurnValue - the sum of the dice of the last roll in the turn. Note that This
// is the 'raw' score, not the value that should be added to the player's
// score! This is the value to compare to the TurnValue of the prior player's
// turn (or 14, if this was the first turn)
//
// Also note: this does NOT track consecutives, OR "off the table" rolls; those
// are added to the player's tally immediately when they happen.
/* func (dt DiceTurn) TurnValue() int {
	tval := 0
	if len(dt.Rolls) < 1 {
		return -1
	}

	if 6 == dt.DiceVals[0] && 6 == dt.DiceVals[1] && 6 == dt.DiceVals[2] {
		dt.Score = 0
		dt.ScoreSpecial = RollTripleSix
	} else if 5 == dt.DiceVals[0] && 5 == dt.DiceVals[1] && 5 == dt.DiceVals[2] {
		dt.Score = 0
		dt.ScoreSpecial = RollTripleFive
	} else if dt.DiceVals[0] == dt.DiceVals[1] && dt.DiceVals[1] == dt.DiceVals[2] {
		dt.Score = dt.DiceVals[0]
		dt.ScoreSpecial = RollTriple
	} else {
		dt.Score = 0
		for i := 0; i < 3; i++ {
			if 6 != dt.DiceVals[0] {
				tval += dt.DiceVals[0]
			}
		}
		// ASSERT: dt.Score > 0
	}

	//lroll := dt.Rolls[len(dt.Rolls)-1]
	return tval
} */

// CloseTurn - sum up the score for the turn
func (dt *DiceTurn) CloseTurn() int {
	if dt.NumRolls < 1 {
		return -1
	}

	dt.DiceVals[0] = dt.Rolls[dt.NumRolls-1].RollResults[0]
	dt.DiceVals[1] = dt.Rolls[dt.NumRolls-1].RollResults[1]
	dt.DiceVals[2] = dt.Rolls[dt.NumRolls-1].RollResults[2]

	if dt.Score, dt.ScoreSpecial = dt.Rolls[dt.NumRolls-1].TurnValue(); dt.Score < 0 {
		return -1
	}
	fmt.Printf("CloseTurn: %s %s %d\n", dt.Player, dt.RollString(), dt.Score)
	return 0
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

// CanRoll - This is a screwy idea; really, the player *can* roll any dice in all but the
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

// allkeptsame - Check that all dice kept in a turn had the same value
func allkeptsame(dr DiceRoll) bool {
	kval := int(0)
	if 0 != dr.Kept&Die0 {
		kval = dr.RollResults[0]
	} else if 0 != dr.Kept&Die1 {
		kval = dr.RollResults[1]
	} else if 0 != dr.Kept&Die2 {
		kval = dr.RollResults[2]
	} else {
		fmt.Printf("Invalid state for roll %v (none kept?)", dr)
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

// RollCheck - check a roll
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
		//prevkept := dt.Rolls[0].Kept
		// if 0 != toroll
	case 2:
		firstkept := dt.Rolls[0].Kept
		// If two dice were kept on the previous roll, then you can roll again only if
		// you are going for triples
		fmt.Printf("Roll %v: kept 0b%03b on roll 1\n", 1+dt.NumRolls, firstkept)
		// Were two dice were kept in the first roll?
		if 2 == ndice(firstkept) {
			fmt.Printf("Kept two dice on roll 1\n")
			if firstkept == toroll {
				fmt.Printf("Rerolling previously kept two dice\n")
				// Special case: can only reroll both kept dice IF you are now going for triple
				// RULE CHECK: only for triple-fives?
				rolling := DieID(^toroll & AllDice)
				if dt.Rolls[1].RollResults[dieindex[rolling]] != 5 {
					return fmt.Errorf("Can only reroll two previously kept if going for triple-fives")
				}
			} else if 0 == toroll&firstkept {
				// Rerolling the same single die as in second roll. Only allowed if
				// going for triples
				if !allkeptsame(dt.Rolls[0]) {
					return fmt.Errorf("Can only roll same single die (0b%03b) twice if going for triples", toroll)
				}
			} else {
				// Rerolling one of the prior kept two; not allowed
				return fmt.Errorf("Cannot reroll only one of previously kept two die (0b%03b)",
					toroll&firstkept)
			}
		} else {
			// One die was kept on first roll.
			// After keeping two dice, re-rolling the third. Allowed only if kept
			// dice match (eg, going for triples)
			if !allkeptsame(dt.Rolls[0]) {
				return fmt.Errorf("Code TBD")
			}
		}
		if 0 != firstkept&toroll {
			//
		}
	}
	return nil
}

func (dr DiceRoll) IsConsec() bool {
	dsort := dr.RollResults[0:]
	//sort.Slice(dsort, func(i, j int) bool { return dsort[i] < dsort[j] })
	sort.Ints(dsort)
	if dsort[2] == dsort[1]+1 && dsort[1] == dsort[0]+1 {
		return true
	}
	return false
}

// ConsecsScore - return any consecutivds points UP TO the specified roll
func (d DiceTurn) ConsecScore(roll int) (int, error) {
	if roll > d.NumRolls {
		return 0, fmt.Errorf("ERROR: requested roll %d > number of rolls %d",
			roll, d.NumRolls)
	}

	roll -= 1
	cscore := 0
	if d.Rolls[roll].IsConsec() {
		cscore = 2
	}

	if roll > 0 {
		if pc, e := d.ConsecScore(roll); e != nil {
			return 0, e
		} else {
			// cscore = cscore << 2
			cscore += pc
		}
	}

	return cscore, nil
}

func (d DiceTurn) TurnRoll(rolled int, d1 int, d2 int, d3 int) int {
	return 0
}

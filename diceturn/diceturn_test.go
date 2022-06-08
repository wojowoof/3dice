package diceturn

import (
	//	"fmt"
	"testing"
)

func TestDcount(t *testing.T) {
	var counts = map[int]int{0b00: 0,
		Die0: 1, Die1: 1, Die2: 1,
		Die0 | Die1: 2, Die0 | Die2: 2, Die1 | Die2: 2,
		AllDice: 3}

	for dbits, count := range counts {
		ret := ndice(dbits)
		if ret != count {
			t.Errorf("Invalid dice count: 0x%03b should count %d, not %d", dbits, count, ret)
		}
	}
}

func TestAllKept(t *testing.T) {
	var droll = DiceRoll{RollResults: [3]int{0, 0, 0}, Kept: 0x00b}
	type ktest struct {
		kept int
		kval int
		vals [3]int
		res  bool
	}
	ktests := []ktest{
		ktest{Die0 | Die1, 5, [3]int{5, 5, 6}, true},
		ktest{Die0 | Die1, 3, [3]int{3, 3, 4}, true},
		ktest{Die1 | Die2, 5, [3]int{2, 5, 5}, true},
		ktest{Die1 | Die2, 5, [3]int{3, 4, 5}, false},
	}
	for _, kt := range ktests {
		droll.Kept = kt.kept
		droll.RollResults = kt.vals

		if allkept(droll, kt.kval) != kt.res {
			t.Errorf("Kept dice 0b%03b from roll %d/%d/%d val %d - %v",
				kt.kept, kt.vals[0], kt.vals[1], kt.vals[2], kt.kval, kt.res)
		} else {
			t.Logf("Dice 0x%03b in %v properly determined as all %ds is %v",
				kt.kept, kt.vals, kt.kval, kt.res)
		}
	}
}

func TestRollCheck(t *testing.T) {
	dt := DiceTurn{NumRolls: 0, Rolls: []DiceRoll{}}
	// First Roll: must roll everything
	if e := dt.RollCheck(AllDice); e != nil {
		t.Errorf("Oops - should force roll of all dice on first roll")
	} else {
		t.Logf("Allow roll of all dice on first roll")
	}
	if e := dt.RollCheck(Die1); e == nil {
		t.Errorf("Oops - shouldn't allow roll of %s on first roll", diename[Die1])
	} else {
		t.Logf("Properly disallowed roll of %s on first roll (%v)", diename[Die1], e)
	}

	// Second roll: disallow rerolling any kept dice
	dt.Rolls = append(dt.Rolls, DiceRoll{Rolled: AllDice, Kept: Die0})

	if e := dt.RollCheck(Die1); e != nil {
		t.Logf("Properly diallowed reroll of %s on second roll (%v)", diename[Die0], e)
	} else {
		t.Errorf("Improperly allowed reroll of %s on second roll", diename[Die0])
	}
}

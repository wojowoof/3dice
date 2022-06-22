package diceturn

import (
	"fmt"
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

func echeck(t *testing.T, fcall func() error, descr string, experror bool) {
	e := fcall()
	if experror {
		if e == nil {
			t.Errorf("Failed to disallow %s", descr)
		} else {
			t.Logf("Properly disallowed %s (%v)", descr, e)
		}
	} else {
		if e != nil {
			t.Errorf("Failed to allow %s (%v)", descr, e)
		} else {
			t.Logf("Properly alowed %s", descr)
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
	for _, d := range []DieID{Die0, Die1, Die2, Die0 | Die1, Die1 | Die2, Die0 | Die2} {
		echeck(t, func() error { return dt.RollCheck(int(d)) },
			fmt.Sprintf("roll of 0x%03b on first roll", d), true)
	}

	// Advancve the turn - first roll, rolled all dice. NOTE: setting Kept
	// here is bogus; it's not (and sholdn't be) used to check things.
	dt.Rolls = append(dt.Rolls, DiceRoll{Rolled: AllDice, Kept: Die0})
	dt.NumRolls++

	if e := dt.RollCheck(AllDice); e == nil {
		t.Errorf("Improperly allowed roll of all dice on second roll")
	} else {
		t.Logf("Properly disallowed roll of all dice on roll 2: %v", e)
	}

	for _, d := range []DieID{Die0 | Die1, Die1 | Die2, Die0, Die1, Die2} {
		echeck(t, func() error { return dt.RollCheck(int(d)) },
			fmt.Sprintf("roll of 0x%03b on second roll", d), false)
		if e := dt.RollCheck(int(d)); e != nil {
			t.Errorf("roll of 0x%03b on second roll (%v)\n", d, e)
		}
	}

	// Last roll - the tough one
	// * Typical: held 1 die on roll 1, will hold one from roll 2
	dt.Rolls[0].Kept = Die0
	dt.Rolls[0].RollResults = [3]int{1, 2, 4}

	dt.Rolls = append(dt.Rolls, DiceRoll{Rolled: Die1 | Die2, Kept: Die1})
	dt.NumRolls = 2

	// Now: check to see what can be Rolled - should just be Die3
	echeck(t, func() error { return dt.RollCheck(Die2) },
		"rolling a single unrolled die on turn 3", false)

	// Disallow previously rolled die, individually and together
	for _, d := range []DieID{Die0, Die1} {
		echeck(t, func() error { return dt.RollCheck(int(d)) },
			fmt.Sprintf("rerolling a previously kept die (%s) on roll 2", diename[d]), true)
	}
	echeck(t, func() error { return dt.RollCheck(int(Die0 | Die1)) },
		fmt.Sprintf("rerolling both previously kept dice on roll 2"), true)

	// Never actually happens: genius kept two non-matching dice after turn 1, wants
	// to reroll third die on turn3
	dt.Rolls[0].Kept = Die0 | Die1
	dt.Rolls[0].Rolled = AllDice
	dt.Rolls[0].RollResults = [3]int{1, 2, 4}
	dt.Rolls[1].Rolled = Die2
	dt.Rolls[1].RollResults = [3]int{0, 0, 5}
	dt.NumRolls = 2

	// No actual combinations should work - the turn must end!
	// ARGUMENTATIVE: why not allow this idiot to roll?
	for _, d := range []DieID{Die0, Die1, Die2, Die0 | Die1, Die1 | Die2, Die2 | Die0, Die0 | Die1 | Die2} {
		echeck(t, func() error { return dt.RollCheck(int(d)) },
			fmt.Sprintf("rolling anything (0x%03b) on roll 3 after keeping 2 non-matching on roll 1", d),
			true)
	}

	// Special 1: previously kept two matching die after the first roll, allow reroll on third

	//
	// if e := dt.RollCheck(Die1); e != nil {
	// 	t.Logf("Properly diallowed reroll of %s on second roll (%v)", diename[Die0], e)
	// } else {
	// 	t.Errorf("Improperly allowed reroll of %s on second roll", diename[Die0])
	// }
}

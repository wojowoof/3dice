package dicegame

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
	"wojones.com/src/dicescore"
	"wojones.com/src/diceturn"
)

type DiceGame struct {
	ID         string                           `json:"game_id"`
	Players    []string                         `json:"players"`
	Scores     map[string]dicescore.PlayerScore `json:"scores"`
	PrevPlayer *string                          `json:"prev_player"`
	PrevTurn   *diceturn.DiceTurn
	CurPlayer  *string             `json:"cur_player"`
	Turns      []diceturn.DiceTurn `json:"turns"`
}

func NewGame(ID string, p1 string, p2 string, p3 string) DiceGame {
	dg := DiceGame{ID: ID, Players: []string{p1, p2, p3},
		Scores: map[string]dicescore.PlayerScore{p1: {Player: p1},
			p2: {Player: p2, Chevrons: []dicescore.Chevron{{}, {}}},
			p3: {Player: p3}}}
	dg.CurPlayer = &dg.Players[0]
	dg.Turns = []diceturn.DiceTurn{{Player: dg.Players[0], Score: 0, NumRolls: 0}}
	for _, player := range dg.Players {
		dg.Scores[player] = dicescore.PlayerScore{Player: player, Chevrons: []dicescore.Chevron{{Count: 0, Filled: false, Paid: false}}}
	}
	return dg
}

// | || ||| |||| +++++ +++++
func centerin(s string, width int) string {
	return fmt.Sprintf("%[1]*s", -width, fmt.Sprintf("%[1]*s", (width+len(s))/2, s))
}

func asticks(count int) string {
	s := ""
	const twidth = 5
	for i := 0; i < 4; i++ {
		if count <= 0 {
			s += strings.Repeat(" ", twidth)
		} else if count > 5 {
			s += "++++ "
		} else {
			s += fmt.Sprintf("%[1]*s", -twidth, strings.Repeat("i", count))
		}
		count -= 5
	}
	return s
}

func (dg DiceGame) String() string {
	headline := fmt.Sprintf("Game %s: %s", dg.ID, strings.Join(dg.Players, ", "))
	return headline
}

func (dg DiceGame) Scorecard() string {
	const pwidth = 20
	scorecard := ""

	for idx, player := range dg.Players {
		if idx > 0 {
			scorecard += "|"
		}
		scorecard += centerin(player, pwidth)
		// scorecard += fmt.Sprintf("%[1]*s", -pwidth, fmt.Sprintf("%[1]*s", (pwidth+len(p))/2, p))
	}

	scorecard += "\n" + strings.Repeat(strings.Repeat("-", pwidth), len(dg.Players)) + "\n"
	for idx, player := range dg.Players {
		if idx > 0 {
			scorecard += "|"
		}
		if len(dg.Scores[player].Chevrons) > 0 {
			scorecard += centerin(asticks(int(dg.Scores[player].Chevrons[0].Count)), pwidth)
			// score := fmt.Sprintf("%d", dg.Scores[player].Chevrons[0].Count)
			// scorecard += centerin(score, pwidth)
		} else {
			scorecard += strings.Repeat(" ", pwidth)
		}
	}

	return scorecard + "\n"
}

func (dg DiceGame) CurrentTurn() diceturn.DiceTurn {
	return dg.Turns[len(dg.Turns)-1]
}
func (dg DiceGame) CurTurn() string {
	tp := dg.Turns[len(dg.Turns)-1]
	s := fmt.Sprintf("%v", tp)
	if dg.PrevTurn != nil {
		s += fmt.Sprintf("\n\tAgainst %s's %s", dg.PrevTurn.Player,
			dg.PrevTurn.RollString())
	} else {
		s += "\n\tto start the game!!!"
	}

	return s
}

func (dg *DiceGame) RollCheck(dmap int) error {
	return dg.Turns[len(dg.Turns)-1].RollCheck(dmap)
}

func (dg *DiceGame) RollWith(d1 int, d2 int, d3 int) error {
	fmt.Printf("Rolling %d %d %d\n", d1, d2, d3)
	tp := &dg.Turns[len(dg.Turns)-1]
	if tp.NumRolls >= 3 {
		return fmt.Errorf("cannot roll (%d rolls already)", tp.NumRolls)
	}
	fmt.Printf("Before: %v\n", tp)

	toroll := 0
	nroll := 0
	if d1 > 0 {
		toroll |= diceturn.Die0
		nroll++
	}
	if d2 > 0 {
		toroll |= diceturn.Die1
		nroll++
	}
	if d3 > 0 {
		toroll |= diceturn.Die2
		nroll++
	}

	if toroll == 0 {
		return fmt.Errorf("rolling no dice is not a roll")
	}
	fmt.Printf("Asking to roll: %03b\n", toroll)

	var prevroll *diceturn.DiceRoll = nil

	switch tp.NumRolls {
	case 0:
		if toroll != diceturn.AllDice {
			return fmt.Errorf("must roll all dice on the first roll")
		}

	case 1:
		prevroll = &tp.Rolls[0]
		prevroll.Kept = ^toroll & diceturn.AllDice
		fmt.Printf("Kept from roll 1: 0x%03b, rolling 0x%03b, reroll 0x%03b\n",
			tp.Rolls[0].Kept, toroll, tp.Rolls[0].Kept&toroll)

		// FWIW: This cannot happen in this code factoring where the inputs determine which
		// dice are being rolled, and in turn which were kept on the prior roll (5 lines above)
		// BUT: Keep the code for when (if) it gets refactored
		if tp.Rolls[0].Kept&toroll != 0 {
			return fmt.Errorf("cannot re-roll a kept die (0x%03b)", tp.Rolls[0].Kept&toroll)
		}

	case 2:
		prevroll = &tp.Rolls[1]
		prevroll.Kept = ^toroll & diceturn.AllDice
		allkept := prevroll.Kept | tp.Rolls[0].Kept

		fmt.Printf("Kept from roll 2: 0x%03b, rolling 0x%03b, reroll 0x%03b\n",
			tp.Rolls[1].Kept, toroll, tp.Rolls[0].Kept&toroll)
		if nroll > 1 {
			// ASSERT: allkept & toroll is nonzero; this is the same case as below, really
		}
		if allkept&toroll != 0 {
			// The special case: after roll 2, they can pick up their original dice IFF rolling for fives
			// FOR NOW: Just error out if rerolling
			return fmt.Errorf("cannot reroll a kept die")
		}
	}

	// Make sure they're keeping at least one die for roll two
	// AND: if they kept two dice on roll one, they can re-roll the remaining
	// die on turn 3
	// Validate they're not re-rolling a kept die

	// if tp.NumRolls == 2 {
	// Make sure nothing kept is being rerolled
	// TODO: The triple-5 exception (can roll one kept die if two rolled die are both fives)
	// }

	// TODO: Make sure they're only rolling one die
	// Update prior roll's kept value
	//fmt.Printf("Set kept for roll %d to 0x%03b\n", tp.NumRolls-1, ^toroll&diceturn.AllDice)
	//tp.Rolls[tp.NumRolls-1].Kept = ^toroll & diceturn.AllDice
	// }

	tp.Rolls = append(tp.Rolls, diceturn.DiceRoll{Rolled: toroll})
	drp := &tp.Rolls[tp.NumRolls]

	if tp.NumRolls != 0 {
		drp.RollResults = tp.Rolls[tp.NumRolls-1].RollResults
	}
	if toroll&diceturn.Die0 != 0 {
		drp.RollResults[0] = d1
	}
	if toroll&diceturn.Die1 != 0 {
		drp.RollResults[1] = d2
	}
	if toroll&diceturn.Die2 != 0 {
		drp.RollResults[2] = d3
	}

	if drp.IsConsec() {
		drp.Consecs = true
	}

	tp.NumRolls++
	fmt.Printf("After: %v\n", tp)

	// TODO: If tp.NumRolls >= 3 then the turn is over!
	return nil
}

func (dg *DiceGame) PassDice(player string) int {
	idx := slices.IndexFunc(dg.Players, func(s string) bool { return s == player })
	if idx < 0 {
		fmt.Printf("PassDice: no player %s?\n", player)
		return -1
	}

	fmt.Printf("PassDice: passing to %s\n", dg.Players[idx])

	// TODO: Cleanup the last turn, assign score, etc
	dg.Turns[len(dg.Turns)-1].CloseTurn()

	dg.PrevPlayer = dg.CurPlayer
	dg.PrevTurn = &dg.Turns[len(dg.Turns)-1]
	dg.CurPlayer = &dg.Players[idx]

	dg.Turns = append(dg.Turns, diceturn.NewTurn(*dg.CurPlayer))
	return 0
}

func (dg *DiceGame) RollDice() int {
	return 0
}

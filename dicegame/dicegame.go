package dicegame

import (
	"fmt"
	"golang.org/x/exp/slices"
	"strings"
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
		Scores: map[string]dicescore.PlayerScore{p1: dicescore.PlayerScore{Player: p1},
			p2: dicescore.PlayerScore{Player: p2, Chevrons: []dicescore.Chevron{dicescore.Chevron{}, dicescore.Chevron{}}},
			p3: dicescore.PlayerScore{Player: p3}}}
	dg.CurPlayer = &dg.Players[0]
	dg.Turns = []diceturn.DiceTurn{diceturn.DiceTurn{Player: dg.Players[0], Score: 0, NumRolls: 0}}
	for _, player := range dg.Players {
		dg.Scores[player] = dicescore.PlayerScore{Player: player, Chevrons: []dicescore.Chevron{dicescore.Chevron{Count: 0, Filled: false, Paid: false}}}
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

func (dg *DiceGame) RollWith(d1 int, d2 int, d3 int) error {
	fmt.Printf("Rolling %d %d %d\n", d1, d2, d3)
	tp := &dg.Turns[len(dg.Turns)-1]
	fmt.Printf("Before: %v\n", tp)

	toroll := 0
	if d1 > 0 {
		toroll |= diceturn.Die0
	}
	if d2 > 0 {
		toroll |= diceturn.Die1
	}
	if d3 > 0 {
		toroll |= diceturn.Die2
	}

	fmt.Printf("Rolled dice: %03b\n", toroll)

	switch tp.NumRolls {
	case 0:
		if toroll != diceturn.AllDice {
			return fmt.Errorf("Must roll all dice on the first roll")
		}

	case 1:
		fmt.Printf("Kept from roll 1: 0x%03b, rolling 0x%03b, reroll 0x%03b\n",
			tp.Rolls[0].Kept, toroll, tp.Rolls[0].Kept&toroll)
		if tp.Rolls[0].Kept&toroll != 0 {
			return fmt.Errorf("Cannot re-roll a kept die (0x%03b)", tp.Rolls[0].Kept&toroll)
		}
	case 2:
		fmt.Printf("Kept from roll 2: 0x%03b, rolling 0x%03b, reroll 0x%03b\n",
			tp.Rolls[1].Kept, toroll, tp.Rolls[1].Kept&toroll)
		if tp.Rolls[1].Kept&toroll != 0 {
			// The special case: after roll 2,
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
	drp.Kept = ^toroll & diceturn.AllDice

	tp.NumRolls++
	fmt.Printf("After: %v\n", tp)

	// TODO: If tp.NumRolls >= 3 then the turn is over!
	return nil
}

func (dg *DiceGame) PassDice(player string) int {
	idx := slices.IndexFunc(dg.Players, func(s string) bool { return s == player })
	if idx < 0 {
		return -1
	}

	// TODO: Cleanup the last turn, assign score, etc

	dg.PrevPlayer = dg.CurPlayer
	dg.PrevTurn = &dg.Turns[len(dg.Turns)-1]
	dg.CurPlayer = &dg.Players[idx]

	dg.Turns = append(dg.Turns, diceturn.NewTurn(*dg.CurPlayer))
	return 0
}

func (dg *DiceGame) RollDice() int {
	return 0
}

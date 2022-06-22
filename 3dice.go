package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"wojones.com/src/dicegame"
	"wojones.com/src/dicescore"
	//"wojones.com/src/diceturn"
)

func dispatch(dg *dicegame.DiceGame, argv []string) (int, error) {
	if 0 == strings.Compare("help", argv[0]) {
		fmt.Println("Well, help yerself!")
	} else if 0 == strings.Compare("exit", argv[0]) {
		fmt.Println("Latah!")
		return 0, nil
	} else if 0 == strings.Compare("status", argv[0]) {
		fmt.Printf("Game: %v\n", *dg)
		fmt.Printf("Turn: %s\n", dg.CurTurn())
	} else if 0 == strings.Compare("score", argv[0]) {
		fmt.Printf("Scorecard:\n%s", dg.Scorecard())
	} else if 0 == strings.Compare("rollcheck", argv[0]) {
		if len(argv) < 2 {
			return 1, fmt.Errorf("ERROR: Must specify roll bitmap")
		} else if len(argv) > 2 {
			return 1, fmt.Errorf("ERROR: usage: rollcheck <dicebits>")
		}
		if dmap, err := strconv.ParseInt(argv[1], 0, 32); err != nil {
			return 1, fmt.Errorf("ERROR: Invalid dicemap specification: %s", argv[1])
		} else if e := dg.RollCheck(int(dmap)); e != nil {
			fmt.Printf("Cannot roll 0x%03b: %v\n", dmap, e)
		} else {
			fmt.Printf("Sure, roll 0x%03b!!!\n", dmap)
		}
	} else if 0 == strings.Compare("roll", argv[0]) {
		if len(argv) < 2 {
			fmt.Printf("ERROR: must specify dice values")
		}
		dice := []int{0, 0, 0}
		for i := 1; i < len(argv); i++ {
			if dval, err := strconv.Atoi(argv[i]); err != nil {
				fmt.Printf("Error converting %s: %v\n", argv[i], err)
				return 1, fmt.Errorf("Invalid integer value %s", argv[i])
			} else if dval > 6 || dval < 0 {
				fmt.Printf("Invalid dice value \"%s\"\n", argv[i])
				return 1, fmt.Errorf("Invalid value for a die: %s", argv[i])
			} else {
				dice[i-1] = dval
			}
		}

		fmt.Printf("Rolling %d/%d/%d\n", dice[0], dice[1], dice[2])
		if err := dg.RollWith(dice[0], dice[1], dice[2]); err != nil {
			fmt.Printf("Whoops - %v\n", err)
		}
	} else if 0 == strings.Compare("passto", argv[0]) {
		if len(argv) < 2 {
			fmt.Println("ERROR: usage: passto <player>")
		}
		ret := dg.PassDice(argv[1])
		if 0 != ret {
			fmt.Printf("FAIL! (%d)\n", ret)
		}
	} else {
		fmt.Printf("Huh? Wassa '%s' mean?\n", argv[0])
	}
	return 1, nil
}

func interact(dg *dicegame.DiceGame) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Talking shit?")
	fmt.Println("-------------")

	for {
		fmt.Print("3d% ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSuffix(text, "\n")
		argv := strings.Fields(text)

		if goon, err := dispatch(dg, argv); err != nil {
			fmt.Printf("Error with %s: %v\n", argv[0], err)
		} else if goon == 0 {
			break
		}
	}
}

func main() {
	fmt.Printf("Talking shit?\n")

	cv := dicescore.NewChevron()
	fmt.Printf("Chevron: %v\n", cv)

	// tdg := dicegame.DiceGame{ID: "Game1",
	// 	Players: []string{"Alpha", "Beta", "Greg"},
	// }
	tdg := dicegame.NewGame("Game001", "Freddy", "Danny", "Smeck")
	fmt.Printf("Game: %v\n", tdg)

	//tdg.Scores["Freddy"].Chevrons[0].Count = 11
	//tdg.Scores["Danny"].Chevrons[0].Count = 17
	buf, err := json.MarshalIndent(tdg, "", " ")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(buf))
	}

	// Try playing a bit
	// tdg.Turns = append(tdg.Turns, diceturn.DiceTurn{Player: "Freddy",
	// 	DiceVals: [3]int{1, 1, 1},
	// 	Score:    1, NumRolls: 3,
	// 	Rolls: [3]diceturn.DiceRoll{
	// 		{Rolled: 0x07,
	// 			RollResults: [3]int{5, 1, 2},
	// 			OffTable:    0, Kept: diceturn.Die1, Consecs: false,
	// 		}, {Rolled: diceturn.Die0 | diceturn.Die2,
	// 			RollResults: [3]int{1, 1, 2},
	// 			OffTable:    0, Kept: diceturn.Die0 | diceturn.Die1, Consecs: false,
	// 		}, {Rolled: diceturn.Die2,
	// 			RollResults: [3]int{1, 1, 1},
	// 			OffTable:    0, Kept: diceturn.AllDie, Consecs: false,
	// 		}}})

	// Thoughts on the processing:
	//   THings will go in player decision steps, as in:
	//   PassDice(toplayer) - not necessarily the next player in the list!
	//   RollDice() - and display result to roller
	//     ... in here deal with consecutives, off the table re-roll, etc.
	//   ChooseKeeps() or EndTurn()
	//   RollDice()
	//   ChooseKeeps() or EndTurn()

	// TODO: Handle someone rolling for someone else, like a community player
	// stepping in for a roller; they roll the dice, the roller gets points

	// TODO: Another screwy thing: some scorers allow a player to pick up a
	// previously kept die after the second roll, to try and complete a triple
	// 5 or 6. I'm *not* one of those scorers, unless they're rolling two dice
	// to get triple-5. However, make sure to allow for this type of thing.

	// tdg.PassDice("Freddy")
	// fmt.Printf("Turn: %v\n", tdg.Turns[0])

	interact(&tdg)
	// tdg.RollDice()
}

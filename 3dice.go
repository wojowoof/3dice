package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"golang.org/x/exp/slices"
	"os"
	//"sort"
	"strconv"
	"strings"
	"wojones.com/src/dicegame"
	"wojones.com/src/dicescore"
	"wojones.com/src/diceturn"
	//"wojones.com/src/diceturn"
)

type dispfunc func(*dicegame.DiceGame, []string) (int, error)

type cmd struct {
	command string
	disp    func(*dicegame.DiceGame, []string) (int, error)
	argstr  string
	usestr  string
}

var cmdz = []cmd{
	{"help", helpme, "", "get help"},
	{"exit", quitme, "", "quit the command loop"},
	{"status", givestatus, "", "show current game status"},
	{"score", showscore, "", "show the scorecard"},
	{"history", showhist, "", "Display the game history"},
	{"rollcheck", rollcheck, "<rollbits>", "check validity of a roll"},
	{"roll", rolldice, "<d0> <d1> <d2>", "roll with given values (0 is a keep)"},
	{"passto", passto, "<player>", "end turn and pass dice to specified player"},
}

type cmdhelp struct {
	command string
	usage   string
	argstr  string
}

var cmddoc = []cmdhelp{}

func showhist(dg *dicegame.DiceGame, argv []string) (int, error) {
	fmt.Printf("History of game: %s (%d turns)\n", dg.ID, len(dg.Turns))
	for turnno := 0; turnno < len(dg.Turns)-1; turnno++ {
		ct := dg.Turns[turnno]
		fmt.Printf("%s rolled a %d ()\n", ct.Player, ct.Score)
	}
	// Show the current (last) turns
	ct := dg.Turns[len(dg.Turns)-1]
	fmt.Printf("Current: %s\n", ct.String())
	return 1, nil
}

func passto(dg *dicegame.DiceGame, argv []string) (int, error) {
	if len(argv) < 2 {
		return 1, fmt.Errorf("Must specify a player")
	}
	ret := dg.PassDice(argv[1])
	if 0 != ret {
		return 1, fmt.Errorf("FAIL! (%d)\n", ret)
	}
	return 1, nil
}

func rolldice(dg *dicegame.DiceGame, argv []string) (int, error) {
	if len(argv) < 2 {
		return 1, fmt.Errorf("ERROR: must specify dice values")
	}
	dice := []int{0, 0, 0}
	for i := 1; i < len(argv); i++ {
		if argv[i] == "-" {
			dice[i-1] = 0
		} else if dval, err := strconv.Atoi(argv[i]); err != nil {
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
	} else {
		dt := dg.Turns[len(dg.Turns)-1]
		cr := dt.Rolls[dt.NumRolls-1]
		rv, special := cr.TurnValue()

		if cr.IsConsec() {
			fmt.Printf("CONSECUTIVES!")
		}
		if rv >= 0 || diceturn.NothingSpecial != special {
			fmt.Printf("That's a %s\n", cr.TurnValueString())

		} else {
			fmt.Printf("Whoops - value %d?\n", rv)
		}
	}

	// TODO: Advance turn if this is third rolls
	ct := dg.Turns[len(dg.Turns)-1]
	if ct.NumRolls > 2 {
		fmt.Printf("Turn over!\n")
	}

	return 1, nil
}

func rollcheck(dg *dicegame.DiceGame, argv []string) (int, error) {
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
	return 1, nil
}

func showscore(dg *dicegame.DiceGame, argv []string) (int, error) {
	fmt.Printf("Scorecard:\n%s", dg.Scorecard())
	return 1, nil
}

func givestatus(dg *dicegame.DiceGame, argv []string) (int, error) {
	fmt.Printf("Game: %v\n", *dg)
	fmt.Printf("Turn: %s\n", dg.CurTurn())
	return 1, nil
}

func quitme(dg *dicegame.DiceGame, argv []string) (int, error) {
	fmt.Println("Adios!")
	return 0, nil
}

func helpme(dg *dicegame.DiceGame, argv []string) (int, error) {
	fmt.Println("Well, please, help yourself!")
	for _, cp := range cmddoc {
		fmt.Printf("%s ", cp.command)
		if len(cp.argstr) > 0 {
			fmt.Printf("%s ", cp.argstr)
		}
		fmt.Printf("- %s\n", cp.usage)
	}
	return 1, nil
}

func givescore(dg *dicegame.DiceGame, argv []string) (int, error) {
	fmt.Printf("Scorecard:\n%s", dg.Scorecard())
	return 1, nil
}

func dispatch(dg *dicegame.DiceGame, argv []string) (int, error) {

	if c := slices.IndexFunc(cmdz, func(c cmd) bool { return c.command == argv[0] }); c < 0 {
		return 1, fmt.Errorf("Unknown command: \"%v\"", argv[0])
	} else {
		ret, err := cmdz[c].disp(dg, argv)
		return ret, err
	}
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

	for _, c := range cmdz {
		cmddoc = append(cmddoc, cmdhelp{c.command, c.argstr, c.usestr})
	}
	// TODO: would this be better in an init function? Or maybe the whole struct
	// command stuff should be broken into its own library ...
	/*sort.Slice(cmdz, func(i int, j int) bool {
		return cmdz[i].command < cmdz[j].command
	})*/
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

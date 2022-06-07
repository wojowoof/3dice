package dicescore

import "fmt"

type Chevron struct {
	Count  int32 `json:"count"`
	Filled bool  `json:"is_filled"`
	Paid   bool  `json:"is_paid"`
}

type PlayerScore struct {
	Player   string    `json:"player_name"`
	Chevrons []Chevron `json:"chevrons"`
}

func NewChevron() Chevron {
	return Chevron{Count: 0, Filled: false, Paid: false}
}

func (c Chevron) String() string {
	return fmt.Sprintf("C: %d", c.Count)
}

func NewPlayerScore(player string) PlayerScore {
	return PlayerScore{Player: player,
		Chevrons: []Chevron{{Count: 0, Filled: false, Paid: false}}}
}

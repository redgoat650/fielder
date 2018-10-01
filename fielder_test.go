package fielder

import (
	"fmt"
	"testing"
)

func TestScheduleGame(t *testing.T) {

	var testPlayers map[string]PlayerGender = map[string]PlayerGender{
		"Nick":  MaleGender,
		"Patty": FemaleGender,
		"Rob":   MaleGender,
		"Scott": MaleGender,
		// "Shawn":     FemaleGender,
		"Craig": MaleGender,
		"Cody":  MaleGender,
		// "Jenna":     FemaleGender,
		// "Sam": MaleGender,
		// "Kira":      FemaleGender,
		// "Shelli":    FemaleGender,
		// "Yev": MaleGender,
		// "Blair":     FemaleGender,
		"Jen": FemaleGender,
		// "Patrick":   MaleGender,
		"Alexandra": FemaleGender,
		// "Son":       MaleGender,
		"Christina": FemaleGender,
		"Andrew":    MaleGender,
		"David":     MaleGender,
		"Brett":     MaleGender,
	}

	roster := new(Roster)
	roster.players = make([]*Player, 0)

	for player, gender := range testPlayers {
		newPlayer := NewPlayer(player, "blab", gender)
		roster.players = append(roster.players, newPlayer)
	}

	for _, v := range roster.players {
		fmt.Printf("%v\n", v)
	}

	innings := 5
	game := ScheduleGame(innings, roster)

	fmt.Println(game)
}

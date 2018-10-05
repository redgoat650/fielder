package fielder

import (
	"fmt"
	"testing"
)

var testPlayers = map[string]PlayerGender{
	"Nick":      MaleGender,
	"Patty":     FemaleGender,
	"Rob":       MaleGender,
	"Scott":     MaleGender,
	"Shawn":     FemaleGender,
	"Craig":     MaleGender,
	"Cody":      MaleGender,
	"Jenna":     FemaleGender,
	"Sam":       MaleGender,
	"Kira":      FemaleGender,
	"Shelli":    FemaleGender,
	"Yev":       MaleGender,
	"Blair":     FemaleGender,
	"Jen":       FemaleGender,
	"Patrick":   MaleGender,
	"Alexandra": FemaleGender,
	"Son":       MaleGender,
	"Christina": FemaleGender,
	"Andrew":    MaleGender,
	"David":     MaleGender,
	"Brett":     MaleGender,
}

func TestScheduleGame(t *testing.T) {

	roster := NewRoster()

	for player, gender := range testPlayers {
		newPlayer := NewPlayer(player, "blab", gender)
		roster.Players = append(roster.Players, newPlayer)
	}

	for _, v := range roster.Players {
		fmt.Printf("%v\n", v)
	}

	innings := 5
	game := NewGame(innings, 0)
	game.SetRoster(roster)
	schedErr := game.ScheduleGame()
	if schedErr != nil {
		panic(schedErr)
	}

	fmt.Println(game)
}

func TestSeasonSchedule(t *testing.T) {

	numGames := 8
	inningsPerGame := 5
	season := NewSeason(numGames, inningsPerGame)

	team := season.Team
	for plName, plGender := range testPlayers {
		pl := NewPlayer(plName, "blub", plGender)

		for _, gamePtr := range season.Games {
			if plName != "Cody" {
				pl.Attendance[gamePtr] = true
			}
		}

		team.AddPlayer(pl)
	}

	season.ScheduleAllGames()

	fmt.Println(season)
}

func TestPreferences(t *testing.T) {

	roster := NewRoster()

	for player, gender := range testPlayers {
		newPlayer := NewPlayer(player, "blab", gender)
		if newPlayer.FirstName == "Nick" {
			newPlayer.Pref = append(newPlayer.Pref, Third)
		}
		if newPlayer.FirstName == "Cody" {
			newPlayer.Pref = append(newPlayer.Pref, First)
		}
		if newPlayer.FirstName == "Rob" {
			newPlayer.Pref = append(newPlayer.Pref, Catcher)
		}
		if newPlayer.FirstName == "Craig" {
			newPlayer.Pref = append(newPlayer.Pref, Pitcher)
		}
		if newPlayer.FirstName == "Sam" {
			newPlayer.Pref = append(newPlayer.Pref, Pitcher)
		}
		roster.Players = append(roster.Players, newPlayer)
	}

	for _, v := range roster.Players {
		fmt.Printf("%v\n", v)
	}

	innings := 5
	game := NewGame(innings, 0)
	game.SetRoster(roster)
	schedErr := game.ScheduleGame()
	if schedErr != nil {
		panic(schedErr)
	}

	fmt.Println(game)
}

package fielder

import (
	_ "net/http/pprof"
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

// func TestScheduleGame(t *testing.T) {
// 	go func() {
// 		log.Println(http.ListenAndServe("localhost:6060", nil))
// 	}()

// 	roster := NewRoster()

// 	for player, gender := range testPlayers {
// 		newPlayer := NewPlayer(player, "blab", gender)
// 		roster.AddPlayer(newPlayer)
// 	}

// 	for v := range roster.Players {
// 		fmt.Printf("%v\n", v)
// 	}

// 	innings := 5
// 	game := NewGame(innings, 0)
// 	game.SetRoster(roster)
// 	schedErr := game.ScheduleGame()
// 	if schedErr != nil {
// 		panic(schedErr)
// 	}

// 	fmt.Println(game)
// }

// func TestSeasonSchedule(t *testing.T) {

// 	numGames := 8
// 	inningsPerGame := 5
// 	season := NewSeason(numGames, inningsPerGame)

// 	team := season.Team
// 	for plName, plGender := range testPlayers {
// 		pl := NewPlayer(plName, "blub", plGender)

// 		// for _, gamePtr := range season.Games {
// 		// 	if plName != "Cody" {
// 		// 		pl.Attendance[gamePtr] = true
// 		// 	}
// 		// }

// 		team.AddPlayer(pl)
// 	}

// 	season.ScheduleAllGames()

// 	fmt.Println(season)
// }

// func TestPreferences(t *testing.T) {

// 	roster := NewRoster()

// 	for player, gender := range testPlayers {
// 		newPlayer := NewPlayer(player, "blab", gender)
// 		if newPlayer.FirstName == "Nick" {
// 			newPlayer.Pref[Third] = 1.0
// 		}
// 		if newPlayer.FirstName == "Cody" {
// 			newPlayer.Pref[First] = 1.0
// 		}
// 		if newPlayer.FirstName == "Rob" {
// 			newPlayer.Pref[Catcher] = 1.0
// 		}
// 		if newPlayer.FirstName == "Craig" {
// 			newPlayer.Pref[Pitcher] = 1.0
// 		}
// 		if newPlayer.FirstName == "Sam" {
// 			newPlayer.Pref[Pitcher] = 1.0
// 		}
// 		roster.AddPlayer(newPlayer)
// 	}

// 	for v := range roster.Players {
// 		fmt.Printf("%v\n", v)
// 	}

// 	innings := 5
// 	game := NewGame(innings, 0)
// 	game.SetRoster(roster)
// 	schedErr := game.ScheduleGame()
// 	if schedErr != nil {
// 		panic(schedErr)
// 	}

// 	fmt.Println(game)
// }

// func TestPreferencesNewScheduler(t *testing.T) {

// 	roster := NewRoster()

// 	for player, gender := range testPlayers {
// 		newPlayer := NewPlayer(player, "blab", gender)
// 		if newPlayer.FirstName == "Nick" {
// 			newPlayer.Pref[Third] = 1.0
// 			newPlayer.Skill = 1.0
// 		}
// 		if newPlayer.FirstName == "Cody" {
// 			newPlayer.Pref[First] = 1.0
// 			newPlayer.Skill = 1.0
// 		}
// 		if newPlayer.FirstName == "Rob" {
// 			newPlayer.Pref[Catcher] = 1.0
// 			newPlayer.Skill = 1.0
// 		}
// 		if newPlayer.FirstName == "Craig" {
// 			newPlayer.Pref[Pitcher] = 1.0
// 			newPlayer.Skill = 1.0
// 		}
// 		if newPlayer.FirstName == "Sam" {
// 			newPlayer.Pref[Pitcher] = 1.0
// 			newPlayer.Skill = 0.2
// 		}
// 		if newPlayer.FirstName == "Patty" {
// 			newPlayer.Pref[RShort] = 1.0
// 			newPlayer.Pref[LShort] = 1.0
// 			newPlayer.Pref[Pitcher] = -8
// 			newPlayer.Skill = 1.0

// 		}
// 		roster.AddPlayer(newPlayer)
// 	}

// 	for v := range roster.Players {
// 		fmt.Printf("%v\n", v)
// 	}

// 	innings := 5
// 	game := NewGame(innings, 0)
// 	game.SetRoster(roster)
// 	schedErr := game.ScheduleGame2()
// 	if schedErr != nil {
// 		panic(schedErr)
// 	}

// 	fmt.Println(game)
// }

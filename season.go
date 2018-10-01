package fielder

import (
	"fmt"
	"strings"
)

type Season struct {
	Team  *Team
	Games []*Game
}

func NewSeason(numGames int, inningsPerGame int) *Season {
	season := new(Season)

	season.Games = make([]*Game, 0)
	for gameNum := 0; gameNum < numGames; gameNum++ {
		season.Games = append(season.Games, NewGame(inningsPerGame))
	}

	season.Team = NewTeam()

	return season
}

func (season *Season) ScheduleAllGames() error {

	if season.Team == nil {
		panic("No team for this season")
	}

	for _, game := range season.Games {

		gameRoster := NewRoster()

		//Add players to the roster if they're marked as attending
		for _, player := range season.Team.Players {
			if player.IsAttending(game) {
				gameRoster.AddPlayer(player)
			}
		}

		//Schedule this game
		game.SetRoster(gameRoster)

		gameSchedErr := game.ScheduleGame()
		if gameSchedErr != nil {
			return gameSchedErr
		}

	}
	return nil
}

func (season Season) String() string {
	str := new(strings.Builder)
	for gameNum, game := range season.Games {

		str.WriteString(fmt.Sprintf("Game %d:\n", gameNum))

		str.WriteString(game.String())

		str.WriteString("---------------------------------\n")

	}
	return str.String()
}

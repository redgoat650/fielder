package fielder

import (
	"fmt"
	"strings"
)

//Season is a struct for describing the Season for a Team.
type Season struct {
	Team  *Team
	Games []*Game
}

//NewSeason initializes a Season with a new Team for a given
//number of games, with a given number of innings per game, and
//returns its pointer
func NewSeason(numGames int, inningsPerGame int) *Season {
	season := new(Season)

	season.Games = make([]*Game, 0)
	for gameNum := 0; gameNum < numGames; gameNum++ {
		season.Games = append(season.Games, NewGame(inningsPerGame))
	}

	season.Team = NewTeam()

	return season
}

//ScheduleAllGames schedules the all games in the season for the provided
//Season.
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

//String implements the Stringer interface and nicely displays the
//schedule for the season in a human-readable form
func (season Season) String() string {
	str := new(strings.Builder)
	for gameNum, game := range season.Games {

		str.WriteString(fmt.Sprintf("Game %d:\n", gameNum))

		str.WriteString(game.String())

		str.WriteString("---------------------------------\n")

	}
	return str.String()
}

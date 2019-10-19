package fielder

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"strings"
)

//Season is a struct for describing the Season for a Team.
type Season struct {
	Desc string

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
		season.Games = append(season.Games, NewGame(inningsPerGame, gameNum))
	}

	season.Team = NewTeam("")

	return season
}

// AddGame adds a game to the Season
func (season *Season) AddGame(innings int, startTime, oppTeam, gameDetails string) {

	game := NewGame(innings, len(season.Games))
	game.SetStartStr(startTime)
	game.SetOppTeam(oppTeam)
	game.SetGameDetails(gameDetails)
	season.Games = append(season.Games, game)

}

//ScheduleAllGames schedules the all games in the season for the provided
//Season.
func (season *Season) ScheduleAllGames() error {

	if season.Team == nil {
		panic("No team for this season")
	}

	for gameNo, game := range season.Games {

		gameRoster := game.Roster
		gameRoster.Reset()

		//Add players to the roster if they're marked as attending
		for player := range season.Team.Active.Players {
			if player.IsAttending(game) {
				gameRoster.AddPlayer(player)
			}
		}

		// Schedule this game
		// game.SetRoster(gameRoster)

		gameSchedErr := game.ScheduleGame2()
		if gameSchedErr != nil {
			fmt.Println("Game", gameNo, gameSchedErr.Error())
		}

	}
	return nil
}

//String implements the Stringer interface and nicely displays the
//schedule for the season in a human-readable form
func (season Season) String() string {
	str := new(strings.Builder)

	str.WriteString("Season:\n")
	str.WriteString(season.Desc)
	str.WriteString("\n")

	for gameNum, game := range season.Games {

		str.WriteString(fmt.Sprintf("Game %d:\n", gameNum))

		str.WriteString(game.String())

		str.WriteString("---------------------------------\n")

	}
	return str.String()
}

// SaveToFile saves a Season to the given file name
func (season *Season) SaveToFile(filename string) error {

	for i := range season.Games {
		season.Games[i].Self = nil
	}

	buf := new(bytes.Buffer)

	encErr := gob.NewEncoder(buf).Encode(season)
	if encErr != nil {
		return encErr
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	i := 0
	n := 0
	for {
		n, err = file.Write(buf.Bytes()[i:])
		if err != nil {
			return err
		}
		i += n

		if i >= len(buf.Bytes()) {
			break
		}
	}

	return nil

}

// LoadSeasonFromFile loads a Season from data in a given file name
func LoadSeasonFromFile(filename string) (season *Season, err error) {

	file, err := os.Open(filename)
	if err != nil {
		return
	}

	loadSeason := &Season{}

	err = gob.NewDecoder(file).Decode(loadSeason)

	for i, v := range loadSeason.Games {
		loadSeason.Games[i].Self = v
	}

	return loadSeason, err
}

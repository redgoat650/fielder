package fielder

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"strings"
)

//Team is a structure containing the information on a Team.
//The players on a Team are a superset of each Game's Roster.
type Team struct {
	Name       string
	SeasonList map[string]Season
}

//NewTeam will initialize a new Team and return its pointer
func NewTeam(name string) *Team {
	return &Team{
		Name:       name,
		SeasonList: make(map[string]Season),
	}
}

func (team *Team) String() string {
	str := new(strings.Builder)
	str.WriteString(fmt.Sprintf("Team %s\n", team.Name))

	//Print seasons
	if len(team.SeasonList) > 0 {
		str.WriteString("---------------\n")
		str.WriteString("Seasons:\n")
		str.WriteString("---------------\n")

		for name, season := range team.SeasonList {
			str.WriteString(fmt.Sprintf("%s:\n%s\n", name, season))
		}
	}

	return str.String()
}

//SetTeamName will set the name of the team
func (team *Team) SetTeamName(name string) {
	team.Name = name
}

//AddPlayer will append a new Player to the Team's player list
// func (team *Team) AddPlayer(player *Player) {
// 	team.Players.AddPlayer(player)
// 	team.Active.AddPlayer(player)
// }

// SaveTeamToFile saves a team to a given file name
func SaveTeamToFile(team *Team, filename string) error {

	// for _, season := range team.SeasonList {
	// 	for i, _ := range season.Games {
	// 		season.Games[i].Self = nil
	// 	}
	// }

	// d, err := yaml.Marshal(team)
	// if err != nil {
	// 	return err
	// }

	// path := storage.GetTeamDirectory(team.TeamName)

	buf := new(bytes.Buffer)

	encErr := gob.NewEncoder(buf).Encode(team)
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

// LoadTeamFromFile loads Team data from a given file name
func LoadTeamFromFile(filename string) (team *Team, err error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	loadTeam := &Team{}

	err = gob.NewDecoder(file).Decode(loadTeam)

	// for _, season := range loadTeam.SeasonList {
	// 	for i, v := range season.Games {
	// 		season.Games[i].Self = v
	// 	}
	// }

	return loadTeam, err
}

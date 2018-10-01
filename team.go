package fielder

type Team struct {
	Players  []*Player
	schedule []*Game
}

func NewTeam() *Team {
	team := new(Team)
	team.Players = make([]*Player, 0)
	team.schedule = make([]*Game, 0)
	return team
}

func (team *Team) AddPlayer(player *Player) {
	team.Players = append(team.Players, player)
}

// func (team *Team) BuildRosters() []*Roster {

// 	for _, game := range team.schedule {

// 	}

// }

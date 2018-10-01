package fielder

//Team is a structure containing the information on a Team.
//The players on a Team are a superset of each Game's Roster.
type Team struct {
	Players  []*Player
	schedule []*Game
}

//NewTeam will initialize a new Team and return its pointer
func NewTeam() *Team {
	team := new(Team)
	team.Players = make([]*Player, 0)
	team.schedule = make([]*Game, 0)
	return team
}

//AddPlayer will append a new Player to the Team's player list
func (team *Team) AddPlayer(player *Player) {
	team.Players = append(team.Players, player)
}

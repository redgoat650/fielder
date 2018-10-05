package fielder

//Team is a structure containing the information on a Team.
//The players on a Team are a superset of each Game's Roster.
type Team struct {
	TeamName string
	Players  *Roster
	Active   *Roster

	SeasonList []*Season
}

//NewTeam will initialize a new Team and return its pointer
func NewTeam() *Team {
	team := new(Team)
	team.Players = NewRoster()
	team.Active = NewRoster()
	team.SeasonList = make([]*Season, 0)
	return team
}

//SetTeamName will set the name of the team
func (team *Team) SetTeamName(name string) {
	team.TeamName = name
}

//AddPlayer will append a new Player to the Team's player list
func (team *Team) AddPlayer(player *Player) {
	team.Players.AddPlayer(player)
}

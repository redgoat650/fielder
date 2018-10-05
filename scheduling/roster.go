package fielder

//Roster is a list of Players that will be participating
//in a game
type Roster struct {
	Players []*Player
}

//NewRoster initializes a Roster and returns its pointer
func NewRoster() *Roster {
	roster := new(Roster)
	roster.Players = make([]*Player, 0)
	return roster
}

//AddPlayer is a roster method that will add a Player to the Roster
func (roster *Roster) AddPlayer(player *Player) {
	roster.Players = append(roster.Players, player)
}

//NumPlayers is a roster method that returns the number of players
//currently in the Roster
func (roster *Roster) NumPlayers() int {
	return len(roster.Players)
}

//CountGenders is a Roster method that returns the number of
//females and males in the Roster
func (roster *Roster) CountGenders() (female, male int) {

	for _, playerInfo := range roster.Players {
		if playerInfo.IsFemale() {
			female++
		}
	}

	male = len(roster.Players) - female

	return

}

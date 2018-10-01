package fielder

//Roster is a list of Players that will be participating
//in a game
type Roster struct {
	players []*Player
}

//NewRoster initializes a Roster and returns its pointer
func NewRoster() *Roster {
	roster := new(Roster)
	roster.players = make([]*Player, 0)
	return roster
}

//AddPlayer is a roster method that will add a Player to the Roster
func (roster *Roster) AddPlayer(player *Player) {
	roster.players = append(roster.players, player)
}

//NumPlayers is a roster method that returns the number of players
//currently in the Roster
func (roster *Roster) NumPlayers() int {
	return len(roster.players)
}

//CountGenders is a Roster method that returns the number of
//females and males in the Roster
func (roster *Roster) CountGenders() (female, male int) {

	for _, playerInfo := range roster.players {
		if playerInfo.IsFemale() {
			female++
		}
	}

	male = len(roster.players) - female

	return

}

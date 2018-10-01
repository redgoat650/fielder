package fielder

type Roster struct {
	players []*Player
}

func NewRoster() *Roster {
	roster := new(Roster)
	roster.players = make([]*Player, 0)
	return roster
}

func (roster *Roster) AddPlayer(player *Player) {
	roster.players = append(roster.players, player)
}

func (roster *Roster) NumPlayers() int {
	return len(roster.players)
}

func (roster *Roster) CountGenders() (female, male int) {

	for _, playerInfo := range roster.players {
		if playerInfo.IsFemale() {
			female++
		}
	}

	male = len(roster.players) - female

	return

}

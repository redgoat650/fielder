package fielder

type Roster struct {
	players []*Player
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

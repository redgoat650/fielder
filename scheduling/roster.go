package fielder

import (
	"fmt"
	"strings"
)

//Roster is a list of Players that will be participating
//in a game
type Roster struct {
	Players map[*Player]struct{}
}

//NewRoster initializes a Roster and returns its pointer
func NewRoster() *Roster {
	roster := new(Roster)
	roster.Players = make(map[*Player]struct{}, 0)
	return roster
}

//AddPlayer is a roster method that will add a Player to the Roster
func (roster *Roster) AddPlayer(player *Player) {
	roster.Players[player] = struct{}{}
}

//DropPlayer is a roster method that will drop a Player from the Roster
func (roster *Roster) DropPlayer(player *Player) {
	delete(roster.Players, player)
}

// Reset rests a roster
func (roster *Roster) Reset() {
	for player := range roster.Players {
		roster.DropPlayer(player)
	}
}

//NumPlayers is a roster method that returns the number of players
//currently in the Roster
func (roster *Roster) NumPlayers() int {
	return len(roster.Players)
}

//CountGenders is a Roster method that returns the number of
//females and males in the Roster
func (roster *Roster) CountGenders() (female, male int) {

	for playerInfo := range roster.Players {
		if playerInfo.IsFemale() {
			female++
		}
	}

	male = len(roster.Players) - female

	return

}

// String meets the stringer interface
func (roster Roster) String() string {
	str := new(strings.Builder)
	for player := range roster.Players {
		str.WriteString(fmt.Sprintf("%s", player))
	}
	return str.String()
}

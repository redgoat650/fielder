package fielder

import (
	"fmt"
	"strings"
)

type Game struct {
	Innings []*Inning
	Players *Roster
}

func NewGame(innings int, roster *Roster) *Game {
	game := new(Game)

	game.Innings = make([]*Inning, 0)

	for inningNum := 0; inningNum < innings; inningNum++ {
		game.Innings = append(game.Innings, NewInning())
	}

	game.Players = roster

	return game
}

func (game Game) String() string {

	//Analysis of players in each position by inning
	s := new(strings.Builder)
	for inningNum, inning := range game.Innings {

		s.WriteString(fmt.Sprintf("Inning %d:\n", inningNum))

		for pos, player := range inning.FieldPositions {
			s.WriteString(fmt.Sprintf("%s: %s (%s)\n", pos, player.FirstName, player.Gender))
		}

		s.WriteString("----------------\n")

	}

	//Analysis for each player
	mostInnings := 0
	leastInnings := len(game.Innings)
	mostInningsMale := 0
	mostInningsFemale := 0
	leastInningsMale := len(game.Innings)
	leastInningsFemale := len(game.Innings)

	for _, player := range game.Players.players {

		inningsThisPlayer := 0

		for inningNum, role := range player.Roles {
			s.WriteString(fmt.Sprintf("Inning %d: %s plays ", inningNum, player.FirstName))
			if role == Bench {
				s.WriteString(fmt.Sprintf("(%s)\n", role))
			} else {
				s.WriteString(fmt.Sprintf("%s\n", role))
			}

			if role != Bench {
				inningsThisPlayer++
			}
		}

		s.WriteString(fmt.Sprintf("%s is playing %d innings\n", player.FirstName, inningsThisPlayer))
		s.WriteString(fmt.Sprintf("----------\n"))

		if inningsThisPlayer > mostInnings {
			mostInnings = inningsThisPlayer
		}
		if inningsThisPlayer < leastInnings {
			leastInnings = inningsThisPlayer
		}
		if player.IsFemale() {
			if inningsThisPlayer > mostInningsFemale {
				mostInningsFemale = inningsThisPlayer
			}
			if inningsThisPlayer < leastInningsFemale {
				leastInningsFemale = inningsThisPlayer
			}
		} else {
			if inningsThisPlayer > mostInningsMale {
				mostInningsMale = inningsThisPlayer
			}
			if inningsThisPlayer < leastInningsMale {
				leastInningsMale = inningsThisPlayer
			}

		}

	}

	s.WriteString(fmt.Sprintf("Most innings played by a player: %d\nLeast innings played by a player: %d\n", mostInnings, leastInnings))
	s.WriteString(fmt.Sprintf("Most innings played by a FEMALE: %d\nLeast innings played by a FEMALE: %d\n", mostInningsFemale, leastInningsFemale))
	s.WriteString(fmt.Sprintf("Most innings played by a MALE: %d\nLeast innings played by a MALE: %d\n", mostInningsMale, leastInningsMale))

	s.WriteString("\n")

	//Analysis of each position by inning

	return s.String()
}

func (game *Game) NumPlayers() int {
	return game.Players.NumPlayers()
}

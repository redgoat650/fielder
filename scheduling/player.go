package fielder

import (
	"fmt"
	"strings"
)

//Player is a struct for storing the information and scheduling
//information for each player
type Player struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Gender    PlayerGender

	Pref      []Position
	Seniority float64
	Skill     float64

	Roles         map[*Game][]Position
	scoreByInning []float64

	Attendance map[*Game]bool
}

//NewPlayer initializes a new Player and returns its pointer
func NewPlayer(first, last string, gender PlayerGender) *Player {
	p := new(Player)
	p.FirstName = first
	p.LastName = last
	p.Gender = gender

	p.Pref = make([]Position, 0)
	// p.Roles = make([]Position, 0)
	p.Roles = make(map[*Game][]Position)
	p.scoreByInning = make([]float64, 0)
	p.Attendance = make(map[*Game]bool)

	return p
}

//IsFemale is a helper method for Player that returns whether
//the player is female
func (player *Player) IsFemale() bool {
	return player.Gender == FemaleGender
}

//IsAttending is a helper method for Player that returns whether
//the player in question is planning to attend the provided Game
func (player *Player) IsAttending(game *Game) bool {
	return player.Attendance[game]
}

//PlayerGender is a type that describes the gender of a player
type PlayerGender int

//String is a helper method on PlayerGender that satisfies the
//Stringer interface to assist with printing the gender of the player
func (gender PlayerGender) String() string {
	switch gender {
	case FemaleGender:
		return fmt.Sprintf("Female")
	case MaleGender:
		return fmt.Sprintf("Male")
	default:
		return fmt.Sprintf("Undefined gender")
	}
}

//List of PlayerGenders
const (
	FemaleGender PlayerGender = iota
	MaleGender
	NumGenders int = iota
)

func (player Player) String() string {
	str := new(strings.Builder)

	str.WriteString(fmt.Sprintf("%s %s, %s\n", player.FirstName, player.LastName, player.Gender))

	return str.String()

}

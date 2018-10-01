package fielder

import "fmt"

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

func NewPlayer(first, last string, gender PlayerGender) *Player {
	p := new(Player)
	p.FirstName = first
	p.LastName = last
	p.Gender = gender

	p.Pref = make([]Position, 0)
	// p.Roles = make([]Position, 0)
	p.Roles = make(map[*Game][]Position)
	p.scoreByInning = make([]float64, 0)

	return p
}

func (player *Player) IsFemale() bool {
	return player.Gender == FemaleGender
}

func (player *Player) IsAttending(game *Game) bool {
	return player.Attendance[game]
}

type PlayerGender int

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

const (
	FemaleGender PlayerGender = iota
	MaleGender
	NumGenders int = iota
)

package main

import "fmt"

type Player struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Gender    PlayerGender

	Pref      []Position
	Seniority int
	Skill     int

	Roles []Position
}

func (player *Player) IsFemale() bool {
	return player.Gender == FemaleGender
}

type PlayerGender int

const (
	FemaleGender PlayerGender = iota
	MaleGender
)

type Position int

const (
	Bench Position = iota
	Pitcher
	Catcher
	First
	Second
	Third
	LShort
	RShort
	LField
	LCenter
	RCenter
	RField
	NumFieldPositions int = iota - 1 //Don't include bench
)

type Roster []*Player

func (roster *Roster) NumPlayers() int {
	return len(*roster)
}

func NewPlayer(first, last string) *Player {
	p := new(Player)
	p.FirstName = first
	p.LastName = last

	p.Pref = make([]Position, 0)
	p.Roles = make([]Position, 0)
	return p
}

func main() {

	numPlayersIn := 20

	roster := Roster(make([]*Player, 0))

	for player := 0; player < numPlayersIn; player++ {
		newPlayer := NewPlayer(fmt.Sprintf("Player%d", player), "blab")
		roster = append(roster, newPlayer)
	}

	for _, v := range roster {
		fmt.Printf("%v", v)
	}

	innings := 1
	game := ScheduleGame(innings, &roster)

	fmt.Println(game)
}

type Inning map[Position]*Player

func NewInning() *Inning {
	inning := Inning(make(map[Position]*Player))
	inning[Pitcher] = nil
	inning[Catcher] = nil
	inning[First] = nil
	inning[Second] = nil
	inning[Third] = nil
	inning[LShort] = nil
	inning[RShort] = nil
	inning[LField] = nil
	inning[LCenter] = nil
	inning[RCenter] = nil
	inning[RField] = nil
	return &inning
}

func (inning *Inning) CountPlayersOnField() (filled, unfilled int) {
	filledCount := 0
	unfilledCount := 0
	for _, position := range *inning {
		if position != nil {
			filledCount++
		} else {
			unfilledCount++
		}
	}

	if filledCount > NumFieldPositions {
		panic("How did we fill more positions than exist?")
	}

	return filledCount, unfilledCount
}

func (inning *Inning) CountGenders() (female, male int) {
	for _, position := range *inning {
		if position != nil {
			if position.IsFemale() {
				female++
			} else {
				male++
			}
		}
	}
	return
}

type Game struct {
	Innings []*Inning
	Players *Roster
}

const (
	MinGenderCount = 4
)

func (game *Game) NumPlayers() int {
	return game.Players.NumPlayers()
}

func (game *Game) verify() error {

	for inningNum, inning := range game.Innings {

		filledPositions, _ := inning.CountPlayersOnField()
		if filledPositions < NumFieldPositions {
			if filledPositions != game.NumPlayers() {
				return fmt.Errorf("Not all positions were filled")
			}
		}

		femaleCount, maleCount := inning.CountGenders()
		if femaleCount < MinGenderCount || maleCount < MinGenderCount {
			return fmt.Errorf("Invalid gender assignment in this inning. Inning %d females: %d/%d males %d/%d", inningNum, femaleCount, MinGenderCount, maleCount, MinGenderCount)
		}

	}

	return nil
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

func ScheduleGame(innings int, roster *Roster) *Game {

	game := NewGame(innings, roster)

	for {

		for _, inning := range game.Innings {

			for position := range *inning {
				fmt.Println("scheduling", position)
			}

		}

	}

	err := game.verify()
	if err != nil {
		panic(err)
	}

	return game
}

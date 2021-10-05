package fielder

//Inning is a struct for tracking the players assigned to each
//field position for this game inning.
type Inning struct {
	FieldPositions map[Position]*Player

	// mtx *ScoringMatrix
}

//NewInning will initialize an Inning and return its pointer
func NewInning() *Inning {
	inning := new(Inning)
	inning.FieldPositions = make(map[Position]*Player)
	// inning.InitializeFieldPositions()
	return inning
}

var fieldPosList = []Position{
	Pitcher,
	Catcher,
	First,
	Second,
	Third,
	LShort,
	RShort,
	LField,
	LCenter,
	RCenter,
	RField,
}

func (inning *Inning) FindPlayerPosition(pl *Player) Position {
	pos := Bench
	for pos, checkPlayer := range inning.FieldPositions {
		if checkPlayer.Name == pl.Name {
			return pos
		}
	}
	return pos
}

//InitializeFieldPositions will initialize the map of valid field positions
//to nil values.
// func (inning *Inning) InitializeFieldPositions() {
// 	inning.FieldPositions[Pitcher] = nil
// 	inning.FieldPositions[Catcher] = nil
// 	inning.FieldPositions[First] = nil
// 	inning.FieldPositions[Second] = nil
// 	inning.FieldPositions[Third] = nil
// 	inning.FieldPositions[LShort] = nil
// 	inning.FieldPositions[RShort] = nil
// 	inning.FieldPositions[LField] = nil
// 	inning.FieldPositions[LCenter] = nil
// 	inning.FieldPositions[RCenter] = nil
// 	inning.FieldPositions[RField] = nil
// }

// DropFieldPositions drops the field positions from the inning
func (inning *Inning) DropFieldPositions() {
	for i := range inning.FieldPositions {
		delete(inning.FieldPositions, i)
	}
}

//CountPlayersOnField is an Inning method that will return the number
//of positions that are filled, and the number of positions that are unfilled.
func (inning *Inning) CountPlayersOnField() (filled, unfilled int) {
	filledCount := len(inning.FieldPositions)
	unfilledCount := len(fieldPosList) - filledCount
	// for _, position := range inning.FieldPositions {
	// 	if position != nil {
	// 		filledCount++
	// 	} else {
	// 		unfilledCount++
	// 	}
	// }

	if filledCount > NumFieldPositions {
		panic("How did we fill more positions than exist?")
	}

	return filledCount, unfilledCount
}

//CountGenders is an Inning method that returns the number of
//male and female players assigned to positions in this inning.
func (inning *Inning) CountGenders(roster *Roster) (female, male int) {
	for _, player := range inning.FieldPositions {
		if player.IsFemale() {
			female++
		} else {
			male++
		}
	}
	return
}

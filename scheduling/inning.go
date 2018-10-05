package fielder

//Inning is a struct for tracking the players assigned to each
//field position for this game inning.
type Inning struct {
	FieldPositions map[Position]*Player

	mtx *ScoringMatrix
}

//NewInning will initialize an Inning and return its pointer
func NewInning() *Inning {
	inning := new(Inning)
	inning.FieldPositions = make(map[Position]*Player)
	inning.InitializeFieldPositions()
	return inning
}

//InitializeFieldPositions will initialize the map of valid field positions
//to nil values.
func (inning *Inning) InitializeFieldPositions() {
	inning.FieldPositions[Pitcher] = nil
	inning.FieldPositions[Catcher] = nil
	inning.FieldPositions[First] = nil
	inning.FieldPositions[Second] = nil
	inning.FieldPositions[Third] = nil
	inning.FieldPositions[LShort] = nil
	inning.FieldPositions[RShort] = nil
	inning.FieldPositions[LField] = nil
	inning.FieldPositions[LCenter] = nil
	inning.FieldPositions[RCenter] = nil
	inning.FieldPositions[RField] = nil
}

//CountPlayersOnField is an Inning method that will return the number
//of positions that are filled, and the number of positions that are unfilled.
func (inning *Inning) CountPlayersOnField() (filled, unfilled int) {
	filledCount := 0
	unfilledCount := 0
	for _, position := range inning.FieldPositions {
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

//CountGenders is an Inning method that returns the number of
//male and female players assigned to positions in this inning.
func (inning *Inning) CountGenders() (female, male int) {
	for _, position := range inning.FieldPositions {
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

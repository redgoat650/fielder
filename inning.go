package fielder

type Inning struct {
	FieldPositions map[Position]*Player

	mtx *ScoringMatrix
}

func NewInning() *Inning {
	inning := new(Inning)
	inning.FieldPositions = make(map[Position]*Player)
	inning.InitializeFieldPositions()
	return inning
}

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

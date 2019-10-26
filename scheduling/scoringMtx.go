package fielder

//ScoringMatrix is a helper struct that contains a 2D matrix
//of position-propensity score for a given position by a given player.
type ScoringMatrix struct {
	PlayerInfoMap map[*Player](map[Position]float64)
}

//NewScoringMatrix initializes a ScoringMatrix and returns its pointer
func NewScoringMatrix(roster *Roster) *ScoringMatrix {
	mtx := new(ScoringMatrix)

	// mtx.PlayerInfoMap = make(map[*Player](map[Position]float64))
	// for player := range roster.Players {
	// 	mtx.PlayerInfoMap[player] = make(map[Position]float64)

	// 	//Initialize the scores to zero
	// 	for posIdx := 0; posIdx < NumFieldPositions; posIdx++ {
	// 		pos := posIdx2Position(posIdx)
	// 		mtx.PlayerInfoMap[player][pos] = 0.0
	// 	}

	// }

	return mtx

}

//copy is a helper method that will copy the score values of another
//ScoringMatrix into the ScoringMatrix receiver
func (origMtx *ScoringMatrix) copy() (newMtx *ScoringMatrix) {

	newMtx = new(ScoringMatrix)

	newMtx.PlayerInfoMap = make(map[*Player](map[Position]float64))

	for playerInfo := range origMtx.PlayerInfoMap {
		newMtx.PlayerInfoMap[playerInfo] = make(map[Position]float64)

		for fieldPos, score := range origMtx.PlayerInfoMap[playerInfo] {
			newMtx.PlayerInfoMap[playerInfo][fieldPos] = score
		}
	}

	return
}

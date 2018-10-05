package fielder

//ScoringMatrix is a helper struct that contains a 2D matrix
//of position-propensity score for a given position by a given player.
type ScoringMatrix struct {
	PlayerIdxByPosition [][]float64
}

//NewScoringMatrix initializes a ScoringMatrix and returns its pointer
func NewScoringMatrix(numPlayers int) *ScoringMatrix {
	mtx := new(ScoringMatrix)

	mtx.PlayerIdxByPosition = make([][]float64, numPlayers)
	for playerIdx := range mtx.PlayerIdxByPosition {
		mtx.PlayerIdxByPosition[playerIdx] = make([]float64, NumFieldPositions)
	}

	return mtx

}

//copy is a helper method that will copy the score values of another
//ScoringMatrix into the ScoringMatrix receiver
func (origMtx *ScoringMatrix) copy() (newMtx *ScoringMatrix) {

	newMtx = NewScoringMatrix(len(origMtx.PlayerIdxByPosition))

	for playerIdx := range origMtx.PlayerIdxByPosition {
		for fieldPos, score := range origMtx.PlayerIdxByPosition[playerIdx] {
			newMtx.PlayerIdxByPosition[playerIdx][fieldPos] = score
		}
	}

	return
}

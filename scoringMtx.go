package fielder

type ScoringMatrix struct {
	PlayerIdxByPosition [][]float64
}

func NewScoringMatrix(numPlayers int) *ScoringMatrix {
	mtx := new(ScoringMatrix)

	mtx.PlayerIdxByPosition = make([][]float64, numPlayers)
	for playerIdx := range mtx.PlayerIdxByPosition {
		mtx.PlayerIdxByPosition[playerIdx] = make([]float64, NumFieldPositions)
	}

	return mtx

}

func (origMtx *ScoringMatrix) copy() (newMtx *ScoringMatrix) {

	newMtx = NewScoringMatrix(len(origMtx.PlayerIdxByPosition))

	for playerIdx := range origMtx.PlayerIdxByPosition {
		for fieldPos, score := range origMtx.PlayerIdxByPosition[playerIdx] {
			newMtx.PlayerIdxByPosition[playerIdx][fieldPos] = score
		}
	}

	return
}

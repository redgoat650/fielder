package main

import (
	"fmt"
	"math/rand"
)

type Player struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Gender    PlayerGender

	Pref      []Position
	Seniority float64
	Skill     float64

	Roles         []Position
	scoreByInning []float64
}

func (player *Player) IsFemale() bool {
	return player.Gender == FemaleGender
}

type PlayerGender int

func (gender *PlayerGender) String() string {
	switch *gender {
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

func (pos *Position) String() string {
	switch *pos {
	case Pitcher:
		return "Pitcher"
	case Catcher:
		return "Catcher"
	case First:
		return "First"
	case Second:
		return "Second"
	case Third:
		return "Third"
	case LShort:
		return "LShort"
	case RShort:
		return "RShort"
	case LField:
		return "LField"
	case LCenter:
		return "LCenter"
	case RCenter:
		return "RCenter"
	case RField:
		return "RField"
	default:
		return "Not a position"
	}

}

func posIdx2Position(posIdx int) Position {

	if posIdx > NumFieldPositions {
		panic(fmt.Sprintf("Invalid position index: %d", posIdx))
	}

	//Skip the bench position, return the next position
	return Position(posIdx + 1)
}
func position2PosIdx(pos Position) int {
	return int(pos) - 1
}

type Roster struct {
	players []*Player
}

func (roster *Roster) NumPlayers() int {
	return len(roster.players)
}

func (roster *Roster) CountGenders() (female, male int) {

	for _, playerInfo := range roster.players {
		if playerInfo.IsFemale() {
			female++
		}
	}

	male = len(roster.players) - female

	return

}

func NewPlayer(first, last string, gender PlayerGender) *Player {
	p := new(Player)
	p.FirstName = first
	p.LastName = last
	p.Gender = gender

	p.Pref = make([]Position, 0)
	p.Roles = make([]Position, 0)
	p.scoreByInning = make([]float64, 0)

	return p
}

func main() {

	numPlayersIn := 20

	roster := new(Roster)
	roster.players = make([]*Player, 0)

	for player := 0; player < numPlayersIn; player++ {

		gender := PlayerGender(rand.Intn(NumGenders))

		newPlayer := NewPlayer(fmt.Sprintf("Player%d", player), "blab", gender)
		roster.players = append(roster.players, newPlayer)
	}

	for _, v := range roster.players {
		fmt.Printf("%v\n", v)
	}

	innings := 2
	game := ScheduleGame(innings, roster)

	fmt.Println(game)
}

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

type Inning struct {
	FieldPositions map[Position]*Player

	mtx *ScoringMatrix
}

func NewInning() *Inning {
	inning := new(Inning)
	inning.FieldPositions = make(map[Position]*Player)
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
	return inning
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

const (
	prefScaleFactor = float64(0.3)
	threshDelta     = float64(0.1)
	benchCredit     = float64(1.0)
	retryThreshold  = 2500
)

func checkRoster(roster *Roster) error {
	female, male := roster.CountGenders()
	if female < MinGenderCount {
		return fmt.Errorf("Not enough females. Must forfeit")
	}
	if male < MinGenderCount {
		return fmt.Errorf("Not enough males. Must forfeit.")
	}

	return nil
}

func ScheduleGame(innings int, roster *Roster) *Game {

	err := checkRoster(roster)
	if err != nil {
		panic(err)
	}

	var game *Game
	tries := 0

	for {

		for _, playerInfo := range roster.players {
			playerInfo.Roles = make([]Position, innings)
		}

		game = NewGame(innings, roster)

		for inningNum, inning := range game.Innings {

			initialMax := float64(0)

			if inningNum == 0 {

				inning.mtx = NewScoringMatrix(roster.NumPlayers())

				scoringMtx := inning.mtx.PlayerIdxByPosition
				for playerIdx, playerScores := range scoringMtx {

					for posIdx := range playerScores {

						//Initialize by seniority and skill
						playerInfo := roster.players[playerIdx]
						scoringMtx[playerIdx][posIdx] += playerInfo.Seniority
						scoringMtx[playerIdx][posIdx] += playerInfo.Skill

						//Initialize by player preference
						thisPos := posIdx2Position(posIdx)
						for prefRank, pref := range playerInfo.Pref {
							if pref == thisPos {
								scoringMtx[playerIdx][posIdx] += float64(prefRank) * prefScaleFactor
							}
						}

						if scoringMtx[playerIdx][posIdx] > initialMax {
							initialMax = scoringMtx[playerIdx][posIdx]
						}

					}

				}

			} else {

				//Start with last inning's matrix
				inning.mtx = game.Innings[inningNum-1].mtx.copy()

				for playerIdx, playerInfo := range roster.players {
					if playerInfo.Roles[inningNum-1] == Bench {

						//Player did not play last inning
						for posIdx := range inning.mtx.PlayerIdxByPosition[playerIdx] {
							inning.mtx.PlayerIdxByPosition[playerIdx][posIdx] += benchCredit

							//Set the initial max so we can start from an appropriate threshold
							if inning.mtx.PlayerIdxByPosition[playerIdx][posIdx] > initialMax {
								initialMax = inning.mtx.PlayerIdxByPosition[playerIdx][posIdx]
							}
						}
					}
				}

			}

			//We're going to establish a threshold and continue lowering it
			//and assigning players to positions until all positions are filled.
			assignedPlayers := make(map[*Player]bool)

			for {

				//See if we've filled all the positions and can break out
				filledPositionCount, _ := inning.CountPlayersOnField()
				if filledPositionCount == NumFieldPositions || filledPositionCount == roster.NumPlayers() {
					//Filled all positions needed
					break
				}

				//Set up a list of player candidates
				listOfPlayerIdxsAboveThresholdByPosition := make([][]int, NumFieldPositions)
				for pos := range listOfPlayerIdxsAboveThresholdByPosition {
					listOfPlayerIdxsAboveThresholdByPosition[pos] = make([]int, 0)
				}

				for playerIdx, playerScores := range inning.mtx.PlayerIdxByPosition {

					for posIdx, score := range playerScores {

						if score >= initialMax {

							playerInfo := roster.players[playerIdx]
							if assignedPlayers[playerInfo] {
								//Player has already been assigned to another position
								continue
							}

							listOfPlayerIdxsAboveThresholdByPosition[posIdx] = append(listOfPlayerIdxsAboveThresholdByPosition[posIdx], playerIdx)

							// fmt.Printf("Player %s is a candidate for position %v because score %f >= threshold %f\n", playerInfo.FirstName, posIdx2Position(posIdx), score, initialMax)
						}

					}
				}

				//Each position
				for posIdx, playerIdxList := range listOfPlayerIdxsAboveThresholdByPosition {

					//Check all player indexes for this position
					for {

						if len(playerIdxList) == 0 {
							//Couldn't find a suitable candidate at this threshold level
							break
						}

						position := posIdx2Position(posIdx)
						if inning.FieldPositions[position] != nil {
							//We've already picked a player for this position
							break
						}

						//Pick a player from the list of players that are above the threshold for this position
						pickedListIdx := rand.Intn(len(playerIdxList))
						pickedPlayerIdx := playerIdxList[pickedListIdx]

						pickedPlayerInfo := roster.players[pickedPlayerIdx]

						copy(listOfPlayerIdxsAboveThresholdByPosition[posIdx][:pickedListIdx], listOfPlayerIdxsAboveThresholdByPosition[posIdx][:pickedListIdx+1])
						listOfPlayerIdxsAboveThresholdByPosition[posIdx] = listOfPlayerIdxsAboveThresholdByPosition[posIdx][:len(listOfPlayerIdxsAboveThresholdByPosition[posIdx])-1]

						if len(listOfPlayerIdxsAboveThresholdByPosition[posIdx]) == 0 {
							//Ran out of candidates: couldn't find a suitable candidate this iteration, try the next position
							break
						}

						if assignedPlayers[pickedPlayerInfo] {
							//Try the next player in the list of candidates
							continue
						}

						pickedPlayerInfo.Roles[inningNum] = position
						assignedPlayers[pickedPlayerInfo] = true
						inning.FieldPositions[position] = pickedPlayerInfo

						fmt.Printf("Picked player %s (%s) to play in position %v for inning %d\n", pickedPlayerInfo.FirstName, pickedPlayerInfo.Gender.String(), position.String(), inningNum)
						break
					}

				}

				initialMax -= threshDelta

			}

		}

		verErr := game.verify()
		if verErr == nil {
			break
		}

		fmt.Println("Iterating because verify failed: ", verErr.Error())

		tries++
		if tries > retryThreshold {
			panic("Something went wrong. Too many tries before convergence.")
		}
	}

	verifErr := game.verify()
	if verifErr != nil {
		panic(verifErr)
	}

	return game
}

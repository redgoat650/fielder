package fielder

import (
	"fmt"
	"math/rand"
)

const (
	prefScaleFactor = float64(0.3)
	threshDelta     = float64(0.1)
	benchCredit     = float64(1.0)
	genderDelta     = float64(0.1)
	retryThreshold  = 2500
)

func ScheduleGame(innings int, roster *Roster) *Game {

	err := checkRoster(roster)
	if err != nil {
		panic(err)
	}

	var game *Game
	tries := 0
	maleGenderOffset := 0.0
	femaleGenderOffset := 0.0

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

						//Scale gender offset
						if playerInfo.IsFemale() {
							old := scoringMtx[playerIdx][posIdx]
							scoringMtx[playerIdx][posIdx] += femaleGenderOffset
							fmt.Println("Female", old, scoringMtx[playerIdx][posIdx], femaleGenderOffset)
						} else {
							old := scoringMtx[playerIdx][posIdx]
							scoringMtx[playerIdx][posIdx] += maleGenderOffset
							fmt.Println("Male", old, scoringMtx[playerIdx][posIdx], maleGenderOffset)
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

						fmt.Printf("Picked player %s (%s) to play in position %v for inning %d\n", pickedPlayerInfo.FirstName, pickedPlayerInfo.Gender, position, inningNum)
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

		if genderErr, ok := verErr.(GenderError); ok {
			if genderErr.gender == FemaleGender {
				femaleGenderOffset += genderDelta
			} else {
				maleGenderOffset += genderDelta
			}
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

package fielder

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

//Game is a struct for holding information on a particular game
type Game struct {
	Innings []*Inning
	Roster  *Roster

	Time     time.Time
	TimeDesc string
	OppTeam  string
	Details  string
	WeekNum  int

	Self *Game
}

//NewGame initializes a new Game with the provided number of innings
//and returns its pointer
func NewGame(innings int, weekNo int) *Game {
	game := new(Game)
	game.Self = game
	game.Roster = NewRoster()

	game.Innings = make([]*Inning, 0)

	for inningNum := 0; inningNum < innings; inningNum++ {
		game.Innings = append(game.Innings, NewInning())
	}

	game.WeekNum = weekNo

	return game
}

func (game *Game) SetStartStr(startTime string) {
	game.TimeDesc = startTime
}
func (game *Game) SetOppTeam(oppTeam string) {
	game.OppTeam = oppTeam
}
func (game *Game) SetGameDetails(gameDetails string) {
	game.Details = gameDetails
}

//SetRoster sets the roster for this game
func (game *Game) SetRoster(roster *Roster) {
	game.Roster = roster
}

//String satisfies the stringer interface for Game
func (game Game) String() string {

	s := new(strings.Builder)
	//Print game info:
	s.WriteString(fmt.Sprintf("Game Week %d:\n", game.WeekNum))
	s.WriteString(fmt.Sprintf("Time: %s\n", game.TimeDesc))
	s.WriteString(fmt.Sprintf("Opposing Team: %s\n", game.OppTeam))
	s.WriteString(fmt.Sprintf("Details: %s\n", game.Details))
	s.WriteString("----------------\n")

	//Analysis of players in each position by inning
	for inningNum, inning := range game.Innings {

		s.WriteString(fmt.Sprintf("Inning %d:\n", inningNum+1))

		for pos, player := range inning.FieldPositions {
			if player == nil {
				panic("shouldn't have a nil position")
			}
			s.WriteString(fmt.Sprintf("%s: %s (%s)\n", pos, player.FirstName, player.Gender))
		}

		s.WriteString("----------------\n")

	}

	//Analysis for each player
	mostInnings := 0
	leastInnings := len(game.Innings)
	mostInningsMale := 0
	mostInningsFemale := 0
	leastInningsMale := len(game.Innings)
	leastInningsFemale := len(game.Innings)

	for player := range game.Roster.Players {

		inningsThisPlayer := 0

		for inningNum, role := range player.Roles[game.Self] {
			s.WriteString(fmt.Sprintf("Inning %d: %s plays ", inningNum+1, player.FirstName))
			if role == Bench {
				s.WriteString(fmt.Sprintf("(%s)\n", role))
			} else {
				s.WriteString(fmt.Sprintf("%s\n", role))
			}

			if role != Bench {
				inningsThisPlayer++
			}
		}

		s.WriteString(fmt.Sprintf("%s is playing %d innings\n", player.FirstName, inningsThisPlayer))
		s.WriteString(fmt.Sprintf("----------\n"))

		//Counter metrics
		if inningsThisPlayer > mostInnings {
			mostInnings = inningsThisPlayer
		}
		if inningsThisPlayer < leastInnings {
			leastInnings = inningsThisPlayer
		}
		if player.IsFemale() {
			if inningsThisPlayer > mostInningsFemale {
				mostInningsFemale = inningsThisPlayer
			}
			if inningsThisPlayer < leastInningsFemale {
				leastInningsFemale = inningsThisPlayer
			}
		} else {
			if inningsThisPlayer > mostInningsMale {
				mostInningsMale = inningsThisPlayer
			}
			if inningsThisPlayer < leastInningsMale {
				leastInningsMale = inningsThisPlayer
			}

		}

	}

	s.WriteString(fmt.Sprintf("Most innings played by a player: %d\nLeast innings played by a player: %d\n", mostInnings, leastInnings))
	s.WriteString(fmt.Sprintf("Most innings played by a FEMALE: %d\nLeast innings played by a FEMALE: %d\n", mostInningsFemale, leastInningsFemale))
	s.WriteString(fmt.Sprintf("Most innings played by a MALE: %d\nLeast innings played by a MALE: %d\n", mostInningsMale, leastInningsMale))

	s.WriteString("\n")

	//Analysis of each position by inning

	return s.String()
}

//NumPlayers is a Game method that returns the number of players in the roster for this game
func (game *Game) NumPlayers() int {
	return game.Roster.NumPlayers()
}

//NumInnings is a Game method that returns the number of innings for this game
func (game *Game) NumInnings() int {
	return len(game.Innings)
}

const (
	maxPreferences  = 3
	prefScaleFactor = float64(0.6)
	threshDelta     = float64(0.1)
	benchCredit     = float64(1.0)
	genderDelta     = float64(0.01)
	retryThreshold  = 25000
)

//ScheduleGame is a Game method that schedules positions for all players
//in the Game roster.
func (game *Game) ScheduleGame() error {

	err := checkRoster(game.Roster)
	if err != nil {
		panic(err)
	}

	tries := 0
	maleGenderOffset := 0.0
	femaleGenderOffset := 0.0

	for {

		for playerInfo := range game.Roster.Players {
			playerInfo.Roles[game] = make([]Position, game.NumInnings())
		}

		for inningNum, inning := range game.Innings {

			inning.DropFieldPositions()

			initialMax := float64(0)

			if inningNum == 0 {

				inning.mtx = NewScoringMatrix(game.Roster)

				scoringMtx := inning.mtx.PlayerInfoMap
				for playerInfo, playerPosScores := range scoringMtx {

					//Normalize player preferences
					playerInfo.normalizePrefs()

					for pos := range playerPosScores {

						//Initialize by seniority and skill
						scoringMtx[playerInfo][pos] += playerInfo.Seniority
						scoringMtx[playerInfo][pos] += playerInfo.Skill

						//Initialize by player preference
						for pref, prefStrength := range playerInfo.PrefNorm {
							if pref == pos {
								//Pref strength is a normalized factor between 0 and 1
								scoringMtx[playerInfo][pos] += prefStrength * prefScaleFactor
							}
						}

						//Scale gender offset
						//Pick a random offset up to the gender offset.
						//This counteracts any ping-pong harmonics in the
						//picking algorithm that will cause it to never converge
						//as the male and female scores trade places on each iteration.
						if playerInfo.IsFemale() {
							scoringMtx[playerInfo][pos] += rand.Float64() * femaleGenderOffset
						} else {
							scoringMtx[playerInfo][pos] += rand.Float64() * maleGenderOffset
						}

						if scoringMtx[playerInfo][pos] > initialMax {
							initialMax = scoringMtx[playerInfo][pos]
						}

					}

				}

			} else {

				//Start with last inning's matrix
				inning.mtx = game.Innings[inningNum-1].mtx.copy()

				scoringMtx := inning.mtx.PlayerInfoMap

				for playerInfo := range game.Roster.Players {
					if playerInfo.Roles[game][inningNum-1] == Bench {

						//Player did not play last inning
						for pos := range scoringMtx[playerInfo] {
							scoringMtx[playerInfo][pos] += benchCredit

							//Set the initial max so we can start from an appropriate threshold
							if scoringMtx[playerInfo][pos] > initialMax {
								initialMax = scoringMtx[playerInfo][pos]
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
				if filledPositionCount == NumFieldPositions || filledPositionCount == game.Roster.NumPlayers() {
					//Filled all positions needed
					break
				}

				//Set up a list of player candidates
				listOfPlayerIdxsAboveThresholdByPosition := make(map[Position](map[*Player]struct{}))
				for posIdx := 0; posIdx < NumFieldPositions; posIdx++ {
					pos := posIdx2Position(posIdx)
					listOfPlayerIdxsAboveThresholdByPosition[pos] = make(map[*Player]struct{})
				}

				// for pos := range listOfPlayerIdxsAboveThresholdByPosition {
				// 	listOfPlayerIdxsAboveThresholdByPosition[pos] = make([]int, 0)
				// }

				scoringMtx := inning.mtx.PlayerInfoMap
				for playerInfo := range scoringMtx {

					for pos, score := range scoringMtx[playerInfo] {

						if score >= initialMax {

							if assignedPlayers[playerInfo] {
								//Player has already been assigned to another position
								continue
							}

							listOfPlayerIdxsAboveThresholdByPosition[pos][playerInfo] = struct{}{}

							// fmt.Printf("Player %s is a candidate for position %v because score %f >= threshold %f\n", playerInfo.FirstName, pos, score, initialMax)
						}

					}
				}

				//Each position
				for pos, playerListAtPos := range listOfPlayerIdxsAboveThresholdByPosition {

					//Check all player indexes for this position
					for {

						if len(playerListAtPos) == 0 {
							// fmt.Println("Couldn't find a suitable candidate for", pos)
							//Couldn't find a suitable candidate at this threshold level
							break
						}

						if _, ok := inning.FieldPositions[pos]; ok {
							//We've already picked a player for this position
							// fmt.Println("Already picked a player for position", pos)
							break
						}

						//Pick a player from the list of players that are above the threshold for this position
						pickedListIdx := rand.Intn(len(playerListAtPos))
						var pickedPlayerInfo *Player
						for pickedPlayerInfo = range playerListAtPos {
							if pickedListIdx <= 0 {
								break
							}
							pickedListIdx--
						}

						// fmt.Println("Picked", pickedPlayerInfo.FirstName)

						delete(listOfPlayerIdxsAboveThresholdByPosition[pos], pickedPlayerInfo)

						if assignedPlayers[pickedPlayerInfo] {
							//Try the next player in the list of candidates
							// fmt.Println("Player already assigned for position", pos, "...try next")
							continue
						}

						pickedPlayerInfo.Roles[game][inningNum] = pos
						assignedPlayers[pickedPlayerInfo] = true
						inning.FieldPositions[pos] = pickedPlayerInfo

						// fmt.Printf("Picked player %s (%s) to play in position %v for inning %d\n", pickedPlayerInfo.FirstName, pickedPlayerInfo.Gender, pos, inningNum)
						break
					}

				}

				initialMax -= threshDelta

			}

		}

		verErr := verifyGame(game)
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

	verifErr := verifyGame(game)
	if verifErr != nil {
		panic(verifErr)
	}

	fmt.Println("Game score is", game.ScoreGame())

	return nil
}

func (game *Game) ScoreGame() float64 {

	score := 0.0

	for player := range game.Roster.Players {

		player.normalizePrefs()

		playerScore := 0.0

		participationDecay := 0.0
		participationPenaltyFactor := 0.9

		playingBonus := 1.0

		for _, pos := range player.Roles[game] {

			if pos != Bench {

				//Calculate the score for this player
				playerScore += (playingBonus + player.PrefNorm[pos]) * math.Pow(participationPenaltyFactor, participationDecay)

				participationDecay++

			}
		}

		score += playerScore
	}

	return score
}

//--------------------------------------------
// Second Algorithm

type ScratchGame struct {
	b          []int
	numInnings int
}

const (
	Unfilled = -1
)

func NewScratchGame(numInnings int) *ScratchGame {
	s := new(ScratchGame)

	s.b = make([]int, numInnings*NumFieldPositions)
	for i := range s.b {
		s.b[i] = Unfilled
	}

	s.numInnings = numInnings

	return s
}

type ScratchInning []int

func (game *ScratchGame) GetInning(inningNum int) ScratchInning {
	return ScratchInning(game.b[inningNum*NumFieldPositions : (inningNum+1)*NumFieldPositions])
}

// type fieldPos int

// const (
// 	PitcherIdx fieldPos = iota
// 	CatcherIdx
// 	FirstIdx
// 	SecondIdx
// 	ThirdIdx
// 	LShortIdx
// 	RShortIdx
// 	LFieldIdx
// 	LCenterIdx
// 	RCenterIdx
// 	RFieldIdx
// 	NumFieldPosIdx int = iota
// )

func (game *ScratchGame) GetPlayerIdxAtPos(inningNum int, posIdx int) int {
	return game.GetInning(inningNum)[posIdx]
}

func (game *ScratchGame) verifyGame(femaleLookup []bool) error {
	for innNum := 0; innNum < game.numInnings; innNum++ {
		maleCount := 0
		femaleCount := 0
		for posIdx := 0; posIdx < NumFieldPositions; posIdx++ {
			playerIdxAtPos := game.GetPlayerIdxAtPos(innNum, posIdx)
			if playerIdxAtPos != Unfilled {
				if femaleLookup[playerIdxAtPos] {
					femaleCount++
				} else {
					maleCount++
				}
			}
		}

		if maleCount < MinGenderCount || femaleCount < MinGenderCount {
			return fmt.Errorf("Not valid positions")
		}
	}

	return nil
}

func (game *ScratchGame) ScoreGame(prefLookup [][]float64, playerTakenTable [][]int, skillLookup []float64) float64 {

	numPlayers := len(prefLookup)
	inningsPlayedLookup := make([]int, numPlayers)

	score := 0.0
	for inningNum := 0; inningNum < game.numInnings; inningNum++ {

		for posIdx := 0; posIdx < NumFieldPositions; posIdx++ {

			playerIdxAtPosition := game.GetPlayerIdxAtPos(inningNum, posIdx)

			posPrefListForPlayer := prefLookup[playerIdxAtPosition]

			prefScore := posPrefListForPlayer[posIdx]

			innPlayed := inningsPlayedLookup[playerIdxAtPosition]

			participationPenaltyFactor := 0.9
			playingBonus := 1.0

			skillFactor := skillLookup[playerIdxAtPosition]

			// inningDistanceFactor := 0.2
			// inningDistance := game.numInnings
			// for checkInn := 0; checkInn < game.numInnings; checkInn++ {
			// 	// if game.
			// 	posInThisInning := playerTakenTable[checkInn][playerIdxAtPosition]
			// 	if posInThisInning != Unfilled {
			// 		distFromThisInning := int(math.Abs(float64(inningNum - checkInn)))
			// 		if inningDistance > distFromThisInning {
			// 			inningDistance = distFromThisInning
			// 		}
			// 	}
			// }
			// innDistScore := float64(inningDistance) * inningDistanceFactor

			playBonusScale := 6.0
			prefScoreScale := 6.0
			skillFactorScale := 1.0
			sum := playBonusScale + prefScoreScale + skillFactorScale

			//Player is happier if
			//1. they play in an inning
			//2. they play a position they really want
			//3. they have played fewer innings already
			playerScoreThisInn := (playingBonus*playBonusScale/sum +
				prefScore*prefScoreScale/sum +
				skillFactor*skillFactorScale/sum)

			playerScoreScaledByPlayTime := playerScoreThisInn * math.Pow(participationPenaltyFactor, float64(innPlayed))

			score += playerScoreScaledByPlayTime

			inningsPlayedLookup[playerIdxAtPosition]++

		}

	}

	return score
}

//ScheduleGame2 is a Game method that schedules positions for all players
//in the Game roster.
func (game *Game) ScheduleGame2() error {

	err := checkRoster(game.Roster)
	if err != nil {
		return err
	}

	//Lookup table to give an index to a Player
	// []*Player
	playerLookup := make([]*Player, game.Roster.NumPlayers())

	//Preference lookup by field pos idx, listed by playerIdx
	// [] ([]float64)
	//playerIdx -> posIdx -> score
	prefLookup := make([][]float64, game.Roster.NumPlayers())
	skillLookup := make([]float64, game.Roster.NumPlayers())

	genderFemaleLookup := make([]bool, game.Roster.NumPlayers())
	malePlayerIdxList := make([]int, 0)
	femalePlayerIdxList := make([]int, 0)
	playerIdxList := make([]int, 0)

	i := 0
	for player := range game.Roster.Players {
		playerLookup[i] = player

		prefLookup[i] = make([]float64, NumFieldPositions)

		player.normalizePrefs()
		for pos, normMag := range player.PrefNorm {

			posIdx := position2PosIdx(pos)
			prefLookup[i][posIdx] = normMag

		}

		skillLookup[i] = player.Skill

		isFemale := player.IsFemale()
		genderFemaleLookup[i] = isFemale
		if isFemale {
			femalePlayerIdxList = append(femalePlayerIdxList, i)
		} else {
			malePlayerIdxList = append(malePlayerIdxList, i)
		}

		playerIdxList = append(playerIdxList, i)

		i++
	}

	//Start filling in stuff

	numInnings := 5
	scratchGame := NewScratchGame(numInnings)

	//Fill the game with stuff:

	playerTakenTable := make([][]int, numInnings)
	for i := 0; i < numInnings; i++ {
		playerTakenTable[i] = make([]int, game.Roster.NumPlayers())
		//Initialize to all players in untaken positions
		for j := range playerTakenTable[i] {
			playerTakenTable[i][j] = Unfilled
		}
	}

	for inningNum := 0; inningNum < numInnings; inningNum++ {

		scrInn := scratchGame.GetInning(inningNum)

		for fieldPosCount := 0; fieldPosCount < NumFieldPositions; fieldPosCount++ {

			//If we want to fill positions by priority we index the list from
			//0-highest priority
			//...
			//otherwise pick a random index
			fieldPosIdx := fieldPosCount

			//Fill each field position
			var genderErr error
			switch {
			case fieldPosCount < 4:
				scrInn[fieldPosIdx], genderErr = pickRandAvailPlayerIdx(malePlayerIdxList, playerTakenTable[inningNum], fieldPosIdx)
			case fieldPosCount < 8:
				scrInn[fieldPosIdx], genderErr = pickRandAvailPlayerIdx(femalePlayerIdxList, playerTakenTable[inningNum], fieldPosIdx)
			default:
				scrInn[fieldPosIdx], _ = pickRandAvailPlayerIdx(playerIdxList, playerTakenTable[inningNum], fieldPosIdx)
			}

			//Didn't have enough girls or guys to fill a valid field
			if genderErr != nil {
				return genderErr
			}

		}

	}

	oldScore := scratchGame.ScoreGame(prefLookup, playerTakenTable, skillLookup)
	callConvergeAt := 100000

	convCount := 0
	for {

		pickRandInningNum := rand.Intn(numInnings)
		scratchInning := scratchGame.GetInning(pickRandInningNum)

		pickNewPlayer := pickRandPlayerIdx(playerIdxList)
		pickRandPosition := pickRandFilledPosIdx(scratchInning)

		swapPlayer := scratchInning[pickRandPosition]
		if swapPlayer == Unfilled {
			panic("How did we get an unfilled position?")
		}
		if swapPlayer == pickNewPlayer {
			continue
		}

		newPlayerPrevPosIdx := playerTakenTable[pickRandInningNum][pickNewPlayer]

		//Set the new player to that position in the scratch inning
		scratchInning[pickRandPosition] = pickNewPlayer
		//Update player taken table with new player's position
		playerTakenTable[pickRandInningNum][pickNewPlayer] = pickRandPosition

		//Set the swapped player's position to the new player's old position
		if newPlayerPrevPosIdx != Unfilled {
			scratchInning[newPlayerPrevPosIdx] = swapPlayer
		}
		//Update player taken table with swapped player's position
		playerTakenTable[pickRandInningNum][swapPlayer] = newPlayerPrevPosIdx

		//Recalculate score
		newScore := scratchGame.ScoreGame(prefLookup, playerTakenTable, skillLookup)

		//increment attempts at convergence
		convCount++

		//If score is smaller or field layout is invalid, swap back
		if newScore < oldScore || scratchGame.verifyGame(genderFemaleLookup) != nil {
			scratchInning[pickRandPosition] = swapPlayer
			playerTakenTable[pickRandInningNum][swapPlayer] = pickRandPosition

			if newPlayerPrevPosIdx != Unfilled {
				scratchInning[newPlayerPrevPosIdx] = pickNewPlayer
			}
			playerTakenTable[pickRandInningNum][pickNewPlayer] = newPlayerPrevPosIdx

		} else if newScore > oldScore {
			fmt.Printf("Score %f -> %f\n", oldScore, newScore)
			oldScore = newScore
			//Got a higher score so reset the convergence counter
			convCount = 0
		}

		if convCount >= callConvergeAt {
			break
		}

	}

	fmt.Println("Done!")

	game.fillFromScratch(scratchGame, playerLookup)

	return nil
}

func (game *Game) fillFromScratch(scrGame *ScratchGame, playerLookup []*Player) {

	for playerInfo := range game.Roster.Players {
		playerInfo.Roles[game] = make([]Position, game.NumInnings())
	}

	for inningNum := 0; inningNum < scrGame.numInnings; inningNum++ {

		for posIdx := 0; posIdx < NumFieldPositions; posIdx++ {

			playerIdx := scrGame.GetPlayerIdxAtPos(inningNum, posIdx)
			if playerIdx == Unfilled {
				continue
			}

			player := playerLookup[playerIdx]
			pos := posIdx2Position(posIdx)

			game.Innings[inningNum].FieldPositions[pos] = player

			fmt.Println(len(player.Roles))
			fmt.Println(len(player.Roles[game]))
			player.Roles[game][inningNum] = pos

		}
	}
}

func pickRandFilledPosIdx(inn ScratchInning) int {
	for {
		randPosIdx := rand.Intn(len(inn))
		if inn[randPosIdx] != Unfilled {
			return randPosIdx
		}
	}
}

func pickRandPlayerIdx(playerIdxList []int) int {
	return playerIdxList[rand.Intn(len(playerIdxList))]
}

func pickRandAvailPlayerIdx(playerIdxList []int, playerPosTable []int, fillPos int) (int, error) {

	somethingAvail := false
	for _, plIdx := range playerIdxList {

		if playerPosTable[plIdx] == Unfilled {
			somethingAvail = true
			break
		}
	}
	if somethingAvail == false {
		return -1, fmt.Errorf("No players are available from this list")
	}

	for {
		pickIdx := pickRandPlayerIdx(playerIdxList)
		if playerPosTable[pickIdx] == Unfilled {
			playerPosTable[pickIdx] = fillPos
			return pickIdx, nil
		}
	}
}

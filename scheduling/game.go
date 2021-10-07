package fielder

import (
	"errors"
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
}

//NewGame initializes a new Game with the provided number of innings
//and returns its pointer
func NewGame(innings int, weekNo int) *Game {
	game := new(Game)
	game.Roster = NewRoster()

	game.Innings = make([]*Inning, 0)

	for inningNum := 0; inningNum < innings; inningNum++ {
		game.Innings = append(game.Innings, NewInning())
	}

	game.WeekNum = weekNo

	return game
}

// SetStartStr sets the start time to the provided string
func (game *Game) SetStartStr(startTime string) {
	game.TimeDesc = startTime
}

// SetOppTeam sets the opponent team name
func (game *Game) SetOppTeam(oppTeam string) {
	game.OppTeam = oppTeam
}

// SetGameDetails sets the game details
func (game *Game) SetGameDetails(gameDetails string) {
	game.Details = gameDetails
}

//SetRoster sets the roster for this game
func (game *Game) SetRoster(roster *Roster) {
	game.Roster = roster
}

//String satisfies the stringer interface for Game
func (game *Game) String() string {

	s := new(strings.Builder)
	//Print game info:
	s.WriteString(fmt.Sprintf("Game Week %d:\n", game.WeekNum))
	s.WriteString(fmt.Sprintf("Time: %s\n", game.TimeDesc))
	s.WriteString(fmt.Sprintf("Opposing Team: %s\n", game.OppTeam))
	s.WriteString(fmt.Sprintf("Details: %s\n", game.Details))
	s.WriteString("----------------\n")

	//Debug info on players
	for _, player := range game.Roster.Players {
		s.WriteString(player.String())
		s.WriteString("----------------\n")
	}

	//Analysis of players in each position by inning
	for inningNum, inning := range game.Innings {
		s.WriteString(fmt.Sprintf("Inning %d:\n", inningNum+1))

		for pos, player := range inning.FieldPositions {
			if player == nil {
				s.WriteString(fmt.Sprintf("%s: NONE", pos))
				continue
			}

			s.WriteString(fmt.Sprintf("%s: %s (%s)\n", pos, player.Name, player.Gender))
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
	mostBenchInARow := 0

	for _, player := range game.Roster.Players {

		inningsThisPlayer := 0

		for inningNum, inning := range game.Innings {
			role := inning.FindPlayerPosition(player)

			s.WriteString(fmt.Sprintf("Inning %d: %s plays ", inningNum+1, player.Name))
			if role == Bench {
				s.WriteString(fmt.Sprintf("(%s)\n", role))
			} else {
				s.WriteString(fmt.Sprintf("%s\n", role))
			}

			if role != Bench {
				inningsThisPlayer++
			}
		}
		benchInARowThisPlayer := game.calcBenchInARowByPlayer(player)
		if benchInARowThisPlayer > mostBenchInARow {
			mostBenchInARow = benchInARowThisPlayer
		}

		s.WriteString(fmt.Sprintf("%s is playing %d innings\n", player.Name, inningsThisPlayer))
		s.WriteString(fmt.Sprintf("and is on the BENCH %v times in a row\n", benchInARowThisPlayer))
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

	for _, pos := range fieldPosList {
		s.WriteString(fmt.Sprintf("%v: ", pos))
		pls := make([]string, 0)
		for _, inning := range game.Innings {
			player := inning.FieldPositions[pos]
			pls = append(pls, player.Name)
		}
		s.WriteString(fmt.Sprintf("%v\n", strings.Join(pls, ", ")))
	}
	s.WriteString("----------------\n")

	s.WriteString(fmt.Sprintf("Most innings played by a player: %d\nLeast innings played by a player: %d\n", mostInnings, leastInnings))
	s.WriteString(fmt.Sprintf("Most innings played by a FEMALE: %d\nLeast innings played by a FEMALE: %d\n", mostInningsFemale, leastInningsFemale))
	s.WriteString(fmt.Sprintf("Most innings played by a MALE: %d\nLeast innings played by a MALE: %d\n", mostInningsMale, leastInningsMale))
	s.WriteString(fmt.Sprintf("Most innings in a row on a bench: %d\n", mostBenchInARow))

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
// func (game *Game) ScheduleGame() error {

// 	err := checkRoster(game.Roster)
// 	if err != nil {
// 		return err
// 	}

// 	tries := 0
// 	maleGenderOffset := 0.0
// 	femaleGenderOffset := 0.0

// 	for {

// 		for playerInfo := range game.Roster.Players {
// 			playerInfo.Roles[game] = make([]Position, game.NumInnings())
// 		}

// 		for inningNum, inning := range game.Innings {

// 			inning.DropFieldPositions()

// 			initialMax := float64(0)

// 			if inningNum == 0 {

// 				inning.mtx = NewScoringMatrix(game.Roster)

// 				scoringMtx := inning.mtx.PlayerInfoMap
// 				for playerInfo, playerPosScores := range scoringMtx {

// 					//Normalize player preferences
// 					playerInfo.normalizePrefs()

// 					for pos := range playerPosScores {

// 						//Initialize by seniority and skill
// 						scoringMtx[playerInfo][pos] += playerInfo.Seniority
// 						scoringMtx[playerInfo][pos] += playerInfo.Skill

// 						//Initialize by player preference
// 						for pref, prefStrength := range playerInfo.PrefNorm {
// 							if pref == pos {
// 								//Pref strength is a normalized factor between 0 and 1
// 								scoringMtx[playerInfo][pos] += prefStrength * prefScaleFactor
// 							}
// 						}

// 						//Scale gender offset
// 						//Pick a random offset up to the gender offset.
// 						//This counteracts any ping-pong harmonics in the
// 						//picking algorithm that will cause it to never converge
// 						//as the male and female scores trade places on each iteration.
// 						if playerInfo.IsFemale() {
// 							scoringMtx[playerInfo][pos] += rand.Float64() * femaleGenderOffset
// 						} else {
// 							scoringMtx[playerInfo][pos] += rand.Float64() * maleGenderOffset
// 						}

// 						if scoringMtx[playerInfo][pos] > initialMax {
// 							initialMax = scoringMtx[playerInfo][pos]
// 						}

// 					}

// 				}

// 			} else {

// 				//Start with last inning's matrix
// 				inning.mtx = game.Innings[inningNum-1].mtx.copy()

// 				scoringMtx := inning.mtx.PlayerInfoMap

// 				for playerInfo := range game.Roster.Players {
// 					if playerInfo.Roles[game][inningNum-1] == Bench {

// 						//Player did not play last inning
// 						for pos := range scoringMtx[playerInfo] {
// 							scoringMtx[playerInfo][pos] += benchCredit

// 							//Set the initial max so we can start from an appropriate threshold
// 							if scoringMtx[playerInfo][pos] > initialMax {
// 								initialMax = scoringMtx[playerInfo][pos]
// 							}
// 						}
// 					}
// 				}
// 			}

// 			//We're going to establish a threshold and continue lowering it
// 			//and assigning players to positions until all positions are filled.
// 			assignedPlayers := make(map[*Player]bool)

// 			for {

// 				//See if we've filled all the positions and can break out
// 				filledPositionCount, _ := inning.CountPlayersOnField()
// 				if filledPositionCount == NumFieldPositions || filledPositionCount == game.Roster.NumPlayers() {
// 					//Filled all positions needed
// 					break
// 				}

// 				//Set up a list of player candidates
// 				listOfPlayerIdxsAboveThresholdByPosition := make(map[Position](map[*Player]struct{}))
// 				for posIdx := 0; posIdx < NumFieldPositions; posIdx++ {
// 					pos := posIdx2Position(posIdx)
// 					listOfPlayerIdxsAboveThresholdByPosition[pos] = make(map[*Player]struct{})
// 				}

// 				// for pos := range listOfPlayerIdxsAboveThresholdByPosition {
// 				// 	listOfPlayerIdxsAboveThresholdByPosition[pos] = make([]int, 0)
// 				// }

// 				scoringMtx := inning.mtx.PlayerInfoMap
// 				for playerInfo := range scoringMtx {

// 					for pos, score := range scoringMtx[playerInfo] {

// 						if score >= initialMax {

// 							if assignedPlayers[playerInfo] {
// 								//Player has already been assigned to another position
// 								continue
// 							}

// 							listOfPlayerIdxsAboveThresholdByPosition[pos][playerInfo] = struct{}{}

// 							// fmt.Printf("Player %s is a candidate for position %v because score %f >= threshold %f\n", playerInfo.FirstName, pos, score, initialMax)
// 						}

// 					}
// 				}

// 				//Each position
// 				for pos, playerListAtPos := range listOfPlayerIdxsAboveThresholdByPosition {

// 					//Check all player indexes for this position
// 					for {

// 						if len(playerListAtPos) == 0 {
// 							// fmt.Println("Couldn't find a suitable candidate for", pos)
// 							//Couldn't find a suitable candidate at this threshold level
// 							break
// 						}

// 						if _, ok := inning.FieldPositions[pos]; ok {
// 							//We've already picked a player for this position
// 							// fmt.Println("Already picked a player for position", pos)
// 							break
// 						}

// 						//Pick a player from the list of players that are above the threshold for this position
// 						pickedListIdx := rand.Intn(len(playerListAtPos))
// 						var pickedPlayerInfo *Player
// 						for pickedPlayerInfo = range playerListAtPos {
// 							if pickedListIdx <= 0 {
// 								break
// 							}
// 							pickedListIdx--
// 						}

// 						// fmt.Println("Picked", pickedPlayerInfo.FirstName)

// 						delete(listOfPlayerIdxsAboveThresholdByPosition[pos], pickedPlayerInfo)

// 						if assignedPlayers[pickedPlayerInfo] {
// 							//Try the next player in the list of candidates
// 							// fmt.Println("Player already assigned for position", pos, "...try next")
// 							continue
// 						}

// 						pickedPlayerInfo.Roles[game][inningNum] = pos
// 						assignedPlayers[pickedPlayerInfo] = true
// 						inning.FieldPositions[pos] = pickedPlayerInfo

// 						// fmt.Printf("Picked player %s (%s) to play in position %v for inning %d\n", pickedPlayerInfo.FirstName, pickedPlayerInfo.Gender, pos, inningNum)
// 						break
// 					}

// 				}

// 				initialMax -= threshDelta

// 			}

// 		}

// 		verErr := verifyGame(game)
// 		if verErr == nil {
// 			break
// 		}

// 		if genderErr, ok := verErr.(GenderError); ok {
// 			if genderErr.gender == FemaleGender {
// 				femaleGenderOffset += genderDelta
// 			} else {
// 				maleGenderOffset += genderDelta
// 			}
// 		}

// 		fmt.Println("Iterating because verify failed: ", verErr.Error())

// 		tries++
// 		if tries > retryThreshold {
// 			panic("Something went wrong. Too many tries before convergence.")
// 		}
// 	}

// 	verifErr := verifyGame(game)
// 	if verifErr != nil {
// 		panic(verifErr)
// 	}

// 	fmt.Println("Game score is", game.ScoreGame())

// 	return nil
// }

// // ScoreGame calculates the score for a game
// func (game *Game) ScoreGame() float64 {

// 	score := 0.0

// 	for player := range game.Roster.Players {

// 		player.normalizePrefs()

// 		playerScore := 0.0

// 		participationDecay := 0.0
// 		participationPenaltyFactor := 0.9

// 		playingBonus := 1.0

// 		for _, pos := range player.Roles[game] {

// 			if pos != Bench {

// 				//Calculate the score for this player
// 				playerScore += (playingBonus + player.PrefNorm[pos]) * math.Pow(participationPenaltyFactor, participationDecay)

// 				participationDecay++

// 			}
// 		}

// 		score += playerScore
// 	}

// 	return score
// }

//--------------------------------------------
// Second Algorithm

// ScratchGame is a scratch game
type ScratchGame struct {
	b          []int
	numInnings int
}

const (
	//unfilled is a special value for describing an unfilled position
	unfilled = -1
)

// NewScratchGame initializes and returns a pointer to a new ScratchGame
func NewScratchGame(numInnings int) *ScratchGame {
	s := new(ScratchGame)

	s.b = make([]int, numInnings*NumFieldPositions)
	for i := range s.b {
		s.b[i] = unfilled
	}

	s.numInnings = numInnings

	return s
}

// ScratchInning is a scratch inning
type ScratchInning []int

// GetInning gets a ScratchInning from a ScratchGame by inning number
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

// GetPlayerIdxAtPos gets a player index at a given position for a given
// inning number
func (game *ScratchGame) GetPlayerIdxAtPos(inningNum int, posIdx int) int {
	return game.GetInning(inningNum)[posIdx]
}

func (game *ScratchGame) verifyGame(femaleLookup []bool) error {
	for innNum := 0; innNum < game.numInnings; innNum++ {
		maleCount := 0
		femaleCount := 0
		for posIdx := 0; posIdx < NumFieldPositions; posIdx++ {
			playerIdxAtPos := game.GetPlayerIdxAtPos(innNum, posIdx)
			if playerIdxAtPos != unfilled {
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

type BetterGame struct {
	Innings []*BetterInning
	Roster  BetterRoster
}

type BetterInning struct {
	PositionsMap map[Position]*Player
}

type ScoringParams struct {
	Player *Player

	CptPref   map[Position]int
	Skill     float64
	Seniority float64

	cptPrefNorm    map[Position]float64
	playerPrefNorm map[Position]float64
}

func NewBetterGame(numInnings int) *BetterGame {
	game := new(BetterGame)

	for i := 0; i < numInnings; i++ {
		game.Innings = append(game.Innings, &BetterInning{
			PositionsMap: make(map[Position]*Player),
		})
	}

	return game
}

func (game *BetterGame) String() string {
	s := new(strings.Builder)
	for inningNum, inning := range game.Innings {
		s.WriteString(fmt.Sprintf("Inning %d:\n", inningNum+1))

		for pos, player := range inning.PositionsMap {
			if player == nil {
				s.WriteString(fmt.Sprintf("%s: NONE", pos))
				continue
			}

			s.WriteString(fmt.Sprintf("%s: %s (%s)\n", pos, player.Name, player.Gender))
		}

		s.WriteString("----------------\n")

	}

	mostInnings := 0
	leastInnings := len(game.Innings)
	mostInningsMale := 0
	mostInningsFemale := 0
	leastInningsMale := len(game.Innings)
	leastInningsFemale := len(game.Innings)
	mostBenchInARow := 0

	for _, playerScoringParams := range game.Roster {
		player := playerScoringParams.Player

		inningsThisPlayer := 0

		for inningNum, inning := range game.Innings {
			role := inning.FindPlayerPosition(player)

			s.WriteString(fmt.Sprintf("Inning %d: %s plays ", inningNum+1, player.Name))
			if role == Bench {
				s.WriteString(fmt.Sprintf("(%s)\n", role))
			} else {
				s.WriteString(fmt.Sprintf("%s\n", role))
			}

			if role != Bench {
				inningsThisPlayer++
			}
		}
		benchInARowThisPlayer := game.calcBenchInARowForPlayer(player)
		if benchInARowThisPlayer > mostBenchInARow {
			mostBenchInARow = benchInARowThisPlayer
		}

		s.WriteString(fmt.Sprintf("%s is playing %d innings\n", player.Name, inningsThisPlayer))
		s.WriteString(fmt.Sprintf("and is on the BENCH %v times in a row\n", benchInARowThisPlayer))
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

	for _, pos := range fieldPosList {
		s.WriteString(fmt.Sprintf("%v: ", pos))
		pls := make([]string, 0)
		for _, inning := range game.Innings {
			player := inning.PositionsMap[pos]
			pls = append(pls, player.Name)
		}
		s.WriteString(fmt.Sprintf("%v\n", strings.Join(pls, ", ")))
	}
	s.WriteString("----------------\n")

	s.WriteString(fmt.Sprintf("Most innings played by a player: %d\nLeast innings played by a player: %d\n", mostInnings, leastInnings))
	s.WriteString(fmt.Sprintf("Most innings played by a FEMALE: %d\nLeast innings played by a FEMALE: %d\n", mostInningsFemale, leastInningsFemale))
	s.WriteString(fmt.Sprintf("Most innings played by a MALE: %d\nLeast innings played by a MALE: %d\n", mostInningsMale, leastInningsMale))
	s.WriteString(fmt.Sprintf("Most innings in a row on a bench: %d\n", mostBenchInARow))

	s.WriteString("\n")

	return s.String()
}

func (game *BetterGame) ScoreGame(scoringParams map[*Player]ScoringParams) float64 {
	inningsPlayedLookup := make(map[*Player]int)

	score := 0.0

	for _, inning := range game.Innings {
		for pos, player := range inning.PositionsMap {
			p := scoringParams[player]

			prefScore := p.playerPrefNorm[pos]
			cptPrefScore := p.cptPrefNorm[pos]

			skillFactor := p.Skill
			seniorityFactor := p.Seniority

			innPlayed := inningsPlayedLookup[player]

			playBonusScale := 6.0
			prefScoreScale := 6.0
			cptPrefScoreScale := 6.0
			skillFactorScale := 1.0
			seniorityFactorScale := 1.0
			participationPenaltyFactor := 0.8
			playingBonus := 1.0
			sum := playBonusScale + prefScoreScale + cptPrefScoreScale + skillFactorScale + seniorityFactorScale

			//Player is happier if
			//1. they play in an inning
			//2. they play a position they really want
			//3. they have played fewer innings already
			playerScoreThisInn := (playingBonus*playBonusScale/sum +
				prefScore*prefScoreScale/sum +
				cptPrefScore*cptPrefScoreScale/sum +
				skillFactor*skillFactorScale/sum +
				seniorityFactor*seniorityFactorScale/sum)

			playerScoreScaledByPlayTime := playerScoreThisInn * math.Pow(participationPenaltyFactor, float64(innPlayed))

			score += playerScoreScaledByPlayTime

			inningsPlayedLookup[player]++
		}
	}

	return score
}

type BetterRoster []ScoringParams

func makeScoringLookup(roster BetterRoster) map[*Player]ScoringParams {
	ret := make(map[*Player]ScoringParams)

	for _, scoringParams := range roster {
		ret[scoringParams.Player] = scoringParams
	}

	return ret
}

var PositionList = []Position{
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

func randomPosition() Position {
	return PositionList[rand.Intn(len(PositionList))]
}

func (game *BetterGame) populateRandomly(roster BetterRoster) {
	for _, inning := range game.Innings {
		for _, pos := range PositionList {
			inning.PositionsMap[pos] = nil
		}

		var availPlayers = make(map[*Player]struct{})
		var availMale = make(map[*Player]struct{})
		var availFemale = make(map[*Player]struct{})

		for _, playerScoringParams := range roster {
			player := playerScoringParams.Player
			availPlayers[player] = struct{}{}

			if player.IsFemale() {
				availFemale[player] = struct{}{}
			} else {
				availMale[player] = struct{}{}
			}
		}

		malesNeeded := MinGenderCount
		femalesNeeded := MinGenderCount

		for pos := range inning.PositionsMap {
			var selectedPlayer *Player

			switch {
			case malesNeeded > 0:
				selectedPlayer = randInMap(availMale)
				delete(availMale, selectedPlayer)
				malesNeeded--

			case femalesNeeded > 0:
				selectedPlayer = randInMap(availFemale)
				delete(availFemale, selectedPlayer)
				femalesNeeded--

			default:
				selectedPlayer = randInMap(availPlayers)
			}

			delete(availPlayers, selectedPlayer)

			inning.PositionsMap[pos] = selectedPlayer
		}
	}
}

func randInMap(m map[*Player]struct{}) *Player {
	n := rand.Intn(len(m))

	for p := range m {
		if n <= 0 {
			return p
		}
		n--
	}

	return nil
}

func (game *BetterGame) ScheduleGame(roster BetterRoster) error {
	game.Roster = roster

	scoringLookup := makeScoringLookup(roster)

	game.populateRandomly(roster)

	oldScore := game.ScoreGame(scoringLookup)

	callConvergeAt := 1000000

	convCount := 0

	for {
		randInning := game.Innings[rand.Intn(len(game.Innings))]
		randPosition := randomPosition()
		swapPlayer := randInning.PositionsMap[randPosition]

		if swapPlayer == nil {
			continue
		}

		randRosterPlayer := roster[rand.Intn(len(roster))].Player

		prevPos := Bench
		for pos, playerAtPos := range randInning.PositionsMap {
			if playerAtPos == randRosterPlayer {
				randInning.PositionsMap[pos] = swapPlayer
				// fmt.Println("swapping", swapPlayer.Name, "to", pos, ", which is", randRosterPlayer.Name, "'s original position")
				prevPos = pos

				break
			}
		}

		randInning.PositionsMap[randPosition] = randRosterPlayer
		// fmt.Println("trying", randRosterPlayer.Name, "at", randPosition)

		newScore := game.ScoreGame(scoringLookup)

		convCount++

		if newScore < oldScore || game.verifyGame() != nil {
			// Swap back
			// fmt.Println("Swapping back")
			if prevPos != Bench {
				randInning.PositionsMap[prevPos] = randRosterPlayer
				// fmt.Println("Swapping back", randRosterPlayer.Name, "to", prevPos)
			}

			randInning.PositionsMap[randPosition] = swapPlayer
			// fmt.Println("Swapping back", swapPlayer.Name, "to", randPosition)

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

	game.minimizeBenchTime()

	return nil
}

func (game *BetterGame) minimizeBenchTime() {
	fmt.Println("Start bench minimization")

	numInnings := len(game.Innings)

	if numInnings == 1 {
		fmt.Println("Only one inning. No work to do.")
		return
	}

	startVal := game.calcBenchInARow()
outerLoop:
	for {
		fmt.Printf("Starting at %v\n", startVal)
		for i := 0; i < 100000; i++ {
			pick1 := rand.Intn(numInnings)
			pick2 := rand.Intn(numInnings)
			for pick1 == pick2 {
				pick2 = rand.Intn(numInnings)
			}
			game.Innings[pick1], game.Innings[pick2] = game.Innings[pick2], game.Innings[pick1]
			finalVal := game.calcBenchInARow()
			if finalVal < startVal {
				startVal = finalVal
				continue outerLoop
			}
			if finalVal > startVal {
				game.Innings[pick1], game.Innings[pick2] = game.Innings[pick2], game.Innings[pick1]
			}
		}
		fmt.Printf("Ending at %v\n", startVal)
		return
	}
}

func (game *BetterGame) verifyGame() error {
	for _, inning := range game.Innings {
		maleCount := 0
		femaleCount := 0

		for _, player := range inning.PositionsMap {
			if player != nil {
				if player.IsFemale() {
					femaleCount++
				} else {
					maleCount++
				}
			}
		}

		if maleCount < MinGenderCount || femaleCount < MinGenderCount {
			return errors.New("not valid positions")
		}
	}

	return nil
}

// ScoreGame scores a scratch game
func (game *ScratchGame) ScoreGame(prefLookup [][]float64, cptPrefLookup [][]float64, playerTakenTable [][]int, skillLookup []float64, seniorityLookup []float64) float64 {

	numPlayers := len(prefLookup)
	inningsPlayedLookup := make([]int, numPlayers)

	score := 0.0
	for inningNum := 0; inningNum < game.numInnings; inningNum++ {

		for posIdx := 0; posIdx < NumFieldPositions; posIdx++ {

			playerIdxAtPosition := game.GetPlayerIdxAtPos(inningNum, posIdx)

			posPrefListForPlayer := prefLookup[playerIdxAtPosition]
			cptPrefListForPlayer := cptPrefLookup[playerIdxAtPosition]

			prefScore := posPrefListForPlayer[posIdx]
			cptPrefScore := cptPrefListForPlayer[posIdx]

			innPlayed := inningsPlayedLookup[playerIdxAtPosition]

			skillFactor := skillLookup[playerIdxAtPosition]

			seniorityFactor := seniorityLookup[playerIdxAtPosition]

			// inningDistanceFactor := 0.2
			// inningDistance := game.numInnings
			// for checkInn := 0; checkInn < game.numInnings; checkInn++ {
			// 	// if game.
			// 	posInThisInning := playerTakenTable[checkInn][playerIdxAtPosition]
			// 	if posInThisInning != unfilled {
			// 		distFromThisInning := int(math.Abs(float64(inningNum - checkInn)))
			// 		if inningDistance > distFromThisInning {
			// 			inningDistance = distFromThisInning
			// 		}
			// 	}
			// }
			// innDistScore := float64(inningDistance) * inningDistanceFactor

			playBonusScale := 6.0
			prefScoreScale := 6.0
			cptPrefScoreScale := 6.0
			skillFactorScale := 1.0
			seniorityFactorScale := 1.0
			participationPenaltyFactor := 0.8
			playingBonus := 1.0
			sum := playBonusScale + prefScoreScale + cptPrefScoreScale + skillFactorScale + seniorityFactorScale

			//Player is happier if
			//1. they play in an inning
			//2. they play a position they really want
			//3. they have played fewer innings already
			playerScoreThisInn := (playingBonus*playBonusScale/sum +
				prefScore*prefScoreScale/sum +
				cptPrefScore*cptPrefScoreScale/sum +
				skillFactor*skillFactorScale/sum +
				seniorityFactor*seniorityFactorScale/sum)

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
	cptPrefLookup := make([][]float64, game.Roster.NumPlayers())
	skillLookup := make([]float64, game.Roster.NumPlayers())
	seniorityLookup := make([]float64, game.Roster.NumPlayers())

	genderFemaleLookup := make([]bool, game.Roster.NumPlayers())
	malePlayerIdxList := make([]int, 0)
	femalePlayerIdxList := make([]int, 0)
	playerIdxList := make([]int, 0)

	i := 0
	for _, player := range game.Roster.Players {
		playerLookup[i] = player

		prefLookup[i] = make([]float64, NumFieldPositions)
		cptPrefLookup[i] = make([]float64, NumFieldPositions)

		prefNorm := player.normalizePrefs()
		for pos, normMag := range prefNorm {

			posIdx := position2PosIdx(pos)
			prefLookup[i][posIdx] = normMag

		}

		// player.normalizeCptPrefs()
		// for pos, cptPrefMag := range player.CptPrefNorm {
		// 	posIdx := position2PosIdx(pos)
		// 	cptPrefLookup[i][posIdx] = cptPrefMag
		// }

		// skillLookup[i] = player.Skill
		// seniorityLookup[i] = player.Seniority

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

	numInnings := game.NumInnings()
	scratchGame := NewScratchGame(numInnings)

	//Fill the game with stuff:

	playerTakenTable := make([][]int, numInnings)
	for i := 0; i < numInnings; i++ {
		playerTakenTable[i] = make([]int, game.Roster.NumPlayers())
		//Initialize to all players in untaken positions
		for j := range playerTakenTable[i] {
			playerTakenTable[i][j] = unfilled
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

	oldScore := scratchGame.ScoreGame(
		prefLookup,
		cptPrefLookup,
		playerTakenTable,
		skillLookup,
		seniorityLookup,
	)
	callConvergeAt := 1000000

	convCount := 0
	for {

		pickRandInningNum := rand.Intn(numInnings)
		scratchInning := scratchGame.GetInning(pickRandInningNum)

		pickNewPlayer := pickRandPlayerIdx(playerIdxList)
		pickRandPosition := pickRandFilledPosIdx(scratchInning)

		swapPlayer := scratchInning[pickRandPosition]
		if swapPlayer == unfilled {
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
		if newPlayerPrevPosIdx != unfilled {
			scratchInning[newPlayerPrevPosIdx] = swapPlayer
		}
		//Update player taken table with swapped player's position
		playerTakenTable[pickRandInningNum][swapPlayer] = newPlayerPrevPosIdx

		//Recalculate score
		newScore := scratchGame.ScoreGame(
			prefLookup,
			cptPrefLookup,
			playerTakenTable,
			skillLookup,
			seniorityLookup,
		)

		//increment attempts at convergence
		convCount++

		//If score is smaller or field layout is invalid, swap back
		if newScore < oldScore || scratchGame.verifyGame(genderFemaleLookup) != nil {
			scratchInning[pickRandPosition] = swapPlayer
			playerTakenTable[pickRandInningNum][swapPlayer] = pickRandPosition

			if newPlayerPrevPosIdx != unfilled {
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
	game.minimizeBenchTime()

	return nil
}

func (game *Game) minimizeBenchTime() {
	fmt.Println("Start bench minimization")

	numInnings := game.NumInnings()
	if numInnings <= 1 {
		fmt.Println("only one inning...")
		return
	}

	initVal := game.calcBenchInARow()
	startVal := initVal
outerLoop:
	for {
		fmt.Printf("Starting at %v\n", startVal)
		for i := 0; i < 100000; i++ {
			pick1 := rand.Intn(numInnings)
			pick2 := rand.Intn(numInnings)
			for pick1 == pick2 {
				pick2 = rand.Intn(numInnings)
			}
			game.Innings[pick1], game.Innings[pick2] = game.Innings[pick2], game.Innings[pick1]
			finalVal := game.calcBenchInARow()
			if finalVal < startVal {
				startVal = finalVal
				continue outerLoop
			}
			if finalVal > startVal {
				game.Innings[pick1], game.Innings[pick2] = game.Innings[pick2], game.Innings[pick1]
			}
		}
		fmt.Printf("Ending at %v\n", startVal)
		return
	}
}

func (game *Game) calcBenchInARow() int {
	mostBenchInARow := 0
	for _, player := range game.Roster.Players {
		mostBenchInARowThisPlayer := game.calcBenchInARowByPlayer(player)
		if mostBenchInARowThisPlayer > mostBenchInARow {
			mostBenchInARow = mostBenchInARowThisPlayer
		}
	}
	return mostBenchInARow
}

func (game *Game) calcBenchInARowByPlayer(player *Player) int {
	mostBenchInARowThisPlayer := 0
	runningBenchInARowThisPlayer := 0
	for _, inning := range game.Innings {
		role := inning.FindPlayerPosition(player)
		if role != Bench {
			runningBenchInARowThisPlayer = 0
		} else {
			runningBenchInARowThisPlayer++
			if runningBenchInARowThisPlayer > mostBenchInARowThisPlayer {
				mostBenchInARowThisPlayer = runningBenchInARowThisPlayer
			}
		}
	}
	return mostBenchInARowThisPlayer
}

func (inning *BetterInning) FindPlayerPosition(player *Player) Position {
	for pos, playerAtPos := range inning.PositionsMap {
		if player == playerAtPos {
			return pos
		}
	}

	return Bench
}

func (game *BetterGame) calcBenchInARowForPlayer(player *Player) int {
	mostBenchInARowThisPlayer := 0
	runningBenchInARowThisPlayer := 0
	for _, inning := range game.Innings {
		role := inning.FindPlayerPosition(player)
		if role != Bench {
			runningBenchInARowThisPlayer = 0
		} else {
			runningBenchInARowThisPlayer++
			if runningBenchInARowThisPlayer > mostBenchInARowThisPlayer {
				mostBenchInARowThisPlayer = runningBenchInARowThisPlayer
			}
		}
	}
	return mostBenchInARowThisPlayer
}

func (game *BetterGame) calcBenchInARow() int {
	mostBenchInARow := 0
	for _, playerInfo := range game.Roster {
		player := playerInfo.Player
		mostBenchInARowThisPlayer := game.calcBenchInARowForPlayer(player)
		if mostBenchInARowThisPlayer > mostBenchInARow {
			mostBenchInARow = mostBenchInARowThisPlayer
		}
	}
	return mostBenchInARow
}

func (game *Game) fillFromScratch(scrGame *ScratchGame, playerLookup []*Player) {
	for inningNum := 0; inningNum < scrGame.numInnings; inningNum++ {

		for posIdx := 0; posIdx < NumFieldPositions; posIdx++ {

			playerIdx := scrGame.GetPlayerIdxAtPos(inningNum, posIdx)
			if playerIdx == unfilled {
				continue
			}

			player := playerLookup[playerIdx]
			pos := posIdx2Position(posIdx)

			game.Innings[inningNum].FieldPositions[pos] = player
		}
	}
}

func pickRandFilledPosIdx(inn ScratchInning) int {
	for {
		randPosIdx := rand.Intn(len(inn))
		if inn[randPosIdx] != unfilled {
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

		if playerPosTable[plIdx] == unfilled {
			somethingAvail = true
			break
		}
	}
	if somethingAvail == false {
		return -1, fmt.Errorf("No players are available from this list")
	}

	for {
		pickIdx := pickRandPlayerIdx(playerIdxList)
		if playerPosTable[pickIdx] == unfilled {
			playerPosTable[pickIdx] = fillPos
			return pickIdx, nil
		}
	}
}

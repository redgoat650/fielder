package fielder

import (
	"fmt"
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

		s.WriteString(fmt.Sprintf("Inning %d:\n", inningNum))

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
			s.WriteString(fmt.Sprintf("Inning %d: %s plays ", inningNum, player.FirstName))
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
	prefScaleFactor = float64(0.2)
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

					for pos := range playerPosScores {

						//Initialize by seniority and skill
						scoringMtx[playerInfo][pos] += playerInfo.Seniority
						scoringMtx[playerInfo][pos] += playerInfo.Skill

						//Initialize by player preference
						for prefRank, pref := range playerInfo.Pref {
							if pref == pos {
								// old := scoringMtx[playerIdx][posIdx]
								scoringMtx[playerInfo][pos] += (maxPreferences - float64(prefRank)) * prefScaleFactor
								// fmt.Println(playerInfo.FirstName, "OLD", old, "NEW", scoringMtx[playerIdx][posIdx])
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

							// fmt.Printf("Player %s is a candidate for position %v because score %f >= threshold %f\n", playerInfo.FirstName, posIdx2Position(posIdx), score, initialMax)
						}

					}
				}

				//Each position
				for pos, playerListAtPos := range listOfPlayerIdxsAboveThresholdByPosition {

					//Check all player indexes for this position
					for {

						if len(playerListAtPos) == 0 {
							// fmt.Println("Couldn't find a suitable candidate for", position)
							//Couldn't find a suitable candidate at this threshold level
							break
						}

						if _, ok := inning.FieldPositions[pos]; ok {
							//We've already picked a player for this position
							// fmt.Println("Already picked a player for position", position)
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
							// fmt.Println("Player already assigned for position", position, "...try next")
							continue
						}

						pickedPlayerInfo.Roles[game][inningNum] = pos
						assignedPlayers[pickedPlayerInfo] = true
						inning.FieldPositions[pos] = pickedPlayerInfo

						// fmt.Printf("Picked player %s (%s) to play in position %v for inning %d\n", pickedPlayerInfo.FirstName, pickedPlayerInfo.Gender, position, inningNum)
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

	return nil
}

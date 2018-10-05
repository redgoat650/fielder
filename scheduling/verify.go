package fielder

import "fmt"

//Constraints for verification
const (
	MinGenderCount = 4
)

//GenderError is a helper struct that implements the error interface
//and includes the PlayerGender that flagged the error
type GenderError struct {
	err    error
	gender PlayerGender
}

//Error implements the Error interface for the GenderError
func (err GenderError) Error() string {
	return err.err.Error()
}

//verify is a helper method that verifies the Game scheduling
//to ensure the result is a valid position distribution
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
			err := fmt.Errorf("Invalid gender assignment in this inning. Inning %d females: %d/%d males %d/%d", inningNum, femaleCount, MinGenderCount, maleCount, MinGenderCount)
			gender := MaleGender
			if femaleCount < MinGenderCount {
				gender = FemaleGender
			}
			return GenderError{err: err, gender: gender}
		}

	}

	return nil
}

//checkRoster is a helper function that ensures the provided Roster
//is enough to field a team without forfeiting
func checkRoster(roster *Roster) error {
	female, male := roster.CountGenders()
	if female < MinGenderCount {
		return fmt.Errorf("Not enough females. Must forfeit")
	}
	if male < MinGenderCount {
		return fmt.Errorf("Not enough males. Must forfeit")
	}

	return nil
}

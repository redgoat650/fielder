package fielder

import "fmt"

const (
	MinGenderCount = 4
)

type GenderError struct {
	err    error
	gender PlayerGender
}

func (err GenderError) Error() string {
	return err.err.Error()
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

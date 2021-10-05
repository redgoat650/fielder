package fielder

import (
	"fmt"
	"strings"
)

type PlayerID uint64

//Player is a struct for storing the information and scheduling
//information for each player
type Player struct {
	Name   string
	Gender PlayerGender

	Pref map[Position]int
	// PrefNorm    map[Position]float64
	// CptPref     map[Position]int
	// CptPrefNorm map[Position]float64
	// Seniority   float64
	// Skill       float64
}

//NewPlayer initializes a new Player and returns its pointer
func NewPlayer(name string, gender PlayerGender) *Player {
	return &Player{
		Name:   name,
		Gender: gender,
		Pref:   prefInit(),
	}
	// p := new(Player)
	// p.Name = name
	// p.Gender = gender

	// p.Pref = make(map[Position]int)
	// // p.CptPref = make(map[Position]int)

	// //Initialize the preferences table
	// for _, pos := range fieldPosList {
	// 	p.Pref[pos] = 0.0
	// 	// p.CptPref[pos] = 0.0
	// }

	// // p.PrefNorm = make(map[Position]float64)
	// // p.CptPrefNorm = make(map[Position]float64)

	// return p
}

func prefInit() map[Position]int {
	pref := make(map[Position]int)
	for _, pos := range fieldPosList {
		pref[pos] = 0
	}

	return pref
}

func (player *Player) normalizePrefs() map[Position]float64 {
	max := 0.0
	min := 0.0
	minSet := false
	for _, prefStrength := range player.Pref {

		strFloat := float64(prefStrength)
		if strFloat > max {
			max = strFloat
		}
		if !minSet || strFloat < min {
			minSet = true
			min = strFloat
		}
	}

	prefNorm := make(map[Position]float64)
	for pos, val := range player.Pref {

		if max == min {
			prefNorm[pos] = 1.0
			continue
		}

		prefNorm[pos] = (float64(val) - min) / (float64(max) - min)

		//Checks on output
		if prefNorm[pos] > 1.0 {
			panic("Normalization error: Preference overflows expected max")
		}
		if prefNorm[pos] < 0.0 {
			panic("Normalization error: Preference underflows expected min")
		}
	}

	return prefNorm
}

// func (player *Player) normalizeCptPrefs() {
// 	max := 0.0
// 	min := 0.0
// 	minSet := false
// 	for _, prefStrength := range player.CptPref {

// 		strFloat := float64(prefStrength)
// 		if strFloat > max {
// 			max = strFloat
// 		}
// 		if !minSet || strFloat < min {
// 			minSet = true
// 			min = strFloat
// 		}
// 	}

// 	for pos, val := range player.CptPref {

// 		if max == min {
// 			player.CptPrefNorm[pos] = 1.0
// 			continue
// 		}

// 		player.CptPrefNorm[pos] = (float64(val) - min) / (float64(max) - min)

// 		//Checks on output
// 		if player.CptPrefNorm[pos] > 1.0 {
// 			panic("Normalization error: Preference overflows expected max")
// 		}
// 		if player.CptPrefNorm[pos] < 0.0 {
// 			panic("Normalization error: Preference underflows expected min")
// 		}
// 	}

// }

//IsFemale is a helper method for Player that returns whether
//the player is female
func (player *Player) IsFemale() bool {
	return player.Gender == FemaleGender
}

//IsAttending is a helper method for Player that returns whether
//the player in question is planning to attend the provided Game
// func (player *Player) IsAttending(game *Game) bool {
// 	return player.Attendance[game]
// }

//PlayerGender is a type that describes the gender of a player
type PlayerGender string

//List of PlayerGenders
const (
	FemaleGender  PlayerGender = "female"
	MaleGender    PlayerGender = "male"
	InvalidGender PlayerGender = "invalid"
)

var femaleGenderStrs = []string{"f", "female", "girl", "woman"}
var maleGenderStrs = []string{"m", "male", "boy", "man"}

func isFemaleGender(s string) bool {
	return isGender(s, femaleGenderStrs)
}

func isMaleGender(s string) bool {
	return isGender(s, maleGenderStrs)
}

func isGender(s string, l []string) bool {
	for _, cmp := range l {
		if s == cmp {
			return true
		}
	}

	return false
}

// ParseGenderString parses a string and returns a PlayerGender.
// It is the complementary operation to PlayerGender.String()
func ParseGenderString(genderStr string) (PlayerGender, error) {
	lowerGenderStr := strings.ToLower(genderStr)

	if isFemaleGender(lowerGenderStr) {
		return FemaleGender, nil
	}

	if isMaleGender(lowerGenderStr) {
		return MaleGender, nil
	}

	return InvalidGender, fmt.Errorf("unable to parse gender string %v", genderStr)
}

func (player Player) String() string {
	str := new(strings.Builder)

	str.WriteString(fmt.Sprintf("%s, %s\n", player.Name, player.Gender))

	prefNorm := player.normalizePrefs()

	for pref, val := range player.Pref {
		str.WriteString(fmt.Sprintf("%s %d (%v)\n", pref, val, prefNorm[pref]))
	}

	str.WriteString("\n")

	// for pref, val := range player.CptPref {
	// 	str.WriteString(fmt.Sprintf("%s %d (%v)\n", pref, val, player.CptPrefNorm[pref]))
	// }

	// str.WriteString("\n")

	return str.String()

}

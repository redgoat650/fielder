package fielder

import (
	"fmt"
)

//Position is a type that describes a field position
type Position int

//List of valid field positions and the bench positions
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

//String is a position method that satisfies the Stringer interface and
//returns a string describing the Position receiver.
func (pos Position) String() string {
	switch pos {
	case Bench:
		return "Bench"
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
		return "Left Shortstop"
	case RShort:
		return "Right Shortstop"
	case LField:
		return "Left Field"
	case LCenter:
		return "Left Center"
	case RCenter:
		return "Right Center"
	case RField:
		return "Right Field"
	default:
		return "Not a position"
	}
}

// ParsePositionStr parses a string and returns a Position.
// It is the complimentary operation to Position.String()
func ParsePositionStr(posStr string) (Position, error) {
	for _, checkPos := range []Position{
		Bench,
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
	} {
		if checkPos.String() == posStr {
			return checkPos, nil
		}
	}
	return Bench, fmt.Errorf("Can't parse position string %v", posStr)
}

//posIdx2Position converts a position index to a Position value,
//assuming the position index is originating from a zero-indexed
//list of valid field positions
func posIdx2Position(posIdx int) Position {

	if posIdx > NumFieldPositions {
		panic(fmt.Sprintf("Invalid position index: %d", posIdx))
	}

	//Skip the bench position, return the next position
	return Position(posIdx + 1)
}

//position2PosIdx returns the position index from an input Position
func position2PosIdx(pos Position) int {
	return int(pos) - 1
}

// ParsePositionGroupString parses a string for a position group
func ParsePositionGroupString(posStr string) ([]Position, error) {
	singlePos, err := ParsePositionStr(posStr)
	if err == nil {
		return []Position{singlePos}, nil
	}

	for _, checkPosGroup := range []PositionGroup{
		Outfield,
		Infield,
		AnyBase,
		AnyShort,
		NoBase,
		NoShort,
		NotPitcher,
		NotCatcher,
		NotPitcherCatcher,
		LiterallyAnything,
		TierOnePositions,
		TierTwoPositions,
		TierThreePositions,
		TierFourPositions,
		TierFivePositions,
	} {
		if checkPosGroup.String() == posStr {
			ps, ok := posGroup2Positions[checkPosGroup]
			if !ok {
				panic("Could not find position group definition")
			}
			return ps, nil
		}
	}
	return nil, fmt.Errorf("Could not parse position group string %q", posStr)
}

// PositionGroup enumerates a named group of positions
type PositionGroup int

// Position group enumeration
const (
	InvalidPosGroup PositionGroup = iota
	Outfield
	Infield
	AnyBase
	AnyShort
	NoBase
	NoShort
	NotPitcher
	NotCatcher
	NotPitcherCatcher
	LiterallyAnything
	TierOnePositions   // RF, LC
	TierTwoPositions   // 2B, RS
	TierThreePositions // LS, RC, LF
	TierFourPositions  // 1B, 3B
	TierFivePositions  // Pitcher, Catcher
)

var posGroup2Positions = map[PositionGroup][]Position{
	Outfield:           []Position{RField, RCenter, LCenter, LField},
	Infield:            []Position{First, RShort, Second, LShort, Third},
	AnyBase:            []Position{First, Second, Third},
	NoBase:             []Position{RShort, LShort, RField, RCenter, LCenter, LField},
	AnyShort:           []Position{RShort, LShort},
	NoShort:            []Position{First, Second, Third, RField, RCenter, LCenter, LField},
	NotPitcher:         []Position{Catcher, First, Second, Third, LShort, RShort, LField, LCenter, RCenter, RField},
	NotCatcher:         []Position{Pitcher, First, Second, Third, LShort, RShort, LField, LCenter, RCenter, RField},
	NotPitcherCatcher:  []Position{First, Second, Third, LShort, RShort, LField, LCenter, RCenter, RField},
	LiterallyAnything:  []Position{Pitcher, Catcher, First, Second, Third, LShort, RShort, LField, LCenter, RCenter, RField},
	TierOnePositions:   []Position{RField, LCenter},
	TierTwoPositions:   []Position{RShort, Second},
	TierThreePositions: []Position{LShort, RCenter, LField},
	TierFourPositions:  []Position{First, Third},
	TierFivePositions:  []Position{Pitcher, Catcher},
}

// String stringifies the position group
func (pg PositionGroup) String() string {
	switch pg {
	case Outfield:
		return "Outfield"
	case Infield:
		return "Infield"
	case AnyBase:
		return "AnyBase"
	case AnyShort:
		return "Any Short"
	case NoBase:
		return "No Base"
	case NoShort:
		return "No Short"
	case NotPitcher:
		return "Not Pitcher"
	case NotCatcher:
		return "Not Catcher"
	case NotPitcherCatcher:
		return "Not Pitcher Catcher"
	case LiterallyAnything:
		return "Literally Anything"
	case TierOnePositions: // RF, LC
		return "Tier One Positions"
	case TierTwoPositions: // 2B, RS
		return "Tier Two Positions"
	case TierThreePositions: // LS, RC, LF
		return "Tier Three Positions"
	case TierFourPositions: // 1B, 3B
		return "Tier Four Positions"
	case TierFivePositions: // Pitcher, Catcher
		return "Tier Five Positions"
	default:
		return "Invalid Position group"
	}
}

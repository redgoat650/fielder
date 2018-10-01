package fielder

import "fmt"

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

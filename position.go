package fielder

import "fmt"

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

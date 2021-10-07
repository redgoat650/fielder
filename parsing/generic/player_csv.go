package generic

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/redgoat650/fielder/parsing/csvhelper"
	fielder "github.com/redgoat650/fielder/scheduling"
)

func ParsePlayerListFromFile(filename string) ([]*fielder.Player, error) {
	return parsePlayerListFromFile(filename)
}

func parsePlayerListFromFile(filename string) ([]*fielder.Player, error) {
	playerData, err := csvhelper.ParseCSVFile(filename)
	if err != nil {
		return nil, err
	}

	return parsePlayerData(playerData)
}

func parsePlayerData(playerDataTable [][]string) ([]*fielder.Player, error) {
	if len(playerDataTable) == 0 {
		return nil, errors.New("unable to parse player data")
	}

	headerRow := playerDataTable[0]

	colMap, err := readPlayerTableHeaderRow(headerRow)
	if err != nil {
		return nil, err
	}

	playerData := playerDataTable[1:]

	return makePlayerList(playerData, colMap)
}

const (
	nameColKey   = "nameCol"
	genderColKey = "genderCol"
)

func makePlayerList(playerData [][]string, colMap map[string]int) (playerList []*fielder.Player, err error) {
	for _, rowData := range playerData {
		pl, err := parseRowToPlayer(rowData, colMap)
		if err != nil {
			return nil, err
		}

		playerList = append(playerList, pl)
	}

	return playerList, nil
}

func parseRowToPlayer(rowData []string, colMap map[string]int) (*fielder.Player, error) {
	name := rowData[colMap[nameColKey]]
	genderStr := rowData[colMap[genderColKey]]
	gender, err := fielder.ParseGenderString(genderStr)
	if err != nil {
		fmt.Printf("Unable to parse gender string: %q\n", genderStr)
		return nil, err
	}

	player := fielder.NewPlayer(name, gender)

	err = parseRowToPlayerPrefs(player, rowData, colMap)
	if err != nil {
		return nil, err
	}

	return player, nil
}

func parseRowToPlayerPrefs(player *fielder.Player, rowData []string, colMap map[string]int) error {
	for colNameKey, colIdx := range colMap {
		switch colNameKey {
		case nameColKey, genderColKey:
			//Ignore
			continue
		}

		// colNameKey should be a parseable position
		pos, err := fielder.ParsePositionStr(colNameKey)
		if err != nil {
			fmt.Printf("Could not parse col name key: %q", colNameKey)
			return err
		}

		// Actual value in the column should be an integer
		prefValueStr := rowData[colIdx]
		prefValue, err := strconv.Atoi(prefValueStr)
		if err != nil {
			fmt.Printf("Could not parse %q's %s preference as integer: %q", player.Name, pos, prefValueStr)
			return err
		}

		player.Pref[pos] = prefValue
	}

	return nil
}

const (
	nameHeaderString   = "Name"
	genderHeaderString = "Gender"
)

func readPlayerTableHeaderRow(headerRow []string) (map[string]int, error) {
	expHeaderRowWidth := len(fielder.PositionList) + 2 // Name + Gender
	if len(headerRow) != expHeaderRowWidth {
		fmt.Println("header row len", len(headerRow), "exp", expHeaderRowWidth)
		return nil, errors.New("header row is unexpected length")
	}

	colMap := make(map[string]int)

	foundPosTable := make(map[fielder.Position]struct{})
	for _, pos := range fielder.PositionList {
		foundPosTable[pos] = struct{}{}
	}

	for colIdx, colHeader := range headerRow {
		switch {
		case colHeader == nameHeaderString:
			colMap[nameColKey] = colIdx
		case colHeader == genderHeaderString:
			colMap[genderColKey] = colIdx
		default:
			pos, err := fielder.ParsePositionStr(colHeader)
			if err != nil {
				fmt.Println("Header entry", colHeader)
				return nil, errors.New("unexpected header entry in player table")
			}

			// Mark position as found
			delete(foundPosTable, pos)

			colMap[pos.String()] = colIdx
		}
	}

	if len(foundPosTable) != 0 {
		fmt.Println("positions without columns in table", foundPosTable)
		return nil, errors.New("table in file lacked columns for some positions")
	}

	return colMap, nil
}

func CreateSamplePlayerFile(filename string, numPlayers int) error {
	playerData := createSamplePlayerPlayerData(numPlayers)

	return csvhelper.WriteCSVFile(filename, playerData)
}

func WriteSamplePlayerTable(w io.Writer, numPlayers int) error {
	playerData := createSamplePlayerPlayerData(numPlayers)

	return csvhelper.WriteCSV(w, playerData)
}

func createSamplePlayerPlayerData(numPlayers int) [][]string {
	header := makePlayerTableHeader()

	return append(
		[][]string{header},
		makeSamplePlayerRows(numPlayers)...,
	)
}

func makePlayerTableHeader() []string {
	header := []string{
		nameHeaderString,
		genderHeaderString,
	}

	for _, pos := range fielder.PositionList {
		header = append(header, pos.String())
	}

	return header
}

const (
	namePlaceholder   = "<ENTER PLAYER NAME>"
	genderPlaceholder = "<ENTER PLAYER GENDER>"
)

func makeSamplePlayerRows(numPlayers int) (ret [][]string) {
	for i := 0; i < numPlayers; i++ {
		ret = append(ret, makeSamplePlayerRow())
	}

	return ret
}

func makeSamplePlayerRow() []string {
	ret := []string{
		namePlaceholder,
		genderPlaceholder,
	}

	for range fielder.PositionList {
		ret = append(ret, strconv.Itoa(0))
	}

	return ret
}

package generic

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/redgoat650/fielder/parsing/csvhelper"
	fielder "github.com/redgoat650/fielder/scheduling"
)

func ParseCaptainListFromFile(filename string) (fielder.BetterRoster, error) {
	return parseCaptainListFromFile(filename)
}

func parseCaptainListFromFile(filename string) (fielder.BetterRoster, error) {
	cptData, err := csvhelper.ParseCSVFile(filename)
	if err != nil {
		return nil, err
	}

	return parseCaptainData(cptData)
}

func parseCaptainData(cptDataTable [][]string) (fielder.BetterRoster, error) {
	if len(cptDataTable) == 0 {
		return nil, errors.New("unable to parse player data")
	}

	headerRow := cptDataTable[0]

	colMap, err := readCaptainTableHeaderRow(headerRow)
	if err != nil {
		return nil, err
	}

	playerData := cptDataTable[1:]

	return makeCaptainList(playerData, colMap)
}

func makeCaptainList(playerData [][]string, colMap map[string]int) (roster fielder.BetterRoster, err error) {
	for _, rowData := range playerData {
		sp, err := parseRowToScoringParams(rowData, colMap)
		if err != nil {
			return nil, err
		}

		roster = append(roster, sp)
	}

	return roster, nil
}

const (
	skillColKey     = "skillCol"
	seniorityColKey = "seniorityCol"
)

func parseRowToScoringParams(rowData []string, colMap map[string]int) (fielder.ScoringParams, error) {
	name := rowData[colMap[nameColKey]]
	skillStr := rowData[colMap[skillColKey]]
	seniorityStr := rowData[colMap[seniorityColKey]]

	skillVal, err := strconv.ParseFloat(skillStr, 64)
	if err != nil {
		return fielder.ScoringParams{}, err
	}

	seniorityVal, err := strconv.ParseFloat(seniorityStr, 64)
	if err != nil {
		return fielder.ScoringParams{}, err
	}

	cptPrefs, err := parseRowToCaptainPrefs(rowData, colMap)
	if err != nil {
		return fielder.ScoringParams{}, err
	}

	return fielder.ScoringParams{
		Player:    &fielder.Player{Name: name},
		CptPref:   cptPrefs,
		Skill:     skillVal,
		Seniority: seniorityVal,
	}, nil
}

func parseRowToCaptainPrefs(rowData []string, colMap map[string]int) (map[fielder.Position]int, error) {
	cptPref := make(map[fielder.Position]int)

	for colNameKey, colIdx := range colMap {
		switch colNameKey {
		case nameColKey, skillColKey, seniorityColKey:
			//Ignore
			continue
		}

		// colNameKey should be a parseable position
		pos, err := fielder.ParsePositionStr(colNameKey)
		if err != nil {
			fmt.Printf("Could not parse col name key: %q", colNameKey)
			return nil, err
		}

		// Actual value in the column should be an integer
		prefValueStr := rowData[colIdx]
		prefValue, err := strconv.Atoi(prefValueStr)
		if err != nil {
			fmt.Printf("Could not parse %s captain preference as integer: %q", pos, prefValueStr)
			return nil, err
		}

		cptPref[pos] = prefValue
	}

	return cptPref, nil
}

const (
	skillHeaderString     = "Skill"
	seniorityHeaderString = "Seniority"
)

func readCaptainTableHeaderRow(headerRow []string) (map[string]int, error) {
	expHeaderRowWidth := len(fielder.PositionList) + 3 // Name + Skill + Seniority

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
		case colHeader == skillHeaderString:
			colMap[skillColKey] = colIdx
		case colHeader == seniorityHeaderString:
			colMap[seniorityColKey] = colIdx
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

func CreateSampleCaptainFile(filename string, numPlayers int) error {
	playerData := createSampleCaptainPlaceholderData(numPlayers)

	return csvhelper.WriteCSVFile(filename, playerData)
}

func WriteSampleCaptainTable(w io.Writer, numPlayers int) error {
	playerData := createSampleCaptainPlaceholderData(numPlayers)

	return csvhelper.WriteCSV(w, playerData)
}

func createSampleCaptainPlaceholderData(numPlayers int) [][]string {
	pl := []*fielder.Player{}
	for i := 0; i < numPlayers; i++ {
		pl = append(pl, &fielder.Player{Name: namePlaceholder})
	}

	return createSampleCaptainData(pl)
}

func WriteSampleCaptainFileFromPlayerFile(w io.Writer, ptFilepath string) error {
	playerList, err := parsePlayerListFromFile(ptFilepath)
	if err != nil {
		return err
	}

	cptData := createSampleCaptainData(playerList)

	return csvhelper.WriteCSV(w, cptData)
}

func createSampleCaptainData(playerList []*fielder.Player) [][]string {
	header := makeCptTableHeader()

	return append(
		[][]string{header},
		makeSampleCaptainRows(playerList)...,
	)
}

func makeCptTableHeader() []string {
	header := []string{
		nameHeaderString,
		skillHeaderString,
		seniorityHeaderString,
	}

	for _, pos := range fielder.PositionList {
		header = append(header, pos.String())
	}

	return header
}

func makeSampleCaptainRows(playerList []*fielder.Player) (ret [][]string) {
	for _, p := range playerList {
		ret = append(ret, makeSampleCaptainRow(p.Name))
	}

	return ret
}

const (
	skillPlaceholder     = 0.0
	seniorityPlaceholder = 0.0
)

var (
	skillPlaceholderStr     = fmt.Sprintf("%f", skillPlaceholder)
	seniorityPlaceholderStr = fmt.Sprintf("%f", seniorityPlaceholder)
)

func makeSampleCaptainRow(name string) []string {
	ret := []string{
		name,
		skillPlaceholderStr,
		seniorityPlaceholderStr,
	}

	for range fielder.PositionList {
		ret = append(ret, strconv.Itoa(0))
	}

	return ret
}

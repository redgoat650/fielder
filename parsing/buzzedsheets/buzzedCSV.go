package buzzedsheets

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	fielder "github.com/redgoat650/fielder/scheduling"
)

const (
	numInnings = 7
)

// ParseBuzzedSheets takes in a legacy buzzed style sheet in CSV format
// and outputs a Game ready for scheduling on the specified day.
func ParseBuzzedSheets(scheduleCSVPath, preferenceCSVPath, captCfgCSV, gameDate string) (*fielder.Game, error) {
	schedule, err := parseCSVFile(scheduleCSVPath)
	if err != nil {
		return nil, err
	}

	nameList, err := getNameList(schedule)
	if err != nil {
		return nil, err
	}

	var pref [][]string
	if preferenceCSVPath != "" {
		pref, err = parseCSVFile(preferenceCSVPath)
		if err != nil {
			return nil, err
		}
	}

	var cptPrefParsed [][]string
	if captCfgCSV != "" {
		err := createCptPrefIfNotExist(captCfgCSV, nameList)
		if err != nil {
			return nil, err
		}
		cptPrefParsed, err = parseCSVFile(captCfgCSV)
		if err != nil {
			return nil, err
		}
	}

	numNames := len(nameList)
	genderList, err := getGenderList(schedule, numNames)
	if err != nil {
		return nil, err
	}
	emailList, err := getEmailList(schedule, numNames)
	if err != nil {
		return nil, err
	}
	attendenceList, err := getAttendenceList(schedule, gameDate, numNames)
	if err != nil {
		return nil, err
	}
	attendenceMap := getAttendenceMap(nameList, attendenceList)

	gameRoster := fielder.NewRoster()
	for i, name := range nameList {
		if attending := attendenceMap[name]; !attending {
			continue
		}
		splitName := strings.Split(name, " ")
		first, last := splitName[0], splitName[1]
		gender, err := fielder.ParseGenderString(genderList[i])
		if err != nil {
			return nil, err
		}
		pl := fielder.NewPlayer(first, last, gender)
		pl.Email = emailList[i]

		if pref != nil {
			err := fillPreferences(pl.Pref, pref, name)
			if err != nil {
				return nil, err
			}
		}

		if cptPrefParsed != nil {
			err := fillPreferences(pl.CptPref, cptPrefParsed, name)
			if err != nil {
				return nil, err
			}
		}
		gameRoster.AddPlayer(pl)
	}

	game := fielder.NewGame(numInnings, 0)
	game.SetRoster(gameRoster)
	return game, nil
}

func parseCSVFile(filepath string) ([][]string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(f)
	parsed, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

const (
	schedHeaderOffset = 1
	prefHeaderOffset  = 0
	numPrefSelections = 3
	nameHeaderStr     = "Name"
	genderHeaderStr   = "Gender"
	emailHeaderStr    = "Email"
	attendingSubStr   = "Yes"
)

func createCptPrefIfNotExist(filepath string, nameList []string) error {
	_, err := os.Stat(filepath)
	if err == nil {
		return nil
	}

	f, err := os.Create(filepath)
	if err != nil {
		return err
	}

	err = writePrefFileHeader(f)
	if err != nil {
		return err
	}

	for _, name := range nameList {
		f.Write([]byte(fmt.Sprintf("%s,,,\n", name)))
	}

	return f.Close()
}

func writePrefFileHeader(f *os.File) error {
	colHeadings := []string{
		"Name",
	}
	for i := 0; i < numPrefSelections; i++ {
		colHeadings = append(colHeadings, getPrefPosHeader(i))
	}
	_, err := f.Write([]byte(fmt.Sprintf("%v\n", strings.Join(colHeadings, ","))))
	if err != nil {
		return err
	}
	return nil
}

func getColumnList(sheet [][]string, colTitle string, minNumRows int, headerRowOffset int) ([]string, error) {
	headerRow := getHeaderRow(sheet, headerRowOffset)
	colIdx, err := findColumnIdx(headerRow, colTitle)
	if err != nil {
		return nil, err
	}
	valList := make([]string, 0)
	for _, v := range sheet[headerRowOffset+1:] {
		colVal := v[colIdx]
		if colVal == "" {
			if len(valList) < minNumRows {
				return nil, fmt.Errorf("%v column is not fully filled out", colTitle)
			}
			return valList, nil
		}
		valList = append(valList, colVal)
	}
	return valList, nil
}

func findColumnIdx(headerRow []string, colTitle string) (int, error) {
	for i, v := range headerRow {
		if strings.TrimSpace(v) == colTitle {
			return i, nil
		}
	}
	return 0, fmt.Errorf("Unable to find column title %v in header row %v", colTitle, headerRow)
}

func getNameList(schedule [][]string) ([]string, error) {
	return getColumnList(schedule, nameHeaderStr, 0, schedHeaderOffset)
}

func getGenderList(schedule [][]string, numNames int) ([]string, error) {
	return getColumnList(schedule, genderHeaderStr, numNames, schedHeaderOffset)
}

func getEmailList(schedule [][]string, numNames int) ([]string, error) {
	return getColumnList(schedule, emailHeaderStr, numNames, schedHeaderOffset)
}

func getHeaderRow(schedule [][]string, headerRowOffset int) []string {
	return schedule[headerRowOffset]
}

func getAttendenceList(schedule [][]string, date string, numNames int) ([]bool, error) {
	return getGameAttendenceList(schedule, date, numNames)
}

func getGameAttendenceList(schedule [][]string, date string, numPlayers int) ([]bool, error) {
	headerRow := getHeaderRow(schedule, schedHeaderOffset)
	dateColIdx, err := findColumnIdx(headerRow, date)
	if err != nil {
		return nil, err
	}
	attendingList := make([]bool, numPlayers)
	for i, v := range schedule[schedHeaderOffset+1:] {
		if i >= len(attendingList) {
			return attendingList, nil
		}
		attendingStr := v[dateColIdx]
		attendingList[i] = (strings.Contains(attendingStr, attendingSubStr) || attendingStr == "")
	}
	return attendingList, nil
}

func getAttendenceMap(nameList []string, attendenceList []bool) map[string]bool {
	if len(nameList) != len(attendenceList) {
		panic("Name list len does not equal attendence list len")
	}
	ret := make(map[string]bool)
	for i := range nameList {
		ret[nameList[i]] = attendenceList[i]
	}
	return ret
}

func getPreferences(pref [][]string, name string) (map[fielder.Position]int, error) {
	prefList, err := getPrefListByName(pref, name)
	if err != nil {
		return nil, err
	}

	ret := make(map[fielder.Position]int)
	for prefOrder, prefPos := range prefList {
		posGroup, err := fielder.ParsePositionGroupString(prefPos)
		if err != nil {
			// Do nothing, whatever
			continue
		}
		for _, pos := range posGroup {
			ret[pos] += getScore(prefOrder)
		}
	}
	return ret, nil
}

func fillPreferences(in map[fielder.Position]int, pref [][]string, name string) error {
	prefList, err := getPrefListByName(pref, name)
	if err != nil {
		return err
	}

	for prefOrder, prefPos := range prefList {
		posGroup, err := fielder.ParsePositionGroupString(prefPos)
		if err != nil {
			// Do nothing, whatever
			continue
		}
		for _, pos := range posGroup {
			in[pos] += getScore(prefOrder)
		}
	}
	return nil
}

func getScore(idx int) int {
	return numPrefSelections - idx
}

func getPrefListByName(pref [][]string, name string) ([]string, error) {
	nameList, err := getColumnList(pref, nameHeaderStr, 0, prefHeaderOffset)
	if err != nil {
		return nil, err
	}
	for i, checkName := range nameList {
		rowNum := prefHeaderOffset + 1 + i
		if name == checkName {
			return getPrefStrsByIdx(pref, rowNum)
		}
	}
	return nil, fmt.Errorf("Name %v not found in preferences list", name)
}

func getPrefStrsByIdx(pref [][]string, rowIdx int) ([]string, error) {
	ret := make([]string, numPrefSelections)
	headerRow := getHeaderRow(pref, prefHeaderOffset)
	for i := 0; i < len(ret); i++ {
		prefIdx, err := findColumnIdx(headerRow, getPrefPosHeader(i))
		if err != nil {
			return nil, err
		}
		ret[i] = pref[rowIdx][prefIdx]
	}
	return ret, nil
}

func getPrefPosHeader(idx int) string {
	return fmt.Sprintf("Preferred Position %v", idx+1)
}

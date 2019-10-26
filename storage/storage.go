package storage

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func GetTeamsDirectory() string {
	return filepath.Join("teams")
}

func GetTeamDirectory(teamName string) string {
	return filepath.Join(GetTeamsDirectory(), teamName)
}

func GetSeasonDirectory(teamName, season string) string {
	return filepath.Join(GetTeamDirectory(teamName), "seasons", season)
}

func GetSeasonCfgFilePath(teamName, season string) string {
	return filepath.Join(GetSeasonDirectory(teamName, season), "seasonCfg.yaml")
}

func GetPlayersDirectory(teamName, season string) string {
	return filepath.Join(GetSeasonDirectory(teamName, season), "players")
}

func GetGamesDirectory(teamName, season string) string {
	return filepath.Join(GetSeasonDirectory(teamName, season), "games")
}

func GetGameDirectoryByIdx(teamName, season string, gameIdx int) string {
	return filepath.Join(GetGamesDirectory(teamName, season), fmt.Sprintf("game%v", gameIdx))
}

func GetLineupFilePath(teamName, season string, gameIdx int, format string) string {
	return filepath.Join(GetGameDirectoryByIdx(teamName, season, gameIdx), fmt.Sprintf("lineup.%v", format))
}

func GetFieldPositionsFilePath(teamName, season string, gameIdx int, format string) string {
	return filepath.Join(GetGameDirectoryByIdx(teamName, season, gameIdx), fmt.Sprintf("field_positions.%v", format))
}

func GetRosterFilePath(teamName, season string, gameIdx int, format string) string {
	return filepath.Join(GetGameDirectoryByIdx(teamName, season, gameIdx), fmt.Sprintf("roster.%v", format))
}

func WriteToFile(r io.Reader, path string) error {

	dir, file := filepath.Split(path)

	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}

	_, copyErr := io.Copy(f, r)
	if copyErr != nil {
		return copyErr
	}

	return nil
}

func SaveGob(path string, thing interface{}) error {
	dir, _ := filepath.Split(path)
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)

	encErr := gob.NewEncoder(buf).Encode(thing)
	if encErr != nil {
		return encErr
	}

	err = ioutil.WriteFile(path, buf.Bytes(), 0777)
	if err != nil {
		return err
	}

	return nil
}

func LoadGob(path string) (*gob.Decoder, error) {

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(b)

	return gob.NewDecoder(buf), nil
}

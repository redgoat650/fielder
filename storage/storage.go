package storage

import (
	"fmt"
	"io"
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

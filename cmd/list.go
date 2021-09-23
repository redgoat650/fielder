/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:`,
	RunE: teamListRunFunc,
}

func teamListRunFunc(cmd *cobra.Command, args []string) error {
	fmt.Println("list called")

	return renderTeamNamesFromDir()
}

func renderTeamNamesFromDir() error {
	names, err := readTeamNamesFromDir()
	if err != nil {
		return err
	}

	fmt.Println("Available team names:")
	for _, name := range names {
		fmt.Println("-", name)
	}

	return nil
}

func readTeamNamesFromDir() ([]string, error) {
	teamsDir := getTeamsDir()

	dirEntries, err := os.ReadDir(teamsDir)
	if err != nil {
		return nil, err
	}

	ret := []string{}
	for _, entry := range dirEntries {
		fileName := entry.Name()
		teamName := strings.TrimSuffix(fileName, ".json")

		var selectedTeamMarker = ""
		if teamName == viper.Get(selectedTeamConfigKey) {
			selectedTeamMarker = " *"
		}

		ret = append(ret, teamName+selectedTeamMarker)
	}

	return ret, nil
}

func init() {
	teamCmd.AddCommand(listCmd)
}

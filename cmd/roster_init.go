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
	"errors"
	"fmt"
	"strings"

	fielder "github.com/redgoat650/fielder/scheduling"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rosterInitCmd represents the init command
var rosterInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a roster",
	Long: `Initialize a roster with a list of players.
	
	Comma separated list of names: <first> <last> <gender>`,
	RunE: rosterInitRunFunc,
}

func rosterInitRunFunc(cmd *cobra.Command, args []string) error {
	fmt.Println("init called")

	if len(args) != 1 {
		return errors.New("expected one list of comma delimited names and genders")
	}

	rosterInput := args[0]

	players := strings.Split(rosterInput, ",")

	selectedSeason := viper.GetString(selectedSeasonConfigKey)
	if selectedSeason == "" {
		return errors.New("no season selected")
	}

	season := gTeam.SeasonList[selectedSeason]
	if season.Roster == nil {
		season.Roster = fielder.NewRoster()
		gTeam.SeasonList[selectedSeason] = season
	}

	for _, playerStr := range players {
		p := strings.Split(playerStr, " ")

		if len(p) != 2 {
			return errors.New("player was not formatted correctly")
		}

		name := p[0]
		gender, err := fielder.ParseGenderString(p[2])
		if err != nil {
			return err
		}

		gTeam.SeasonList[selectedSeason].Roster.AddPlayer(fielder.NewPlayer(name, gender))
	}

	return nil
}

func init() {
	rosterCmd.AddCommand(rosterInitCmd)
}

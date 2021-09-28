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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	clearTeam bool
)

// teamSwitchCmd represents the switch command
var teamSwitchCmd = &cobra.Command{
	Use:   "switch",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:`,
	RunE: switchRunFunc,
}

func switchRunFunc(cmd *cobra.Command, args []string) error {
	fmt.Println("switch called")

	if clearTeam {
		return clearSelectedTeamCheckArgs(args)
	}

	if len(args) != 1 {
		return errors.New("expecting one argument, team name")
	}

	name := args[0]
	fmt.Printf("Switching to team %q\n", name)

	if !teamExists(name) {
		renderTeamNamesFromDir()
		return errors.New("team not found")
	}

	err := loadTeamByName(name)
	if err != nil {
		return err
	}

	err = viperSetTeam(name)
	if err != nil {
		return err
	}

	return renderTeamNamesFromDir()
}

func clearSelectedTeamCheckArgs(args []string) error {
	if len(args) != 0 {
		return errors.New("expecting zero arguments")
	}

	return clearSelectedTeam()
}

func clearSelectedTeam() error {
	fmt.Println("Clearing selected team")
	viper.Set(selectedTeamConfigKey, "")

	err := viperUpdateOrCreate()
	if err != nil {
		return err
	}

	return renderTeamNamesFromDir()
}

func init() {
	teamCmd.AddCommand(teamSwitchCmd)
	teamSwitchCmd.Aliases = append(teamSwitchCmd.Aliases, "select")
	teamSwitchCmd.Flags().BoolVarP(&clearTeam, "none", "X", false, "Clear the currently selected team")
}

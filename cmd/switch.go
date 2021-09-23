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
)

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:`,
	RunE: switchRunFunc,
}

func switchRunFunc(cmd *cobra.Command, args []string) error {
	fmt.Println("switch called")
	name := cmd.Flag(teamNameFlagName).Value.String()

	if !teamExists(name) {
		return errors.New("team not found")
	}

	err := loadTeamByName(name)
	if err != nil {
		return err
	}

	return viperSetTeam(name)
}

func init() {
	teamCmd.AddCommand(switchCmd)

	switchCmd.Flags().StringP(teamNameFlagName, "n", "", "Team name")
	switchCmd.MarkFlagRequired(teamNameFlagName)
}

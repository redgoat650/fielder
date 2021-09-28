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
	seasonDeleteConfirm bool
)

// seasonDeleteCmd represents the delete command
var seasonDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a season",
	Long:  `Delete a season permanently.`,
	RunE:  seasonDeleteRunCmd,
}

func seasonDeleteRunCmd(cmd *cobra.Command, args []string) error {
	fmt.Println("delete called")

	if len(args) != 1 {
		return errors.New("expecting one argument, team name")
	}

	name := args[0]

	if !seasonNameExists(name) {
		fmt.Printf("Could not find season name: %q\n", name)
		return errors.New("season not found")
	}

	fmt.Printf("Will permanently delete season: %q\n", name)

	if !seasonDeleteConfirm {
		fmt.Println("Rerun with '--delete' flag to confirm deletion.")
		return nil
	}

	return deleteSeasonByName(name)
}

func deleteSeasonByName(name string) error {
	fmt.Printf("Removing season %q\n", name)

	delete(gTeam.SeasonList, name)

	if viper.GetString(selectedSeasonConfigKey) == name {
		return clearSelectedSeason()
	}

	return nil
}

func init() {
	seasonCmd.AddCommand(seasonDeleteCmd)
	seasonDeleteCmd.Flags().BoolVarP(&seasonDeleteConfirm, "delete", "D", false, "Confirm deletion of this season")
}

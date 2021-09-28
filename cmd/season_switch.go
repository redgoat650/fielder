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
	clearSeason bool
)

// seasonSwitchCmd represents the switch command
var seasonSwitchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch seasons",
	Long:  `Switch seasons to another valid season for the active team.`,
	RunE:  seasonSwitchRunFunc,
}

func seasonSwitchRunFunc(cmd *cobra.Command, args []string) error {
	fmt.Println("season switch called")

	if clearSeason {
		return clearSelectedSeasonCheckArgs(args)
	}

	if len(args) != 1 {
		return errors.New("expecting one argument, season name")
	}

	name := args[0]
	fmt.Printf("Switching to season %q\n", name)

	if !seasonNameExists(name) {
		renderListSeasons()
		return errors.New("season not found")
	}

	if name == viper.GetString(selectedSeasonConfigKey) {
		fmt.Printf("Season %q was already selected.\n", name)
		renderListSeasons()
		return nil
	}

	viper.Set(selectedSeasonConfigKey, name)
	if err := viperUpdateOrCreate(); err != nil {
		return err
	}

	renderListSeasons()
	return nil
}

func clearSelectedSeasonCheckArgs(args []string) error {
	if len(args) != 0 {
		return errors.New("expecting zero arguments")
	}

	return clearSelectedSeason()
}

func clearSelectedSeason() error {
	fmt.Println("Clearing selected season")
	viper.Set(selectedSeasonConfigKey, "")

	err := viperUpdateOrCreate()
	if err != nil {
		return err
	}

	renderListSeasons()

	return nil
}

func init() {
	seasonCmd.AddCommand(seasonSwitchCmd)
	seasonSwitchCmd.Aliases = append(seasonSwitchCmd.Aliases, "select")
	seasonSwitchCmd.Flags().BoolVarP(&clearSeason, "none", "X", false, "Clear the currently selected season")
}

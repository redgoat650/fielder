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

// seasonAddCmd represents the add command
var seasonAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a season.",
	Long:  `Add a season for the currently selected team.`,
	RunE:  seasonAddRunFunc,
}

func seasonAddRunFunc(cmd *cobra.Command, args []string) error {
	fmt.Println("season add called")

	if gTeam != nil {
		fmt.Println("Select a team with 'fielder team select <name>' or create a new one with 'fielder team create <name>'")
		return errors.New("no team selected")
	}

	return nil
}

func init() {
	seasonCmd.AddCommand(seasonAddCmd)
}

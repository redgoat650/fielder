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
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	teamDeleteConfirm bool
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:`,
	RunE: teamDeleteRunFunc,
}

func teamDeleteRunFunc(cmd *cobra.Command, args []string) error {
	fmt.Println("delete called")

	if len(args) != 1 {
		return errors.New("expecting one argument, team name")
	}

	name := args[0]

	if !teamExists(name) {
		fmt.Printf("Could not find team: %q\n", name)
		return errors.New("team not found")
	}

	fmt.Printf("Will permanently delete team: %q\n", name)

	if !teamDeleteConfirm {
		fmt.Println("Rerun with '--delete' flag to confirm.")
	}

	return deleteTeam(name)
}

func deleteTeam(teamName string) error {
	filename := getFullTeamFilePath(teamName)

	fmt.Println("filename", filename)

	err := os.Remove(filename)
	if err != nil {
		return err
	}

	if viper.Get(selectedTeamConfigKey) == teamName {
		viper.Set(selectedTeamConfigKey, "")
	}

	return nil
}

func init() {
	teamCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	deleteCmd.Flags().BoolVarP(&teamDeleteConfirm, "delete", "D", false, "Confirm deletion")
}

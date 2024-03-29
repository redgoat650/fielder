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

	fielder "github.com/redgoat650/fielder/scheduling"
)

// teamCreateCmd represents the create command
var teamCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:`,

	RunE: teamCreateRunFunc,
}

func teamCreateRunFunc(cmd *cobra.Command, args []string) error {
	fmt.Println("create called")

	if len(args) != 1 {
		return errors.New("expecting one argument, team name")
	}

	name := args[0]
	fmt.Printf("Creating team %q\n", name)

	if teamExists(name) {
		return errors.New("team already exists")
	}

	return teamCreate(name)
}

func teamCreate(name string) error {
	gTeam = fielder.NewTeam(name)

	return viperSetTeam(name)
}

func viperSetTeam(teamName string) error {
	viper.Set(selectedTeamConfigKey, teamName)

	return viperUpdateOrCreate()
}

func teamExists(teamName string) bool {
	filename := getFullTeamFilePath(teamName)

	_, err := os.Stat(filename)
	return err == nil
}

func viperUpdateOrCreate() error {
	err := viper.WriteConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return viper.SafeWriteConfig()
		}
		return err
	}
	return nil
}

func init() {
	teamCmd.AddCommand(teamCreateCmd)
}

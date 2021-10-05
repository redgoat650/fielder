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

	fielder "github.com/redgoat650/fielder/scheduling"
	"github.com/spf13/cobra"
)

var (
	genderFlag string
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a player",
	Long:  `Create a player.`,
	RunE:  playerCreateRunFunc,
}

func playerCreateRunFunc(cmd *cobra.Command, args []string) error {
	fmt.Println("player create called")

	if len(args) != 1 {
		return errors.New("expecting single name argument")
	}

	name := args[0]

	gender, err := fielder.ParseGenderString(genderFlag)
	if err != nil {
		return err
	}

	player := fielder.NewPlayer(name, gender)

	fmt.Println(player)

	return nil
}

func init() {
	playerCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&genderFlag, "gender", "g", "", "Player gender")
	createCmd.MarkFlagRequired("gender")
}

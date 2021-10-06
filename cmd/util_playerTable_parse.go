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

	"github.com/redgoat650/fielder/parsing/generic"
	"github.com/spf13/cobra"
)

var (
	filename string
)

// utilPlayerTableParseCmd represents the parse command
var utilPlayerTableParseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse a player table",
	Long:  `Attempt to parse a player table file given by file name.`,
	RunE:  utilPlayerTableParseRunFunc,
}

func utilPlayerTableParseRunFunc(cmd *cobra.Command, args []string) error {
	fmt.Println("parse called", filename)

	playerList, err := generic.ParsePlayerListFromFile(filename)
	if err != nil {
		fmt.Printf("Error: %q\n", err)
		return errors.New("error parsing file")
	}

	for _, player := range playerList {
		fmt.Println(player)
	}

	return nil
}

func init() {
	utilPlayerTableCmd.AddCommand(utilPlayerTableParseCmd)

	utilPlayerTableParseCmd.Flags().StringVarP(&filename, "filepath", "f", "", "file for fielder to read in and parse")
	utilPlayerTableParseCmd.MarkFlagRequired("filepath")
}

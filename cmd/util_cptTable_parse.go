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

// utilCptTableParseCmd represents the parse command
var utilCptTableParseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse a captain's table",
	Long:  `Parse a captain's table.`,
	RunE:  utilCptTableParseRunFunc,
}

func utilCptTableParseRunFunc(cmd *cobra.Command, args []string) error {
	fmt.Println("parse called")

	scoringParams, err := generic.ParseCaptainListFromFile(filename)
	if err != nil {
		fmt.Printf("Error: %q\n", err)
		return errors.New("error parsing file")
	}

	for _, sp := range scoringParams {
		fmt.Println(sp)
	}

	return nil
}

func init() {
	utilCptTableCmd.AddCommand(utilCptTableParseCmd)

	utilCptTableParseCmd.Flags().StringVarP(&filename, "filepath", "f", "", "file for fielder to read in and parse")
	utilCptTableParseCmd.MarkFlagRequired("filepath")
}

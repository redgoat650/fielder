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
	"fmt"
	"io"
	"os"

	"github.com/redgoat650/fielder/parsing/generic"
	"github.com/spf13/cobra"
)

var (
	outputFile string
	numPlayers int
)

// utilPlayerTableGenerateCmd represents the generate command
var utilPlayerTableGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates a player table",
	Long:  `Generates a CSV player table to be filled out easily in an editor.`,
	RunE:  utilPlayerTableGenerateRunFunc,
}

func utilPlayerTableGenerateRunFunc(cmd *cobra.Command, args []string) error {
	fmt.Println("generate called")

	var w io.Writer

	switch outputFile {
	case "":
		w = os.Stdout

	default:
		var err error
		w, err = os.Create(outputFile)
		if err != nil {
			return err
		}

	}

	return generic.WriteSamplePlayerTable(w, numPlayers)
}

func init() {
	utilPlayerTableCmd.AddCommand(utilPlayerTableGenerateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	utilPlayerTableGenerateCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write the output to a file")
	utilPlayerTableGenerateCmd.Flags().IntVarP(&numPlayers, "num-players", "n", 3, "Number of players to add to table")
}

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
	playerTableFile string
)

// utilCptTableGenerateCmd represents the generate command
var utilCptTableGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate captain tables",
	Long:  `Generate a captain table from scratch or from a player table file`,
	RunE:  utilCptTableGenerateRunFunc,
}

func utilCptTableGenerateRunFunc(cmd *cobra.Command, args []string) error {
	fmt.Println("generate called")

	var w io.Writer = os.Stdout

	if outputFile != "" {
		var err error
		w, err = os.Create(outputFile)
		if err != nil {
			return err
		}
	}

	if playerTableFile != "" {
		return generic.WriteSampleCaptainFileFromPlayerFile(w, playerTableFile)
	}

	return generic.WriteSampleCaptainTable(w, numPlayers)
}

func init() {
	utilCptTableCmd.AddCommand(utilCptTableGenerateCmd)
	utilCptTableGenerateCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write the output to a file")
	utilCptTableGenerateCmd.Flags().IntVarP(&numPlayers, "num-players", "n", 3, "Number of players to add to table")
	utilCptTableGenerateCmd.Flags().StringVarP(&playerTableFile, "player-table", "p", "", "Player table file to use as a basis for generation")
}

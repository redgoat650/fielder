package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an item to data",
	Long:  `A common action to add something to the fielder data`,
}

func init() {

	var (
		numInnings    int
		startTimeDesc string
		opponentTeam  string
		gameDetails   string
	)

	var addGameCmd = &cobra.Command{
		Use:   "game",
		Short: "Add a game to a season",
		Long:  `Add a game to a given season`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if gSeason == nil {
				return fmt.Errorf("Season couldn't be created - can't add a game")
			}

			gSeason.AddGame(numInnings, startTimeDesc, opponentTeam, gameDetails)

			return nil

		},
	}

	addGameCmd.Flags().IntVarP(&numInnings, "numInnings", "i", 5, "Number of innings this game will have.")
	addGameCmd.Flags().StringVarP(&startTimeDesc, "startTime", "t", "", "Start time of the game")
	addGameCmd.Flags().StringVarP(&opponentTeam, "oppTeam", "o", "", "Name of the opposing team")
	addGameCmd.Flags().StringVarP(&gameDetails, "details", "d", "", "Any additional details for the game")

	rootCmd.AddCommand(addCmd)

	addCmd.AddCommand(addGameCmd)
}

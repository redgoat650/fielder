package cmd

import (
	"fmt"

	fielder "github.com/redgoat650/fielder/scheduling"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a fielder item",
	Long:  `A common action to create a fielder action`,
}

func init() {

	var (
		seasonName  string
		numGames    int
		innsPerGame int
	)

	var newSeasonCmd = &cobra.Command{
		Use:   "season",
		Short: "Create a new season",
		Long:  `Create a new season`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("creating a new season")

			gSeason = fielder.NewSeason(numGames, innsPerGame)
			gSeason.Desc = seasonName

		},
	}

	var newTeamCmd = &cobra.Command{
		Use:   "team",
		Short: "Create a new team",
		Long:  `Create a new team`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("creating a new team")

			gSeason = fielder.NewSeason(numGames, innsPerGame)
			gSeason.Desc = seasonName

		},
	}

	newSeasonCmd.Flags().StringVarP(&seasonName, "name", "n", "my-season", "Pick a name for this season.")
	newSeasonCmd.Flags().StringVarP(&innsPerGame, "inningPerGame", "i", 5, "Number of innings per game.")
	newSeasonCmd.Flags().IntVarP(&numGames, "numGames", "g", 0, "Number of games for this season.")
	newSeasonCmd.Flags().IntVarP(&innsPerGame, "inningPerGame", "i", 5, "Number of innings per game.")

	rootCmd.AddCommand(createCmd)

	createCmd.AddCommand(newSeasonCmd)
}

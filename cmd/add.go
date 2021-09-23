package cmd

// import (
// 	"fmt"
// 	"math"

// 	"github.com/spf13/cobra"
// )

// var addCmd = &cobra.Command{
// 	Use:   "add",
// 	Short: "Add an item to data",
// 	Long:  `A common action to add something to the fielder data`,
// 	PreRunE: func(cmd *cobra.Command, args []string) error {
// 		if gTeam == nil {
// 			return fmt.Errorf("Team couldn't be loaded, so can't add any properties")
// 		}
// 		return nil
// 	},
// }

// func init() {

// 	var (
// 		addGameToSeason int
// 		numInnings      int
// 		startTimeDesc   string
// 		opponentTeam    string
// 		gameDetails     string
// 	)

// 	var addGameCmd = &cobra.Command{
// 		Use:   "game",
// 		Short: "Add a game to a season",
// 		Long:  `Add a game to a given season`,
// 		RunE: func(cmd *cobra.Command, args []string) error {

// 			if len(gTeam.SeasonList) == 0 {
// 				return fmt.Errorf("No seasons have been added for this team: fielder add season --help")
// 			}

// 			season := gTeam.SeasonList[len(gTeam.SeasonList)-1-int(math.Abs(float64(addGameToSeason)))]

// 			season.AddGame(numInnings, startTimeDesc, opponentTeam, gameDetails)

// 			return nil

// 		},
// 	}

// 	var addSeasonCmd = &cobra.Command{
// 		Use:   "season",
// 		Short: "Add a season for this team",
// 		Long: `Add a season to a team's career. Provide "season" parameter
// 		as the number of seasons from the most recent to add this game. Defaults
// 		to zero, which adds the game to the most recent season. Setting it to 1
// 		will add it to the previous season`,
// 		RunE: func(cmd *cobra.Command, args []string) error {

// 			season := gTeam.SeasonList[len(gTeam.SeasonList)-1-int(math.Abs(float64(addGameToSeason)))]

// 			season.AddGame(numInnings, startTimeDesc, opponentTeam, gameDetails)

// 			return nil

// 		},
// 	}

// 	addGameCmd.Flags().IntVarP(&addGameToSeason, "season", "s", 0, "Season to add game to. Zero adds to the most recent season, 1/-1 adds to last season, etc")
// 	addGameCmd.Flags().IntVarP(&numInnings, "numInnings", "i", 5, "Number of innings this game will have.")
// 	addGameCmd.Flags().StringVarP(&startTimeDesc, "startTime", "t", "", "Start time of the game")
// 	addGameCmd.Flags().StringVarP(&opponentTeam, "oppTeam", "o", "", "Name of the opposing team")
// 	addGameCmd.Flags().StringVarP(&gameDetails, "details", "d", "", "Any additional details for the game")

// 	rootCmd.AddCommand(addCmd)

// 	addCmd.AddCommand(addGameCmd)
// 	addCmd.AddCommand(addSeasonCmd)
// }

package cmd

import (
	"fmt"

	"github.com/redgoat650/fielder/parsing/buzzedsheets"
	"github.com/spf13/cobra"
)

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse input data from file",
	Long:  `Parse data from input files to do fielder compute`,
}

func init() {

	var (
		scheduleCSV string
		prefCSV     string
		gameDate    string
	)

	var buzzedCmd = &cobra.Command{
		Use:   "buzzed",
		Short: "Parse buzzed lightyear format CSV",
		Long: `Parse a roster and schedule from a buzzed lightyear
format legacy spreadsheet in CSV format.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			if scheduleCSV == "" {
				return fmt.Errorf("No CSV schedule provided")
			}
			if prefCSV == "" {
				return fmt.Errorf("No CSV preferences provided")
			}
			if gameDate == "" {
				return fmt.Errorf("No game date provided")
			}

			game, err := buzzedsheets.ParseBuzzedSheets(scheduleCSV, prefCSV, gameDate)
			if err != nil {
				return err
			}
			err = game.ScheduleGame2()
			if err != nil {
				return err
			}
			fmt.Println(game)
			return nil

		},
	}

	buzzedCmd.Flags().StringVarP(&scheduleCSV, "schedule", "s", "", "Path to schedule file")
	buzzedCmd.Flags().StringVarP(&prefCSV, "preferences", "p", "", "Path to preferences file")
	buzzedCmd.Flags().StringVarP(&gameDate, "gameDate", "d", "", "Date of the game to generate schedule data for")

	rootCmd.AddCommand(parseCmd)

	parseCmd.AddCommand(buzzedCmd)
}

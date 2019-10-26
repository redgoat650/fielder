package cmd

import (
	"fmt"
	"time"

	"github.com/redgoat650/fielder/parsing/buzzedsheets"
	"github.com/redgoat650/fielder/storage"
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
		cptPrefCSV  string
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
			if gameDate == "" {
				return fmt.Errorf("No game date provided")
			}

			game, err := buzzedsheets.ParseBuzzedSheets(scheduleCSV, prefCSV, cptPrefCSV, gameDate)
			if err != nil {
				return err
			}
			err = game.ScheduleGame2()
			if err != nil {
				return err
			}
			fmt.Println(game)

			savePath := fmt.Sprintf("testdata/gen/game_%v", time.Now().Unix())
			err = storage.SaveGob(savePath, game)
			if err != nil {
				return err
			}

			// newGame := &fielder.Game{}
			// dec, err := storage.LoadGob(savePath)
			// if err != nil {
			// 	return err
			// }
			// err = dec.Decode(newGame)
			// if err != nil {
			// 	return err
			// }

			// if game.String() != newGame.String() {
			// 	fmt.Println(newGame)
			// 	panic("what")
			// }

			return nil

		},
	}

	buzzedCmd.Flags().StringVarP(&scheduleCSV, "schedule", "s", "", "Path to schedule file")
	buzzedCmd.Flags().StringVarP(&prefCSV, "preferences", "p", "", "Path to preferences file")
	buzzedCmd.Flags().StringVarP(&gameDate, "gameDate", "d", "", "Date of the game to generate schedule data for")
	buzzedCmd.Flags().StringVarP(&cptPrefCSV, "cptPref", "c", "", "Custom configuration file for finer control of player position preferences. If the file doesn't exist it will be created")

	rootCmd.AddCommand(parseCmd)

	parseCmd.AddCommand(buzzedCmd)
}

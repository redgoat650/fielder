package cmd

import (
	"fmt"
	"os"

	fielder "github.com/redgoat650/fielder/scheduling"
	"github.com/spf13/cobra"
)

var (
	filename string
	// gSeason  *fielder.Season
	gTeam *fielder.Team
)

var rootCmd = &cobra.Command{
	Use:   "fielder",
	Short: "Schedule positions for a team in a field",
	Long: `A scheduler that distributes players into positions on the field
		based on preference, seniority, and equal playing time.`,
	Args: cobra.MinimumNArgs(1),
	// PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {

	// 	if strings.Contains(cmd.CommandPath(), "create") {
	// 		return nil
	// 	}
	// 	if filename == "" {
	// 		return fmt.Errorf("No filename provided")
	// 	}

	// 	gTeam, err = fielder.LoadTeamFromFile(filename)
	// 	return err

	// },
	// PersistentPostRunE: func(cmd *cobra.Command, args []string) (err error) {

	// 	if strings.Contains(cmd.CommandPath(), "print") {
	// 		return nil
	// 	}

	// 	if filename != "" {
	// 		fmt.Println("saving season to file", filename)

	// 		err = fielder.SaveTeamToFile(gTeam, filename)
	// 	} else {
	// 		fmt.Println(gTeam)
	// 	}

	// 	return

	// },
}

// Execute executes the command line framework
func Execute() {

	//Add args
	rootCmd.PersistentFlags().StringVarP(&filename, "filename", "f", "", "Filename for the operation")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

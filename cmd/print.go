package cmd

// import (
// 	"fmt"

// 	"github.com/spf13/cobra"
// )

// var printCmd = &cobra.Command{
// 	Use:   "print",
// 	Short: "print a fielder item",
// 	Long:  `A common action to print an aspect of the fielder data set`,
// }

// func init() {

// 	var printTeamCmd = &cobra.Command{
// 		Use:   "team",
// 		Short: "print team info",
// 		Long:  `Print the information related to the team`,
// 		RunE: func(cmd *cobra.Command, args []string) (err error) {

// 			fmt.Println(gTeam)

// 			return nil
// 		},
// 	}
// 	// var printSeasonCmd = &cobra.Command{
// 	// 	Use:   "season",
// 	// 	Short: "print season info",
// 	// 	Long:  `Print the information related to the season`,
// 	// 	RunE: func(cmd *cobra.Command, args []string) (err error) {

// 	// 		fmt.Println(gSeason)

// 	// 		return nil
// 	// 	},
// 	// }

// 	rootCmd.AddCommand(printCmd)

// 	printCmd.AddCommand(printTeamCmd)
// }

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fielder",
	Short: "Schedule positions for a team in a field",
	Long: `A scheduler that distributes players into positions on the field
		based on preference, seniority, and equal playing time.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

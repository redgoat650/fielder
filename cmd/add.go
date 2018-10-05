package cmd

import (
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an item to data",
	Long:  `A common action to add something to the fielder data`,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

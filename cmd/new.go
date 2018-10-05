package cmd

import (
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a fielder item",
	Long:  `A common action to create a fielder action`,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

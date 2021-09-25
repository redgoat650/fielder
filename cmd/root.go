/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	fielder "github.com/redgoat650/fielder/scheduling"
)

var (
	cfgFile       string
	gTeam         *fielder.Team
	dataDirParent string
)

const (
	teamsDirName          = "teams"
	selectedTeamConfigKey = "selectedTeam"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fielder",
	Short: "Schedule positions for a team in a field",
	Long: `A scheduler that distributes players into positions on the field
	based on preference, seniority, and equal playing time.`,

	PersistentPreRunE: rootPersistentPreRunFunc,

	PersistentPostRunE: rootPersistentPostRunFunc,
}

func skipLoadTeam(cmd *cobra.Command) bool {
	switch cmd.Name() {
	case "create", "switch", "list", "delete":
		return cmd.Parent().Name() == "team"
	}
	return false
}

func rootPersistentPreRunFunc(cmd *cobra.Command, args []string) error {
	fmt.Println("fielder pre run")

	if skipLoadTeam(cmd) {
		return nil
	}

	return loadTeam()
}

func rootPersistentPostRunFunc(cmd *cobra.Command, args []string) error {
	fmt.Println("fielder post run")
	if gTeam != nil {
		return writeGlobalTeam()
	}

	return nil
}

func writeGlobalTeam() error {
	b, err := json.Marshal(gTeam)
	if err != nil {
		return err
	}

	teamsDir := filepath.Join(dataDirParent, teamsDirName)
	_ = os.MkdirAll(teamsDir, 0755)

	filename := getFullTeamFilePath(gTeam.TeamName)

	return os.WriteFile(filename, b, 0755)
}

func getFullTeamFilePath(teamName string) string {
	teamsDir := getTeamsDir()
	filename := filepath.Join(teamsDir, teamName+".json")

	return filename
}

func getTeamsDir() string {
	return filepath.Join(dataDirParent, teamsDirName)
}

func loadTeam() error {
	fmt.Println("Loading team")

	teamName, err := getTeamNameFromViper()
	if err != nil {
		return err
	}

	return loadTeamByName(teamName)
}

func loadTeamByName(teamName string) error {
	filename := getFullTeamFilePath(teamName)

	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &gTeam)
	if err != nil {
		return err
	}

	fmt.Println("Using existing team", gTeam.TeamName)
	return nil
}

func getTeamNameFromViper() (string, error) {
	if !viper.IsSet(selectedTeamConfigKey) {
		return "", errors.New("no team selected")
	}

	teamName := viper.Get(selectedTeamConfigKey)

	teamNameStr, ok := teamName.(string)
	if !ok {
		return "", errors.New("team name is not string format")
	}

	return teamNameStr, nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fielder.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		dataDirParent = filepath.Dir(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".fielder" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".fielder")

		dataDirParent = home
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

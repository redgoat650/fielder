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
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// seasonDescribeCmd represents the describe command
var seasonDescribeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe season",
	Long:  `Describe the selected season.`,
	RunE:  seasonDescribeRunFunc,
}

func seasonDescribeRunFunc(cmd *cobra.Command, args []string) error {
	fmt.Println("describe called")

	if !viper.IsSet(selectedSeasonConfigKey) {
		fmt.Println("Select a season with 'fielder season select <name>'")
		return errors.New("no season selected")
	}

	seasonName := viper.GetString(selectedSeasonConfigKey)

	if _, ok := gTeam.SeasonList[seasonName]; !ok {
		return errors.New("season not found")
	}

	fmt.Println(gTeam.SeasonList[seasonName])

	return nil
}

func init() {
	seasonCmd.AddCommand(seasonDescribeCmd)
}

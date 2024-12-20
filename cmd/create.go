/*
Copyright © 2024 Jan Nawrocki jan.nawrocki06@gmail.com

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
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a notification or a reminder",
	Long: `Create a notification or a reminder. Supply a descriptive title for the notifications.
Supply a time string interval to the "--interval" flag to create a reminder that
fires repetitively.

Examples:

$ terminder create take the trash out

$ terminder create do my taxes -i 1y`,
	RunE: func(cmd *cobra.Command, args []string) error {
		interval, err := cmd.Flags().GetString("interval")
		if err != nil {
			return err
		}
		switch interval {
		case "":
			notification, err := app.CreateNotification(cmd.Context(), strings.Join(args, " "))
			if err != nil {
				return err
			}
			fmt.Println("Succesfully created notification: ", notification.Body)
		default:
			reminder, err := app.CreateReminder(cmd.Context(), strings.Join(args, " "), interval)
			if err != nil {
				return err
			}
			fmt.Println("Succesfully created reminder: ", reminder.Body)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringP("interval", "i", "", "the interval for fireing a notification repetitively")
}

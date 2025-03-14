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
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// dismissCmd represents the dismiss command
var dismissCmd = &cobra.Command{
	Use:   "dismiss",
	Short: "Dismiss a notification",
	Long: `Dismiss a notification. Expects either no args or the ID of the notification 
you'd like to dimiss. When given no args a selection menu appears.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		forever, err := cmd.Flags().GetBool("forever")
		if err != nil {
			return err
		}

		if len(args) == 0 {
			notifications, err := app.GetAllActiveNotifications(cmd.Context())
			if err != nil {
				return err
			}

			prompt := promptui.Select{
				Label:     "Select to dismiss",
				IsVimMode: true,
				Items:     notifications,
				Templates: &promptui.SelectTemplates{
					Active:   `> {{ .ID }} {{ .Body }}`,
					Inactive: `  {{ .ID }} {{ .Body }}`,
					Selected: `  {{ .ID }} {{ .Body }}`,
				},
			}

			index, _, err := prompt.Run()
			if err != nil {
				return err
			}
			notification, err := app.DismissNotification(cmd.Context(), notifications[index].ID)
			if err != nil {
				return err
			}
			fmt.Println("Successfully dismissed notification:", notification.Body)
			if forever {
				_, err := app.DismissReminder(cmd.Context(), notification)
				if err != nil {
					return err
				}
			}
			return nil
		}

		eg, egContext := errgroup.WithContext(cmd.Context())
		for _, arg := range args {
			eg.Go(func() error {
				a, err := strconv.Atoi(arg)
				if err != nil {
					return err
				}
				notification, err := app.DismissNotification(egContext, int64(a))
				if err != nil {
					return err
				}
				fmt.Println("Successfully dismissed notification:", notification.Body)
				if forever {
					_, err := app.DismissReminder(egContext, notification)
					if err != nil {
						return err
					}
				}
				return nil
			})
		}

		return eg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(dismissCmd)
	dismissCmd.Flags().BoolP("forever", "f", false, "Dismisses the parent reminder in addition to the notification.")
}

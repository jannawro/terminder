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
	"context"
	"os"
	"path/filepath"
	"terminder/repository"
	"terminder/terminder"

	"github.com/spf13/cobra"
)

const terminderDir = ".terminder"

var repo *repository.Repository
var app *terminder.App

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "terminder",
	Short: "Terminder creates reminders in your terminal.",
	Long: `Terminder creates reminders in your terminal.

General workflow:
1. Use 'terminder create' to create notifications and reminders.
2. Do 'echo "\nterminder" >> $HOME/.bashrc' or whatever file you use to configure your shell.
   Now you'll see notifications each time you open a new terminal session.
3. Use 'terminder dismiss' when you no longer need a notification.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		repo, err = repository.NewRepo(filepath.Join(homeDir, terminderDir))
		if err != nil {
			return err
		}
		app = terminder.New(repo)

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := app.FireNotifications(cmd.Context())
		if err != nil {
			return err
		}

		notifications, err := app.GetAllActiveNotifications(cmd.Context())
		if err != nil {
			return err
		}

		return terminder.PrettyPrintNotifications(notifications)
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if repo != nil {
			return repo.Close()
		}
		return nil
	},
}

func Execute() {
	err := rootCmd.ExecuteContext(context.Background())
	if err != nil {
		os.Exit(1)
	}
}

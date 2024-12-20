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
package terminder

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"terminder/repository"
	"text/tabwriter"
	"time"

	tparse "github.com/karrick/tparse/v2"
)

type App struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *App {
	return &App{repo}
}

func (a *App) CreateNotification(ctx context.Context, body string) (repository.Notification, error) {
	notification, err := a.repo.CreateNotification(ctx, repository.CreateNotificationParams{
		ReminderID: sql.NullInt64{Int64: 0, Valid: false},
		Body:       body,
	})
	if err != nil {
		return repository.Notification{}, errors.Join(CreateNotificationErr, err)
	}
	return notification, nil
}

func (a *App) CreateReminder(ctx context.Context, body string, interval string) (repository.Reminder, error) {
	_, err := tparse.AddDuration(time.Now(), interval)
	if err != nil {
		return repository.Reminder{}, errors.Join(InvalidIntervalErr, err)
	}
	reminder, err := a.repo.CreateReminder(ctx, repository.CreateReminderParams{
		Interval: interval,
		Body:     body,
	})
	if err != nil {
		return repository.Reminder{}, errors.Join(CreateReminderErr, err)
	}
	return reminder, nil
}

func (a *App) DismissNotification(ctx context.Context, id int64) (repository.Notification, error) {
	notification, err := a.repo.DismissNotification(ctx, id)
	if err != nil {
		return repository.Notification{}, errors.Join(DismissNotificationErr, err)
	}
	return notification, nil
}

func (a *App) DismissReminder(ctx context.Context, childNotification repository.Notification) (repository.Reminder, error) {
	reminder, err := a.repo.DismissReminder(ctx, childNotification.ReminderID.Int64)
	if err != nil {
		return repository.Reminder{}, errors.Join(DismissReminderErr, err)
	}
	return reminder, nil
}

func (a *App) GetAllActiveNotifications(ctx context.Context) ([]repository.Notification, error) {
	notifications, err := a.repo.GetAllActiveNotifications(ctx)
	if err != nil {
		return nil, errors.Join(GetActiveNotificationsErr, err)
	}
	return notifications, nil
}

func (a *App) FireNotifications(ctx context.Context) ([]repository.Notification, error) {
	activeReminders, err := a.repo.GetAllActiveReminders(ctx)
	if err != nil {
		return nil, errors.Join(GetActiveRemindersErr, err)
	}

	var newNotifications []repository.Notification
	for _, reminder := range activeReminders {
		now := time.Now()
		lastFired := reminder.LastFired
		lastFiredPlusInterval, err := tparse.AddDuration(lastFired.Time, reminder.Interval)
		if err != nil {
			return nil, errors.Join(InvalidIntervalErr, err)
		}

		shouldFire := !lastFired.Valid || now.After(lastFiredPlusInterval)

		if shouldFire {
			notification, err := a.repo.CreateNotification(
				ctx, repository.CreateNotificationParams{
					ReminderID: sql.NullInt64{
						Int64: reminder.ID,
						Valid: true,
					},
					Body: reminder.Body,
				})
			if err != nil {
				return nil, errors.Join(CreateNotificationErr, err)
			}
			newNotifications = append(newNotifications, notification)

			_, err = a.repo.UpdateReminderLastFired(ctx, reminder.ID)
			if err != nil {
				return nil, errors.Join(UpdateReminderErr, err)
			}
		}
	}

	return newNotifications, nil
}

func PrettyPrintNotifications(notifications []repository.Notification) error {
	if len(notifications) == 0 {
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNOTIFICATION")
	for _, notification := range notifications {
		fmt.Fprintf(w, "%d\t%s\n", notification.ID, notification.Body)
	}
	err := w.Flush()
	if err != nil {
		return errors.Join(PrettyPrintErr, err)
	}
	return nil
}

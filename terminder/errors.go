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

import "errors"

var (
	CreateNotificationErr     = errors.New("failed to create a notification")
	CreateReminderErr         = errors.New("failed to create a reminder")
	DismissNotificationErr    = errors.New("failed to dismiss a notification")
	DismissReminderErr        = errors.New("failed to dismiss a reminder")
	GetActiveNotificationsErr = errors.New("failed to get all active notifications")
	GetActiveRemindersErr     = errors.New("failed to get all active reminders")
	InvalidIntervalErr        = errors.New("interval for this reminder is invalid")
	UpdateReminderErr         = errors.New("failed to update a reminder")
	PrettyPrintErr            = errors.New("failed to pretty print notifications")
)

-- name: GetAllActiveReminders :many
SELECT 
    id,
    interval,
    creation_date,
    body,
    dismissal_date,
    last_fired
FROM reminders
WHERE dismissal_date IS NULL
ORDER BY creation_date ASC;

-- name: GetAllActiveNotifications :many
SELECT 
    id,
    reminder_id,
    body,
    creation_date,
    dismissal_date
FROM notifications
WHERE dismissal_date IS NULL
ORDER BY creation_date ASC;

-- name: CreateNotification :one
INSERT INTO notifications (
    reminder_id,
    body,
    creation_date
) VALUES (
    ?,
    ?,
    CURRENT_TIMESTAMP
)
RETURNING id, reminder_id, body, creation_date, dismissal_date;

-- name: CreateReminder :one
INSERT INTO reminders (
    interval,
    body,
    creation_date
) VALUES (
    ?,
    ?,
    CURRENT_TIMESTAMP
)
RETURNING id, interval, creation_date, body, dismissal_date, last_fired;

-- name: DismissNotification :one
UPDATE notifications
SET dismissal_date = CURRENT_TIMESTAMP
WHERE id = ? AND dismissal_date IS NULL
RETURNING id, reminder_id, body, creation_date, dismissal_date;


-- name: DismissReminder :one
UPDATE reminders
SET dismissal_date = CURRENT_TIMESTAMP
WHERE id = ? AND dismissal_date IS NULL
RETURNING id, interval, creation_date, body, dismissal_date, last_fired;

-- name: UpdateReminderLastFired :one
UPDATE reminders 
SET last_fired = CURRENT_TIMESTAMP
WHERE id = ? AND dismissal_date IS NULL
RETURNING id, interval, creation_date, body, dismissal_date, last_fired;

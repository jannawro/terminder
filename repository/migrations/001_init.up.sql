CREATE TABLE IF NOT EXISTS reminders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    interval TEXT NOT NULL,
    creation_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    body TEXT NOT NULL,
    dismissal_date DATETIME,
    last_fired DATETIME
);

CREATE INDEX IF NOT EXISTS idx_reminders_creation_date ON reminders(creation_date);

CREATE TABLE IF NOT EXISTS notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    reminder_id INTEGER,
    body TEXT NOT NULL,
    creation_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    dismissal_date DATETIME,
    FOREIGN KEY (reminder_id)
        REFERENCES reminders(id)
        ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_notifications_creation_date ON notifications(creation_date);
CREATE INDEX IF NOT EXISTS idx_notifications_reminder_id ON notifications(reminder_id);

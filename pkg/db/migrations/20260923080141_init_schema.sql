-- +goose Up
CREATE TABLE scheduler
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(256) NOT NULL,
    comment TEXT,
    repeat VARCHAR(128)
);
CREATE INDEX scheduler_date on scheduler(date)

-- +goose Down
DROP TABLE scheduler;

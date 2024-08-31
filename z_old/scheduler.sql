DROP TABLE IF EXISTS scheduler;
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date VARCHAR(8) DEFAULT "", -- handle YYYYMMDD format
    title TEXT DEFAULT "",
    comment TEXT DEFAULT "",
    repeat VARCHAR(128) DEFAULT ""
);
CREATE INDEX idx_scheduler_date ON scheduler(date);

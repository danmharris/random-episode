CREATE TABLE watched_episodes(
    id INTEGER PRIMARY KEY,
    show_id INTEGER REFERENCES shows(id),
    season INTEGER,
    episode INTEGER,
    title TEXT,
    timestamp INTEGER);

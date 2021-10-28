
CREATE TABLE tasks_queue (
    task_id INTEGER PRIMARY KEY,
    tile_name TEXT NOT NULL,
    req_date TEXT NOT NULL,
    user_id INTEGER NOT NULL
);
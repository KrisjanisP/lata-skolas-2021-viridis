
CREATE TABLE tile_urls (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    tfw_url TEXT,
    rgb_url TEXT,
    cir_url TEXT
);
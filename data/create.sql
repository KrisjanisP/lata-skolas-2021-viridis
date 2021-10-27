CREATE TABLE tile_coordinates (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    ulcx REAL,
    ulcy REAL,
    urcx REAL,
    urcy REAL,
    brcx REAL,
    brcy REAL,
    blcx REAL,
    blcy REAL
);

CREATE TABLE tile_urls (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    tfw_url TEXT,
    rgb_url TEXT,
    cir_url TEXT
);
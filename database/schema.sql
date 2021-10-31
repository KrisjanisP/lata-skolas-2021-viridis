CREATE TABLE tiles (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);
CREATE TABLE users (
    id INTEGER PRIMARY KEY
);
CREATE TABLE finishedtiles (
    tileid INTEGER REFERENCES tiles(id),
    rgb INTEGER NOT NULL,
    cir INTEGER NOT NULL,
    ndv INTEGER NOT NULL,
    ove INTEGER NOT NULL
);
CREATE TABLE tilepossesion (
    tileid INTEGER REFERENCES tiles(id),
    userid INTEGER REFERENCES users(id)
);
CREATE TABLE tileurls (
    tileid INTEGER REFERENCES tiles(id) UNIQUE,
    tfwurl TEXT,
    rgburl TEXT,
    cirurl TEXT
);
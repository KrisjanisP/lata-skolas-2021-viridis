CREATE TABLE tiles (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);
CREATE UNIQUE INDEX tiles_id_index
ON tiles(id);
CREATE UNIQUE INDEX tiles_name_index
ON tiles(name);
CREATE TABLE finishedtiles (
    tileid INTEGER REFERENCES tiles(id),
    rgb INTEGER NOT NULL,
    cir INTEGER NOT NULL,
    ndv INTEGER NOT NULL,
    ove INTEGER NOT NULL
);
CREATE UNIQUE INDEX finishedtiles_tileid_index
ON finishedtiles(tileid);
CREATE TABLE tilepossesion (
    tileid INTEGER REFERENCES tiles(id),
    userid TEXT NOT NULL
);
CREATE INDEX tilepossesion_tileid_index
ON tilepossesion(tileid);
CREATE INDEX tilepossesion_userid_index
ON tilepossesion(userid);
CREATE TABLE tileurls (
    tileid INTEGER REFERENCES tiles(id) UNIQUE,
    tfwurl TEXT,
    rgburl TEXT,
    cirurl TEXT
);
CREATE UNIQUE INDEX tileurls_tileid_index
ON tileurls(tileid);
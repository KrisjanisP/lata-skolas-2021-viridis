INSERT INTO tile_coordinates(
    id,
    name,
    ulcx,
    ulcy,
    urcx,
    urcy,
    brcx,
    brcy,
    blcx,
    blcy
)
SELECT
    id,
    name,
    ulcx,
    ulcy,
    urcx,
    urcy,
    brcx,
    brcy,
    blcx,
    blcy
FROM tilesold;

INSERT INTO tile_urls(
    id,
    name,
    tfw_url,
    rgb_url,
    cir_url
)
SELECT
    id,
    name,
    tfwURL,
    rgbURL,
    cirURL
FROM tilesold;
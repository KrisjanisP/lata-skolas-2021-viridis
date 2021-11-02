SELECT DISTINCT tilepossesion.tileid,
    tiles.name,
    ifnull(finishedtiles.rgb, 0),
    ifnull(finishedtiles.cir, 0),
    ifnull(finishedtiles.ndv, 0),
    ifnull(finishedtiles.ove, 0)
FROM tilepossesion
    INNER JOIN tiles ON tilepossesion.tileid = tiles.id
    LEFT JOIN finishedtiles on tilepossesion.tileid = finishedtiles.tileid

    WHERE tilepossesion.userid = "auth0|6180299a2e145b006fe5c44b"
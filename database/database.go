package database

// database schema can be found under scripts/database/gen-database.go

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const (
	// tileurls table
	getTileURLsRecordsSQL    = "SELECT * FROM tileurls"
	getTileURLsRecordSQL     = "SELECT * FROM tileurls WHERE tileid = ?"
	updateTileURLsRecordSQL  = "UPDATE tileurls set tfwurl=?,rgburl=?,cirurl=? WHERE tileid=?"
	insertTileURLsRecordSQL  = "INSERT INTO tileurls(tileid,tfwurl,rgburl,cirurl) VALUES(?,?,?,?)"
	replaceTileURLsRecordSQL = "REPLACE INTO tileurls(tileid,tfwurl,rgburl,cirurl) VALUES(?,?,?,?)"

	// tile table
	getTileRecordsSQL           = "SELECT * FROM tiles"
	getTileRecordByNameSQL      = "SELECT * FROM tiles WHERE name = ?"
	getTileRecordByIdSQL        = "SELECT * FROM tiles WHERE id = ?"
	insertTileRecordSQL         = "INSERT INTO tiles(name) VALUES(?)"
	insertOrIgnoreTileRecordSQL = "INSERT OR IGNORE INTO tiles(name) VALUES(?)"
	replaceTileRecordSQL        = "REPLACE INTO tiles(name) VALUES(?)"

	// finishedtiles table
	getFinishedTileRecordsSQL   = "SELECT * FROM finishedtiles"
	getFinishedTileRecordSQL    = "SELECT * FROM finishedtiles WHERE tileid = ?"
	insertFinishedTileRecordSQL = "INSERT INTO finishedtiles(tileid,rgb,cir,ndv,ove) VALUES(?,?,?,?,?)"
	updateFinishedTileRecordSQL = "UPDATE finishedtiles set rgb=?,cir=?,ndv=?,ove=? WHERE tileid=?"

	//tilepossesion table
	insertOrIgnoreTilePossesionRecordSQL   = "INSERT OR IGNORE INTO tilepossesion(tileid, userid) VALUES(?,?)"
	selectPossesionRecordSQL               = "SELECT tileid,userid FROM tilepossesion WHERE tileid=? AND userid=?"
	selectPossesionRecordTileIdsSQL        = "SELECT DISTINCT tileid FROM tilepossesion"
	joinPossesionRecordTileIdsNamesURLsSQL = `
	SELECT DISTINCT
		tilepossesion.tileid,
		tiles.name,
		tileurls.tfwurl,
		tileurls.rgburl,
		tileurls.cirurl
	FROM tilepossesion
	INNER JOIN tiles
		ON tilepossesion.tileid = tiles.id
	INNER JOIN tileurls
		ON tilepossesion.tileid = tileurls.tileid`
	joinPossesionRecordTileIdsNamesFinishedSQL = `
	SELECT DISTINCT
		tilepossesion.tileid,
		tiles.name,
		ifnull(finishedtiles.rgb,0),
		ifnull(finishedtiles.cir,0),
		ifnull(finishedtiles.ndv,0),
		ifnull(finishedtiles.ove,0)
	FROM tilepossesion
		INNER JOIN tiles ON tilepossesion.tileid = tiles.id
		LEFT JOIN finishedtiles on tilepossesion.tileid = finishedtiles.tileid
	    WHERE tilepossesion.userid = ?`
)

type DBAPI struct {
	database *sql.DB

	// tileurls table
	getTileURLsRecordsStmt    *sql.Stmt
	getTileURLsRecordStmt     *sql.Stmt
	updateTileURLsRecordStmt  *sql.Stmt
	insertTileURLsRecordStmt  *sql.Stmt
	replaceTileURLsRecordStmt *sql.Stmt

	// tile table
	getTileRecordsStmt           *sql.Stmt
	getTileRecordByNameStmt      *sql.Stmt
	getTileRecordByIdStmt        *sql.Stmt
	insertTileRecordStmt         *sql.Stmt
	replaceTileRecordStmt        *sql.Stmt
	insertOrIgnoreTileRecordStmt *sql.Stmt

	// finishedtiles table
	getFinishedTileRecordStmt    *sql.Stmt
	getFinishedTileRecordsStmt   *sql.Stmt
	insertFinishedTileRecordStmt *sql.Stmt
	updateFinishedTileRecordStmt *sql.Stmt

	// tilepossesion table
	insertOrIgnoreTilePossesionRecordStmt       *sql.Stmt
	selectPossesionRecordTileIdsStmt            *sql.Stmt
	selectPossesionRecordStmt                   *sql.Stmt
	joinPossesionRecordTileIdsNamesURLsStmt     *sql.Stmt
	joinPossesionRecordTileIdsNamesFinishedStmt *sql.Stmt
}

func NewDB() (*DBAPI, error) {
	check := func(e error) {
		if e != nil {
			log.Panic(e)
		}
	}

	dbFile := "./data/db.sqlite?cache=shared&mode=rwc"
	sqlDB, err := sql.Open("sqlite3", dbFile)
	check(err)

	// tileurls table
	getTileURLsRecords, err := sqlDB.Prepare(getTileURLsRecordsSQL)
	check(err)
	getTileURLsRecord, err := sqlDB.Prepare(getTileURLsRecordSQL)
	check(err)
	updateTileURLsRecord, err := sqlDB.Prepare(updateTileURLsRecordSQL)
	check(err)
	insertTileURLsRecord, err := sqlDB.Prepare(insertTileURLsRecordSQL)
	check(err)
	replaceTileURLsRecord, err := sqlDB.Prepare(replaceTileURLsRecordSQL)
	check(err)

	// tile table
	getTileRecords, err := sqlDB.Prepare(getTileRecordsSQL)
	check(err)
	getTileRecordByName, err := sqlDB.Prepare(getTileRecordByNameSQL)
	check(err)
	getTileRecordById, err := sqlDB.Prepare(getTileRecordByIdSQL)
	check(err)
	insertTileRecord, err := sqlDB.Prepare(insertTileRecordSQL)
	check(err)
	insertOrIgnoreTileRecord, err := sqlDB.Prepare(insertOrIgnoreTileRecordSQL)
	check(err)
	replaceTileRecord, err := sqlDB.Prepare(replaceTileRecordSQL)
	check(err)

	// finishedtiles table
	getFinishedTileRecord, err := sqlDB.Prepare(getFinishedTileRecordSQL)
	check(err)
	getFinishedTileRecords, err := sqlDB.Prepare(getFinishedTileRecordsSQL)
	check(err)
	insertFinishedTileRecord, err := sqlDB.Prepare(insertFinishedTileRecordSQL)
	check(err)
	updateFinishedTileRecord, err := sqlDB.Prepare(updateFinishedTileRecordSQL)
	check(err)

	// tilepossesion table
	insertOrIgnoreTilePossesionRecord, err := sqlDB.Prepare(insertOrIgnoreTilePossesionRecordSQL)
	check(err)
	selectPossesionRecordTileIds, err := sqlDB.Prepare(selectPossesionRecordTileIdsSQL)
	check(err)
	selectPossesionRecord, err := sqlDB.Prepare(selectPossesionRecordSQL)
	check(err)
	joinPossesionRecordTileIdsNamesURLs, err := sqlDB.Prepare(joinPossesionRecordTileIdsNamesURLsSQL)
	check(err)
	joinPossesionRecordTileIdsNamesFinished, err := sqlDB.Prepare(joinPossesionRecordTileIdsNamesFinishedSQL)
	check(err)

	dbapi := DBAPI{
		database: sqlDB,
		// tilurls table
		getTileURLsRecordsStmt:    getTileURLsRecords,
		getTileURLsRecordStmt:     getTileURLsRecord,
		updateTileURLsRecordStmt:  updateTileURLsRecord,
		insertTileURLsRecordStmt:  insertTileURLsRecord,
		replaceTileURLsRecordStmt: replaceTileURLsRecord,
		// tile table
		getTileRecordsStmt:           getTileRecords,
		getTileRecordByNameStmt:      getTileRecordByName,
		getTileRecordByIdStmt:        getTileRecordById,
		insertTileRecordStmt:         insertTileRecord,
		insertOrIgnoreTileRecordStmt: insertOrIgnoreTileRecord,
		replaceTileRecordStmt:        replaceTileRecord,
		// finished table
		getFinishedTileRecordStmt:    getFinishedTileRecord,
		getFinishedTileRecordsStmt:   getFinishedTileRecords,
		insertFinishedTileRecordStmt: insertFinishedTileRecord,
		updateFinishedTileRecordStmt: updateFinishedTileRecord,
		// tilepossesion table
		insertOrIgnoreTilePossesionRecordStmt:       insertOrIgnoreTilePossesionRecord,
		selectPossesionRecordTileIdsStmt:            selectPossesionRecordTileIds,
		selectPossesionRecordStmt:                   selectPossesionRecord,
		joinPossesionRecordTileIdsNamesURLsStmt:     joinPossesionRecordTileIdsNamesURLs,
		joinPossesionRecordTileIdsNamesFinishedStmt: joinPossesionRecordTileIdsNamesFinished,
	}

	return &dbapi, nil
}

func (dbapi *DBAPI) Close() {
	dbapi.database.Close()
	// tilurls table
	dbapi.getTileURLsRecordsStmt.Close()
	dbapi.getTileURLsRecordStmt.Close()
	dbapi.updateTileURLsRecordStmt.Close()
	dbapi.insertTileURLsRecordStmt.Close()
	dbapi.replaceTileURLsRecordStmt.Close()
	// tile table
	dbapi.getTileRecordsStmt.Close()
	dbapi.getTileRecordByNameStmt.Close()
	dbapi.getTileRecordByIdStmt.Close()
	dbapi.insertTileRecordStmt.Close()
	dbapi.insertOrIgnoreTileRecordStmt.Close()
	dbapi.replaceTileRecordStmt.Close()
	// finished table
	dbapi.getFinishedTileRecordStmt.Close()
	dbapi.getFinishedTileRecordsStmt.Close()
	dbapi.insertFinishedTileRecordStmt.Close()
	dbapi.updateFinishedTileRecordStmt.Close()
	// tilepossesion table
	dbapi.insertOrIgnoreTilePossesionRecordStmt.Close()
	dbapi.selectPossesionRecordTileIdsStmt.Close()
	dbapi.selectPossesionRecordStmt.Close()
	dbapi.joinPossesionRecordTileIdsNamesURLsStmt.Close()
	dbapi.joinPossesionRecordTileIdsNamesFinishedStmt.Close()
}

func (dbapi *DBAPI) GetTileURLsRecords() ([]TileURLs, error) {
	stmt := dbapi.getTileURLsRecordsStmt
	rows, err := stmt.Query()
	if err != nil {
		return []TileURLs{}, err
	}
	defer rows.Close()

	var result []TileURLs
	for rows.Next() {
		var tileURLs TileURLs
		err = rows.Scan(
			&tileURLs.TileId,
			&tileURLs.TfwURL,
			&tileURLs.RgbURL,
			&tileURLs.CirURL)
		if err != nil {
			return []TileURLs{}, err
		}
		result = append(result, tileURLs)
	}

	err = rows.Err()
	if err != nil {
		return []TileURLs{}, err
	}

	return result, nil
}

func (dbapi *DBAPI) GetTileURLsRecord(tileId int64) (TileURLs, error) {
	stmt := dbapi.getTileURLsRecordStmt
	row := stmt.QueryRow(tileId)

	var result TileURLs
	err := row.Scan(&result.TileId, &result.TfwURL, &result.RgbURL, &result.CirURL)

	if err != nil {
		return TileURLs{}, err
	}

	return result, nil
}

// Returns either number of rows affected by the update or error.
func (dbapi *DBAPI) UpdateTileURLsRecord(tileURLs TileURLs) (int64, error) {
	stmt := dbapi.updateTileURLsRecordStmt

	res, err := stmt.Exec(tileURLs.TfwURL, tileURLs.RgbURL, tileURLs.CirURL, tileURLs.TileId)
	if err != nil {
		return 0, err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affect, nil
}

func (dbapi *DBAPI) InsertTileURLsRecord(tileURLs TileURLs) error {
	stmt := dbapi.insertTileURLsRecordStmt

	_, err := stmt.Exec(tileURLs.TileId, tileURLs.TfwURL, tileURLs.RgbURL, tileURLs.CirURL)
	if err != nil {
		return err
	}
	return nil
}

func (dbapi *DBAPI) ReplaceTileURLsRecords(tileURLsArr []TileURLs) error {
	tx, err := dbapi.database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, tile := range tileURLsArr {
		stmt := tx.Stmt(dbapi.replaceTileURLsRecordStmt)
		_, err := stmt.Exec(tile.TileId, tile.TfwURL, tile.RgbURL, tile.CirURL)
		if err != nil {
			tx.Rollback()
			log.Println(err)
			return err
		}
	}

	return tx.Commit()
}

func (dbapi *DBAPI) InsertTileRecords(tiles []Tile) error {
	tx, err := dbapi.database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, tile := range tiles {
		stmt := tx.Stmt(dbapi.insertTileRecordStmt)
		_, err := stmt.Exec(tile.Name)
		if err != nil {
			tx.Rollback()
			log.Println(err)
			return err
		}
	}

	return tx.Commit()
}
func (dbapi *DBAPI) InsertOrIgnoreTileRecords(tiles []Tile) error {
	tx, err := dbapi.database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, tile := range tiles {
		stmt := tx.Stmt(dbapi.insertOrIgnoreTileRecordStmt)
		_, err := stmt.Exec(tile.Name)
		if err != nil {
			tx.Rollback()
			log.Println(err)
			return err
		}
	}

	return tx.Commit()
}
func (dbapi *DBAPI) ReplaceTileRecords(tiles []Tile) error {
	tx, err := dbapi.database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, tile := range tiles {
		stmt := tx.Stmt(dbapi.replaceTileRecordStmt)
		_, err := stmt.Exec(tile.Name)
		if err != nil {
			tx.Rollback()
			log.Println(err)
			return err
		}
	}

	return tx.Commit()
}

func (dbapi *DBAPI) GetTileId(tileName string) (int64, error) {
	stmt := dbapi.getTileRecordByNameStmt
	var result Tile
	row := stmt.QueryRow(tileName)
	err := row.Scan(&result.Id, &result.Name)

	if err != nil {
		return 0, err
	}
	return result.Id, nil
}

func (dbapi *DBAPI) GetTileName(tileId int64) (string, error) {
	stmt := dbapi.getTileRecordByIdStmt
	var result Tile
	row := stmt.QueryRow(tileId)
	err := row.Scan(&result.Id, &result.Name)

	if err != nil {
		return "", err
	}
	return result.Name, nil
}

func (dbapi *DBAPI) InsertTile(tileName string) (int64, error) {
	stmt := dbapi.insertTileRecordStmt

	res, err := stmt.Exec(tileName)

	if err != nil {
		return 0, nil
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return id, nil
}

func (dbapi *DBAPI) GetFinishedTilesRecord(tileId int64) (FinishedTile, error) {
	stmt := dbapi.getFinishedTileRecordStmt
	row := stmt.QueryRow(tileId)

	var result FinishedTile
	err := row.Scan(&result.TileId, &result.Rgb, &result.Cir, &result.Ndv, &result.Ove)

	if err != nil {
		return FinishedTile{}, err
	}

	return result, nil
}

func (dbapi *DBAPI) InsertFinishedTilesRecord(finishedTile FinishedTile) error {
	stmt := dbapi.insertFinishedTileRecordStmt

	_, err := stmt.Exec(finishedTile.TileId, finishedTile.Rgb, finishedTile.Cir, finishedTile.Ndv, finishedTile.Ove)
	if err != nil {
		return err
	}
	return nil
}

// Returns either number of rows affected by the update or error.
func (dbapi *DBAPI) UpdateFinishedTileRecord(finishedTile FinishedTile) (int64, error) {
	stmt := dbapi.updateFinishedTileRecordStmt

	res, err := stmt.Exec(finishedTile.Rgb, finishedTile.Cir, finishedTile.Ndv, finishedTile.Ove, finishedTile.TileId)
	if err != nil {
		return 0, err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affect, nil
}

func (dbapi *DBAPI) InsertOrIgnoreTilePossesionRecord(tilePossesion TilePossesion) error {
	stmt := dbapi.insertOrIgnoreTilePossesionRecordStmt

	_, err := stmt.Exec(tilePossesion.TileId, tilePossesion.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (dbapi *DBAPI) InsertOrIgnoreTilePossesionRecords(tilePossesions []TilePossesion) error {
	tx, err := dbapi.database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, tilePossesion := range tilePossesions {
		stmt := tx.Stmt(dbapi.insertOrIgnoreTilePossesionRecordStmt)
		_, err := stmt.Exec(tilePossesion.TileId, tilePossesion.UserId)
		if err != nil {
			tx.Rollback()
			log.Println(err)
			return err
		}
	}
	return tx.Commit()
}

func (dbapi *DBAPI) SelectDistinctPossesionRecordTileIds() ([]int64, error) {
	stmt := dbapi.selectPossesionRecordTileIdsStmt
	rows, err := stmt.Query()
	if err != nil {
		return []int64{}, err
	}
	defer rows.Close()

	var tileIds []int64
	for rows.Next() {
		var tileId int64
		err = rows.Scan(&tileId)
		if err != nil {
			return []int64{}, err
		}
		tileIds = append(tileIds, tileId)
	}

	err = rows.Err()
	if err != nil {
		return []int64{}, err
	}

	return tileIds, nil
}

func (dbapi *DBAPI) JoinPossesionRecordTileIdsNamesURLs() ([]Tile, []TileURLs, error) {
	stmt := dbapi.joinPossesionRecordTileIdsNamesURLsStmt
	rows, err := stmt.Query()
	if err != nil {
		return []Tile{}, []TileURLs{}, err
	}
	defer rows.Close()

	var tiles []Tile
	var tileURLsArr []TileURLs
	for rows.Next() {
		var tileId int64
		var tileName string
		var tfwURL string
		var rgbURL string
		var cirURL string
		err = rows.Scan(&tileId, &tileName, &tfwURL, &rgbURL, &cirURL)
		if err != nil {
			return []Tile{}, []TileURLs{}, err
		}
		tiles = append(tiles, Tile{Id: tileId, Name: tileName})
		tileURLsArr = append(tileURLsArr, TileURLs{TileId: tileId,
			TfwURL: tfwURL, RgbURL: rgbURL, CirURL: cirURL})
	}

	err = rows.Err()
	if err != nil {
		return []Tile{}, []TileURLs{}, err
	}

	return tiles, tileURLsArr, nil
}

func (dbapi *DBAPI) JoinPossesionRecordTileIdsNamesFinished(userId string) ([]Tile, []FinishedTile, error) {
	stmt := dbapi.joinPossesionRecordTileIdsNamesFinishedStmt
	rows, err := stmt.Query(userId)
	if err != nil {
		return []Tile{}, []FinishedTile{}, err
	}
	defer rows.Close()

	var tiles []Tile
	var finishedTiles []FinishedTile
	for rows.Next() {
		var tileId int64
		var tileName string
		var rgb int64
		var cir int64
		var ndv int64
		var ove int64
		err = rows.Scan(&tileId, &tileName, &rgb, &cir, &ndv, &ove)
		if err != nil {
			return []Tile{}, []FinishedTile{}, err
		}
		tiles = append(tiles, Tile{Id: tileId, Name: tileName})
		finishedTiles = append(finishedTiles, FinishedTile{TileId: tileId, Rgb: int(rgb), Cir: int(cir), Ndv: int(ndv), Ove: int(ove)})
	}

	err = rows.Err()
	if err != nil {
		return []Tile{}, []FinishedTile{}, err
	}

	return tiles, finishedTiles, nil
}

func (dbapi *DBAPI) SelectPossesionRecordByUserId(tileid int64, userid string) (TilePossesion, error) {
	stmt := dbapi.selectPossesionRecordStmt
	row := stmt.QueryRow(tileid, userid)

	var result TilePossesion
	err := row.Scan(&result.TileId, &result.UserId)

	if err != nil {
		return TilePossesion{}, err
	}

	return result, nil
}

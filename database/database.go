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
	getTileRecordSQL            = "SELECT * FROM tiles WHERE name = ?"
	insertTileRecordSQL         = "INSERT INTO tiles(name) VALUES(?)"
	insertOrIgnoreTileRecordSQL = "INSERT OR IGNORE INTO tiles(name) VALUES(?)"
	replaceTileRecordSQL        = "REPLACE INTO tiles(name) VALUES(?)"

	// finishedtiles table
	getFinishedTileRecordsSQL   = "SELECT * FROM finishedtiles"
	getFinishedTileRecordSQL    = "SELECT * FROM finishedtiles WHERE tileid = ?"
	insertFinishedTileRecordSQL = "INSERT INTO finishedtiles(tileid,rgb,cir,ndv,ove) VALUES(?,?,?,?,?)"
	updateFinishedTileRecordSQL = "UPDATE finishedtiles set rgb=?,cir=?,ndv=?,ove=? WHERE tileid=?"
)

type DBAPI struct {
	database *sql.DB

	// tileurls table
	getTileURLsRecords    *sql.Stmt
	getTileURLsRecord     *sql.Stmt
	updateTileURLsRecord  *sql.Stmt
	insertTileURLsRecord  *sql.Stmt
	replaceTileURLsRecord *sql.Stmt

	// tile table
	getTileRecords           *sql.Stmt
	getTileRecord            *sql.Stmt
	insertTileRecord         *sql.Stmt
	replaceTileRecord        *sql.Stmt
	insertOrIgnoreTileRecord *sql.Stmt

	// finishedtiles table
	getFinishedTileRecord    *sql.Stmt
	getFinishedTileRecords   *sql.Stmt
	insertFinishedTileRecord *sql.Stmt
	updateFinishedTileRecord *sql.Stmt
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
	getTileRecord, err := sqlDB.Prepare(getTileRecordSQL)
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
	updateFinishedTileRecord, err := sqlDB.Prepare(updateTileURLsRecordSQL)
	check(err)

	dbapi := DBAPI{
		database: sqlDB,
		// tilurls table
		getTileURLsRecords:    getTileURLsRecords,
		getTileURLsRecord:     getTileURLsRecord,
		updateTileURLsRecord:  updateTileURLsRecord,
		insertTileURLsRecord:  insertTileURLsRecord,
		replaceTileURLsRecord: replaceTileURLsRecord,
		// tile table
		getTileRecords:           getTileRecords,
		getTileRecord:            getTileRecord,
		insertTileRecord:         insertTileRecord,
		insertOrIgnoreTileRecord: insertOrIgnoreTileRecord,
		replaceTileRecord:        replaceTileRecord,
		// finished table
		getFinishedTileRecord:    getFinishedTileRecord,
		getFinishedTileRecords:   getFinishedTileRecords,
		insertFinishedTileRecord: insertFinishedTileRecord,
		updateFinishedTileRecord: updateFinishedTileRecord,
	}

	return &dbapi, nil
}

func (dbapi *DBAPI) Close() {
	dbapi.database.Close()
	// tilurls table
	dbapi.getTileURLsRecords.Close()
	dbapi.getTileURLsRecord.Close()
	dbapi.updateTileURLsRecord.Close()
	dbapi.insertTileURLsRecord.Close()
	dbapi.replaceTileURLsRecord.Close()
	// tile table
	dbapi.getTileRecords.Close()
	dbapi.getTileRecord.Close()
	dbapi.insertTileRecord.Close()
	dbapi.insertOrIgnoreTileRecord.Close()
	dbapi.replaceTileRecord.Close()
	// finished table
	dbapi.getFinishedTileRecord.Close()
	dbapi.getFinishedTileRecords.Close()
	dbapi.insertFinishedTileRecord.Close()
	dbapi.updateFinishedTileRecord.Close()
}

func (dbapi *DBAPI) GetTileURLsRecords() ([]TileURLs, error) {
	stmt := dbapi.getTileURLsRecords
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
	stmt := dbapi.getTileURLsRecord
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
	stmt := dbapi.updateTileURLsRecord

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
	stmt := dbapi.insertTileURLsRecord

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
		stmt := tx.Stmt(dbapi.replaceTileURLsRecord)
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
		stmt := tx.Stmt(dbapi.insertTileRecord)
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
		stmt := tx.Stmt(dbapi.insertOrIgnoreTileRecord)
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
		stmt := tx.Stmt(dbapi.replaceTileRecord)
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
	stmt := dbapi.getTileRecord
	var result Tile
	row := stmt.QueryRow(tileName)
	err := row.Scan(&result.Id, &result.Name)

	if err != nil {
		return 0, err
	}
	return result.Id, nil
}

func (dbapi *DBAPI) InsertTile(tileName string) (int64, error) {
	stmt := dbapi.insertTileRecord

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
	stmt := dbapi.getFinishedTileRecord
	row := stmt.QueryRow(tileId)

	var result FinishedTile
	err := row.Scan(&result.TileId, &result.Rgb, &result.Cir, &result.Ndv, &result.Ove)

	if err != nil {
		return FinishedTile{}, err
	}

	return result, nil
}

func (dbapi *DBAPI) InsertFinishedTilesRecord(finishedTile FinishedTile) error {
	stmt := dbapi.insertFinishedTileRecord

	_, err := stmt.Exec(finishedTile.TileId, finishedTile.Rgb, finishedTile.Cir, finishedTile.Ndv, finishedTile.Ove)
	if err != nil {
		return err
	}
	return nil
}

// Returns either number of rows affected by the update or error.
func (dbapi *DBAPI) UpdateFinishedTileRecord(finishedTile FinishedTile) (int64, error) {
	stmt := dbapi.updateFinishedTileRecord

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

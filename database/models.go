package database

// type for tiles table
type Tile struct {
	Id   int64
	Name string
}

// tpye for finishedtiles table
type FinishedTile struct {
	TileId int64
	Rgb    int
	Cir    int
	Ndv    int
	Ove    int
}

// type for tilepossesion table
type TilePossesion struct {
	TileId int64
	UserId string
}

// type for tileurls table
type TileURLs struct {
	TileId int64
	TfwURL string
	RgbURL string
	CirURL string
}

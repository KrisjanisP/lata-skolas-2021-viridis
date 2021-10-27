package models

type TileCoordinates struct {
	ulcx float32
	ulcy float32
	urcx float32
	urcy float32
	brcx float32
	brcy float32
}

type TileURLs struct {
	TfwURL string
	RgbURL string
	CirURL string
}

type Tile struct {
	Id   int
	Name string
	TileCoordinates
	TileURLs
}

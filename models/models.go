package models

type TileURLs struct {
	TfwURL string
	RgbURL string
	CirURL string
}

type Tile struct {
	Id   int
	Name string
	TileURLs
}

type QueueTask struct {
	Id       int
	TileName string
	ReqDate  string
	UserId   int
}

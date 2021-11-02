package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/KrisjanisP/viridis/database"
	"github.com/schollz/progressbar/v3"

	_ "github.com/mattn/go-sqlite3"
)

const (
	urlsUrl = "https://s3.storage.pub.lvdc.gov.lv/lgia-opendata/ortofoto_rgb_v6/LGIA_OpenData_Ortofoto_rgb_v6_saites.txt"
)

func main() {
	dbapi, err := database.NewDB()
	check(err)
	defer dbapi.Close()

	resp, err := http.Get(urlsUrl)
	check(err)
	defer resp.Body.Close()

	urls, err := deRower(resp.Body)
	check(err)

	tileNameRE := regexp.MustCompile(`[\d]{4}-[\d]{2}_[\d]`) // get tile name.

	cnt := len(urls)
	bar := progressbar.Default(int64(cnt))
	var tileNames []string
	for _, url := range urls {
		tileName := tileNameRE.FindString(url)
		tileNames = append(tileNames, tileName)
		bar.Add(1)
	}
	tileNames = removeDuplicateValues(tileNames)
	var tiles []database.Tile
	for _, tileName := range tileNames {
		tiles = append(tiles, database.Tile{Name: tileName})
	}
	err = dbapi.InsertOrIgnoreTileRecords(tiles)
	check(err)
}

func removeDuplicateValues(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func check(e error) {
	if e != nil {
		log.Panic(e)
	}
}

func deRower(data io.ReadCloser) (rows []string, err error) {
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		rows = append(rows, scanner.Text())
	}
	return rows, scanner.Err()
}

package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"regexp"
	"sort"

	"github.com/KrisjanisP/viridis/database"
	"github.com/KrisjanisP/viridis/models"
	"github.com/schollz/progressbar/v3"

	_ "github.com/mattn/go-sqlite3"
)

const (
	rgbUrlsUrl = "https://s3.storage.pub.lvdc.gov.lv/lgia-opendata/ortofoto_rgb_v6/LGIA_OpenData_Ortofoto_rgb_v6_saites.txt"
	cirUrlsUrl = "https://s3.storage.pub.lvdc.gov.lv/lgia-opendata/ortofoto_cir_v6/LGIA_OpenData_Ortofoto_cir_v6_saites.txt"
)

func main() {
	dbapi, err := database.NewDB()
	check(err)
	defer dbapi.Close()
	m := make(map[int64]models.TileURLs)
	procFirstSet(dbapi, m)
	procSecondSet(dbapi, m)
	var tileURLsArr []models.TileURLs
	for _, tileUrls := range m {
		tileURLsArr = append(tileURLsArr, tileUrls)
	}
	sort.Slice(tileURLsArr, func(i, j int) bool {
		return tileURLsArr[i].TileId < tileURLsArr[j].TileId
	})
	dbapi.ReplaceTileURLsRecords(tileURLsArr)
}

func procFirstSet(dbapi *database.DBAPI, m map[int64]models.TileURLs) {
	resp, err := http.Get(rgbUrlsUrl)
	check(err)
	defer resp.Body.Close()

	rgbURLs, err := deRower(resp.Body)
	check(err)

	tileNameRE := regexp.MustCompile(`[\d]{4}-[\d]{2}_[\d]`) // get tile name.
	tfwFileRE := regexp.MustCompile(`.tfw`)
	tifFileRE := regexp.MustCompile(`.tif`)

	cnt := len(rgbURLs)
	bar := progressbar.Default(int64(cnt))
	for _, rgbURL := range rgbURLs {
		bar.Add(1)
		if !urlIsCorrect(rgbURL) {
			continue
		}
		// Extract tile name
		tileName := tileNameRE.FindString(rgbURL)
		// Query for tileid
		tileId, err := dbapi.GetTileId(tileName)
		check(err)

		tileURLs := m[tileId]

		tileURLs.TileId = tileId
		if tifFileRE.MatchString(rgbURL) {
			tileURLs.RgbURL = rgbURL
		} else if tfwFileRE.MatchString(rgbURL) {
			tileURLs.TfwURL = rgbURL
		}
		m[tileId] = tileURLs
	}
}

func procSecondSet(dbapi *database.DBAPI, m map[int64]models.TileURLs) {
	resp, err := http.Get(cirUrlsUrl)
	check(err)
	defer resp.Body.Close()

	cirURLs, err := deRower(resp.Body)
	check(err)

	tileNameRE := regexp.MustCompile(`[\d]{4}-[\d]{2}_[\d]`) // get tile name.
	tfwFileRE := regexp.MustCompile(`.tfw`)
	tifFileRE := regexp.MustCompile(`.tif`)

	cnt := len(cirURLs)
	bar := progressbar.Default(int64(cnt))
	for _, cirURL := range cirURLs {
		bar.Add(1)
		if !urlIsCorrect(cirURL) {
			continue
		}
		// Extract tile name
		tileName := tileNameRE.FindString(cirURL)
		// Query for tileid
		tileId, err := dbapi.GetTileId(tileName)
		check(err)

		tileURLs := m[tileId]

		tileURLs.TileId = tileId
		if tifFileRE.MatchString(cirURL) {
			tileURLs.CirURL = cirURL
		} else if tfwFileRE.MatchString(cirURL) {
			tileURLs.TfwURL = cirURL
		}
		m[tileId] = tileURLs
	}
}

func urlIsCorrect(url string) bool {
	xmlFileRe := regexp.MustCompile(`.xml`)
	errFileRe := regexp.MustCompile(`5412/5411`)
	if xmlFileRe.MatchString(url) {
		return false
	}
	if errFileRe.MatchString(url) {
		return false
	}
	return true
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

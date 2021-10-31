package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/zenthangplus/goccm"

	"github.com/KrisjanisP/viridis/database"
	"github.com/KrisjanisP/viridis/models"
	"github.com/KrisjanisP/viridis/utils"
	geojson "github.com/paulmach/go.geojson"
	"github.com/schollz/progressbar/v3"
)

const (
	geoJSONPath = "./data/tiles.geojson"
)

func check(e error) {
	if e != nil {
		log.Panic(e)
	}
}

func main() {
	dbapi, err := database.NewDB()
	check(err)
	defer dbapi.Close()

	if utils.FileExists(geoJSONPath) {
		fmt.Println("GeoJSON file already exists.")
		fmt.Println("The script is exiting.")
		return
	}

	tiles, err := dbapi.GetTileURLsRecords()
	check(err)

	fc := geojson.NewFeatureCollection()

	cnt := len(tiles)
	bar := progressbar.Default(int64(cnt))
	c := goccm.New(10) // Limit 10 goroutines to run concurrently.
	for _, tile := range tiles {
		// This function has to call before any goroutine
		c.Wait()
		go func(tile models.TileURLs) {
			name := tile.TileId
			//fmt.Println(name)
			tfwURL := tile.TfwURL
			coords := fetchCoords(tfwURL)
			ulc := []float64{coords[0], coords[1]}
			urc := []float64{coords[2], coords[3]}
			brc := []float64{coords[4], coords[5]}
			blc := []float64{coords[6], coords[7]}
			polygon := [][]float64{ulc, urc, brc, blc, ulc}
			polygons := [][][]float64{polygon}
			feature := geojson.NewPolygonFeature(polygons)
			feature.SetProperty("id", name)
			fc.AddFeature(feature)
			bar.Add(1)
			c.Done()
		}(tile)
	}
	c.WaitAllDone()

	rawJSON, err := fc.MarshalJSON()
	check(err)

	err = ioutil.WriteFile(geoJSONPath, rawJSON, 0755)
	check(err)
}

func fetchCoords(tfwURL string) []float64 {

	tfw, err := http.Get(tfwURL)
	check(err)

	defer tfw.Body.Close()

	raw_data, err := ioutil.ReadAll(tfw.Body)
	check(err)

	data, err := deRower(string(raw_data))
	check(err)

	width, err := strconv.ParseFloat(data[0], 64)
	check(err)

	height, err := strconv.ParseFloat(data[3], 64)
	check(err)

	width, height = width*10000, height*10000

	cmd := exec.Command("python3", "./utils/transform.py", data[5], data[4], fmt.Sprintf("%g", width), fmt.Sprintf("%g", height))

	out, err := cmd.Output()
	check(err)

	strCoords, err := deRower(string(out))
	check(err)

	var floatCoords []float64
	for _, strCoord := range strCoords {
		value, err := strconv.ParseFloat(strCoord, 64)
		check(err)
		floatCoords = append(floatCoords, value)
	}
	return floatCoords
}

func deRower(data string) (rows []string, err error) {
	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		rows = append(rows, scanner.Text())
	}
	return rows, scanner.Err()
}

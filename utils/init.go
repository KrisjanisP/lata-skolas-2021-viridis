package utils

import (
	"bufio"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/zenthangplus/goccm"

	"github.com/KrisjanisP/viridis/models"
	geojson "github.com/paulmach/go.geojson"
	"github.com/schollz/progressbar/v3"
)

func getGeoJSONPath() string {
	return "./data/tiles.geojson"
}

func InitGeoJSON(db *sql.DB) {
	if fileExists(getGeoJSONPath()) {
		return
	}

	tiles, err := GetTileUrlsRecords(db)

	if err != nil {
		log.Fatal(err)
	}
	fc := geojson.NewFeatureCollection()

	cnt := len(tiles)
	bar := progressbar.Default(int64(cnt))
	// Limit 3 goroutines to run concurrently.
	c := goccm.New(500)
	for _, tile := range tiles {
		// This function has to call before any goroutine
		c.Wait()
		go func(tile models.Tile) {
			name := tile.Name
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

	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(getGeoJSONPath(), rawJSON, 0755)
	if err != nil {
		log.Fatal(err)
	}
}

func fetchCoords(tfwURL string) []float64 {

	tfw, err := http.Get(tfwURL)

	if err != nil {
		fmt.Println("url sucks ass")
		log.Fatal(err)
		os.Exit(1)
	}

	defer tfw.Body.Close()

	raw_data, err := ioutil.ReadAll(tfw.Body)

	if err != nil {
		fmt.Println("data sucs ass")
		os.Exit(1)
	}

	data, err := deRower(string(raw_data))

	if err != nil {
		fmt.Println("scanner sucs ass")
		os.Exit(1)
	}

	width, err := strconv.ParseFloat(data[0], 64)

	if err != nil {
		fmt.Println("parsefloat1 sucs ass")
		os.Exit(1)
	}

	height, err := strconv.ParseFloat(data[3], 64)

	if err != nil {
		fmt.Println("parsefloat2 sucs ass")
		os.Exit(1)
	}

	width, height = width*10000, height*10000

	cmd := exec.Command("python3", "./utils/transform.py", data[5], data[4], fmt.Sprintf("%g", width), fmt.Sprintf("%g", height))

	out, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	strCoords, err := deRower(string(out))
	if err != nil {
		fmt.Println("go can go and ...")
	}

	var floatCoords []float64
	for _, strCoord := range strCoords {
		value, err := strconv.ParseFloat(strCoord, 64)
		if err != nil {
			fmt.Println("float succ ass")
		}
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

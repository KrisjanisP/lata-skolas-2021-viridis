import "./style.css";
import { Map, View } from "ol";
import TileLayer from "ol/layer/Tile";
import OSM from "ol/source/OSM";
import VectorLayer from "ol/layer/Vector";
import VectorSource from "ol/source/Vector";
import GeoJSON from "ol/format/GeoJSON";
import Select from "ol/interaction/Select";
import Feature from "ol/Feature";
import Point from "ol/geom/Point";
import { Icon, Style } from "ol/style";
import { click } from "ol/events/condition";
import { fromLonLat } from "ol/proj";

const geojsonObject = {
  "type": "FeatureCollection",
  "features": [
    {
      'type': 'Feature',
      "id":"hello1",
      'geometry': {
        'type': 'Point',
        'coordinates': [2924375.3909816, 7507759.76255841],
      },
    }
  ],
};

const vectorSource = new VectorSource({
  features: new GeoJSON().readFeatures(geojsonObject),
});

const vectorLayer = new VectorLayer({
  source: vectorSource,
});

const raster = new TileLayer({
  source: new OSM(),
});


const vector = new VectorLayer({
  source: new VectorSource({
    url: "http://127.0.0.1:5000/",
    format: new GeoJSON({dataProjection: 'EPSG:3857'}),
  }),
});

const map = new Map({
  layers: [raster, vector],
  target: "map",
  view: new View({
    // center on Latvia. these coordinates can be found by map.getView().getCenter()
    center: [2734027.654715377, 7718479.091241797],
    //zoom: 8,
    zoom: 3,
  }),
});

// https://openlayers.org/en/latest/examples/select-features.html

// select interaction working on "click"
const selectClick = new Select({
  condition: click,
});

map.addInteraction(selectClick);

selectClick.on("select", function (e) {
  var feature = e.target.getFeatures().getArray()[0];

  console.log(feature.id_);
});

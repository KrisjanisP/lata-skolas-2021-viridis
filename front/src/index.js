import "./style.css";
import { Map, View } from "ol";
import TileLayer from "ol/layer/Tile";
import OSM from "ol/source/OSM";
import VectorLayer from "ol/layer/Vector";
import VectorSource from "ol/source/Vector";
import VectorImageLayer from 'ol/layer/VectorImage';
import GeoJSON from "ol/format/GeoJSON";
import Select from "ol/interaction/Select";
import { click } from "ol/events/condition";

const raster = new TileLayer({
  source: new OSM(),
});

const vector = new VectorImageLayer({
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
    zoom: 9,
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

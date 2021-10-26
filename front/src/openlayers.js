
import "ol/ol.css";
import { Map, View } from "ol";
import TileLayer from "ol/layer/Tile";
import OSM from "ol/source/OSM";
import VectorSource from "ol/source/Vector";
import VectorImageLayer from 'ol/layer/VectorImage';
import GeoJSON from "ol/format/GeoJSON";
import Style from 'ol/style/Style';
import Fill from 'ol/style/Fill';
import Stroke from 'ol/style/Stroke';
import {
  defaults as defaultInteractions,
} from 'ol/interaction';

export function InitMap(){

    const raster = new TileLayer({
        source: new OSM(),
      });
      
      const vector = new VectorImageLayer({
        source: new VectorSource({
          url: "http://127.0.0.1:5000/tks93tiles",
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
        interactions: defaultInteractions({doubleClickZoom:false})
      });
      
      // https://openlayers.org/en/latest/examples/select-multiple-features.html
      
      const highlightStyle = new Style({
        fill: new Fill({
          color: 'rgb(120,160,71,0.8)',
        }),
        stroke: new Stroke({
          color: '#3399CC',
          width: 3,
        }),
      });
      
      let selected = [];
      
      map.on('click', function (e) {
        map.forEachFeatureAtPixel(e.pixel, function (f) {
          const selIndex = selected.indexOf(f);
          if (selIndex < 0) {
            selected.push(f);
            f.setStyle(highlightStyle);
          } else {
            selected.splice(selIndex, 1);
            f.setStyle(undefined);
          }
          if(selected.length) document.getElementById('fields-submit-btn').classList.remove('disabled')
          else document.getElementById('fields-submit-btn').classList.add('disabled')
          let selectedFields = selected.map((data)=>data.id_)
          console.log(selectedFields)
          document.getElementById('selected-fields').innerHTML=selectedFields.join(", ")
        });
      });
}
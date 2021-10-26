import 'bootstrap'
import 'bootstrap/dist/css/bootstrap.min.css'
import "./style.css";

import { initMap, handleMapClicks, getSelectedFields } from './openlayers.js';

initMap();
handleMapClicks();

window.submitSelectedFields = function(){
  console.log("hello");
  let selected = getSelectedFields()
  console.log(JSON.stringify(selected));
  fetch('http://127.0.0.1:5000/tks93tiles', {
    method: 'POST', // or 'PUT'
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(selected),
  })
  .then(response => response.json())
  .then(data => {
    console.log('Success:', data);
  })
  .catch((error) => {
    console.error('Error:', error);
  });
}
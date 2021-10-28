import 'bootstrap'
import 'bootstrap/dist/css/bootstrap.min.css'
import "./style.css";

import { initMap, handleMapClicks, getSelectedFields } from './openlayers.js';

initMap();
handleMapClicks();

window.submitSelectedFields = function(){

  let selected = getSelectedFields()
  fetch('/tiles', {
    method: 'POST', // or 'PUT'
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(selected),
  })
  .then(response => response.json())
  .then(data => {
    console.log('Success:', data);
    window.location.href = "profile.html";
  })
  .catch((error) => {
    console.error('Error:', error);
  });
}
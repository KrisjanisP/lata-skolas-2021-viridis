import 'bootstrap'
import 'bootstrap/dist/css/bootstrap.min.css'
import "./style.css";
import 'js-cookie';

import { initMap, handleMapClicks, getSelectedFields } from './openlayers.js';
import Cookies from 'js-cookie';

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

window.logOut = function(){
  Cookies.remove('auth-session');
}
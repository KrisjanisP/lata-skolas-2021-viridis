import "bootstrap";
import "js-cookie";

import { initMap, handleMapClicks, getSelectedFields } from "./openlayers.js";
import Cookies from "js-cookie";

initMap();
handleMapClicks();

window.submitSelectedFields = function () {
  let selected = getSelectedFields();
  fetch("/tiles", {
    method: "POST", // or 'PUT'
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(selected),
  })
    .then((response) => response.json())
    .then((data) => {
      console.log("Success:", data);
      window.location.href = "profile.html";
    })
    .catch((error) => {
      console.error("Error:", error);
    });
};

window.logOut = function () {
  Cookies.remove("auth-session");
};

window.setInterval(function () {
  var wait = document.querySelectorAll(".wait");
  for (var i = 0, len = wait.length; i < len; i++) {
    if(wait[i].innerHTML.length>3){
      wait[i].innerHTML = wait[i].innerHTML.substring(
        1,
        wait[i].innerHTML.length
      );
    }
    else if (Math.random() < 0.7) wait[i].innerHTML += ".";
    else
      wait[i].innerHTML = wait[i].innerHTML.substring(
        1,
        wait[i].innerHTML.length
      );
  }
}, 500);
package tile

import (
	"io"
	"net/http"
)

func (i *Implementation) MapHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = io.WriteString(w, `
		<!DOCTYPE html>
		<html>
		<head>
		  <title>Tiler demo map</title>
		
		  <meta charset="utf-8" />
		  <meta name="viewport" content="width=device-width, initial-scale=1.0">
		
		  <link rel="stylesheet" href="https://unpkg.com/leaflet@1.3.1/dist/leaflet.css" integrity="sha512-Rksm5RenBEKSKFjgI3a41vrjkw4EVPlJ3+OiI65vTjIdo9brlAacEuKOiQ5OFh7cOI1bkDwLqdLw3Zg0cRJAAQ==" crossorigin=""/>
		  <script src="https://unpkg.com/leaflet@1.3.1/dist/leaflet.js" integrity="sha512-/Nsx9X4HebavoBvEBuyp3I7od5tA0UzAxs+j83KgC8PU0kgB4XiK4Lfe4y4cgBtaRJQEIFCW+oC506aPT2L1zw==" crossorigin=""></script>
		
		  <style>
			html, body, #map {
			  width: 100%;
			  height: 100%;
			  margin: 0;
			  padding: 0;
			}
		  </style>
		</head>
		
		<body>
		<div id="map"></div>
		
		<script>
		  var map = L.map('map').setView([48.700001, 44.516666], 12);
		
		  L.tileLayer('/tile/default/{z}/{x}/{y}.png', {
			maxZoom: 18,
			attribution: '',
			id: 'base'
		  }).addTo(map);
		</script>
		</body>
		</html>`,
	)
}

# OpenStreetMap Display Debugging Guide

## What Was Fixed

### Root Cause
The map container div was using `v-if="showMap"` which completely removes the element from the DOM when `showMap` is false. This created a chicken-and-egg problem:
1. `initMap()` tries to find the map container with `getElementById()`
2. But the container doesn't exist in the DOM yet because `showMap` is false
3. So `initMap()` returns early and never sets `showMap` to true

### Solution
Changed `v-if="showMap"` to `v-show="showMap"` in EventLocation.vue. This keeps the element in the DOM but hides it with CSS until the map is ready.

## How to Debug Map Display Issues

### 1. Open Browser Console
Open your browser's developer tools (F12) and go to the Console tab.

### 2. Look for EventLocation Logs
When viewing an event with a location, you should see these log messages:

```
[EventLocation] Component mounted { hasLocation: true, location: "New York, NY", canEdit: true, mapId: "map-abc123" }
[EventLocation] initMap called { hasLocation: true, location: "New York, NY", showMap: false, mapId: "map-abc123" }
[EventLocation] Geocoding address: New York, NY
[EventLocation] Geocoding response: [{ lat: "40.7127", lon: "-74.0059", ... }]
[EventLocation] Creating map at coordinates: { lat: "40.7127", lon: "-74.0059" }
[EventLocation] Map initialized successfully
```

### 3. Common Issues and Solutions

#### Issue: "Map container not found"
**Log:** `[EventLocation] Map container not found: map-abc123`
**Cause:** The map container div doesn't exist in the DOM
**Solution:** Ensure you're using `v-show` instead of `v-if` for the map container

#### Issue: "No geocoding results found"
**Log:** `[EventLocation] No geocoding results found for: [location]`
**Cause:** Nominatim couldn't geocode the address
**Solution:** 
- Try a more specific address
- Check if the address format is valid
- Nominatim works best with real addresses (not landmarks or POIs)

#### Issue: "Geocoding failed: 429"
**Log:** `[EventLocation] Error geocoding location: Error: Geocoding failed: 429`
**Cause:** Rate limited by Nominatim (max 1 request per second)
**Solution:** Wait a few seconds and refresh the page

#### Issue: "Error geocoding location: TypeError: Failed to fetch"
**Log:** `[EventLocation] Error geocoding location: TypeError: Failed to fetch`
**Cause:** Network error or CORS issue
**Solution:** 
- Check internet connection
- Check browser console for CORS errors
- Ensure Nominatim API is accessible

### 4. Verify Map Container Exists
In the browser console, run:
```javascript
document.getElementById('map-abc123') // Replace with actual mapId from logs
```
This should return the div element if it exists in the DOM.

### 5. Test Geocoding Manually
In the browser console, run:
```javascript
fetch('https://nominatim.openstreetmap.org/search?format=json&q=New%20York,%20NY', {
  headers: { 'User-Agent': 'Timeful.app (https://timeful.app)' }
})
  .then(r => r.json())
  .then(console.log)
```
This tests if Nominatim geocoding works from your browser.

### 6. Check OpenStreetMap Embed
If geocoding succeeds but the map doesn't display, check the iframe src URL:
```javascript
// Find the iframe in the map container
const mapContainer = document.getElementById('map-abc123') // Replace with actual mapId
const iframe = mapContainer?.querySelector('iframe')
console.log(iframe?.src)
```
The URL should look like:
```
https://www.openstreetmap.org/export/embed.html?bbox=-74.0159,-73.9959&layer=mapnik&marker=40.7127,-74.0059
```

## Testing Different Address Types

Good addresses for testing:
- ✅ "New York, NY" - City names
- ✅ "1600 Pennsylvania Avenue, Washington, DC" - Full street addresses
- ✅ "Eiffel Tower, Paris" - Famous landmarks
- ✅ "San Francisco, California, USA" - Full location names

## Nominatim Rate Limits
- Maximum 1 request per second
- Be respectful of the free service
- Consider caching geocoding results if implementing frequent lookups

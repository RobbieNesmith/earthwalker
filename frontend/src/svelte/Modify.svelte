<script>
    // TODO: this file is getting out of hand

    import { onMount } from 'svelte';
    import { ewapi, globalMap, globalChallenge, globalResult } from '../js/stores.js';
    import { calcTotalScore, showPolygonOnMap } from '../js/earthwalker';
    import L from 'leaflet';
    import '../modify_frontend/modify.css'

    // data fetched from server
    let tileServerURL;
    // score
    let totalScore = 0;
    // map sizing
    const mapSizes = [
        {
            dimensions: [180, 135],  // 1x (can't be used)
            zoom: 0, // Zoom can only be set to integers without using fractional zoom: https://leafletjs.com/examples/zoom-levels/
        },
        {
            dimensions: [300, 225], // 1.666x
            zoom: 0,
        },
        {
            dimensions: [500, 375], // 2.777x
            zoom: 1,
        },
        {
            dimensions: [800, 600], // 4.444x
            zoom: 1,
        },
        {
            dimensions: [1200, 900], // 6.666x
            zoom: 2,
        },
        {
            dimensions: [1600, 1200], // 8.888x
            zoom: 2,
        },
    ];
    // sets title to "earthwalker"
    let titleInterval;
    // decrements timeRemaining once per second
    let timerInterval;

    // DOM elements
    let floatingContainer;
    let guessButton;

    let leafletMap = null;
    let leafletMapPolyGroup;

    // resize StreetView elements
    let innerWidth, innerHeight;
    $: document.body.style.height = innerHeight + "px";
    $: document.body.style.width = innerWidth + "px";

    // settings
    let locStorage = window.localStorage;

    let shrinkMap = locStorage.shrinkMap !== undefined ? JSON.parse(locStorage.shrinkMap) : true;
    $: locStorage.shrinkMap = shrinkMap;

    let storedMapSize = locStorage.storedMapSize !== undefined ? JSON.parse(locStorage.storedMapSize) : 1;
    $: locStorage.storedMapSize = storedMapSize;
    $: curMapSize = shrinkMap && !mapFocused ? 1 : storedMapSize;
    $: if (curMapSize && leafletMap) {
        floatingContainer.style.width = mapSizes[curMapSize].dimensions[0] + "px";
        floatingContainer.style.height = mapSizes[curMapSize].dimensions[1] + "px";
        leafletMap.invalidateSize()
    };

    let showPolygon = locStorage.showPolygon !== undefined ? JSON.parse(locStorage.showPolygon) : true;
    $: locStorage.showPolygon = showPolygon;
    $: if (leafletMapPolyGroup) {setPolygonVisibility(showPolygon);}
    
    // remaining time
    let timeRemaining = undefined;
    if (sessionStorage.getItem("timeRemaining") !== undefined && 
        sessionStorage.getItem("timeRemainingGameID") !== undefined &&
        sessionStorage.getItem("roundNumber") !== undefined &&
        sessionStorage.getItem("timeRemainingGameID") == $globalResult.ChallengeResultID &&
        sessionStorage.getItem("timeRemainingRoundNumber") == $globalResult.Guesses.length
    ) {
        timeRemaining = sessionStorage.getItem("timeRemaining");
    } else {
        timeRemaining = $globalMap.TimeLimit;
    }
    sessionStorage.setItem("timeRemainingGameID", $globalResult.ChallengeResultID);
    sessionStorage.setItem("timeRemainingRoundNumber", $globalResult.Guesses.length);

    // state
    let hasGuessed = false;
    let marker = null;
    let mapFocused = !shrinkMap;
    let unfocusMapInterval;
    let showSettings = false;

    // I assume these are part of the streetview sorcery, I'm not messing with them
    let replaceStateLocal = history.replaceState;
    history.replaceState = function() {
    }

    let pushStateLocal = history.pushState;
    history.pushState = function() {
    }

    onMount(async () => {
        tileServerURL = (await $ewapi.getTileServer($globalMap.ShowLabels)).tileserver;
        totalScore = calcTotalScore($globalResult.Guesses, $globalChallenge.Places, $globalMap.GraceDistance, $globalMap.Area);
        titleInterval = setInterval(setTitle, 100);
        createMinimap();
    });

    // Sometimes, the google scripts crash on startup. Just reload the page if that happens.
    window.onerror = function(e) {
        if (e.includes("Timer")) {
            location.reload(false);
        }
    };

    function setTitle() {
        try {
            new MutationObserver(function(mutations) {
                if (document.title != "Earthwalker") {
                    document.title = "Earthwalker";
                }
            }).observe(
                document.querySelector('title'),
                { childList: true }
            );
            clearInterval(interval);
        } catch (e) {
            // Title element is not there yet.
            // Wait for the next poll...
        }
    }

    function makeGuess(latlng) {
        if (hasGuessed) {
            return;
        }
        hasGuessed = true;
        latlng = latlng.wrap();
        let guess = {
            ChallengeResultID: $globalResult.ChallengeResultID,
            RoundNum: $globalResult.Guesses.length,
            Location: {Lat: latlng.lat, Lng: latlng.lng},
        };
        $ewapi.postGuess(guess).then((response) => {
            if (response) {
                window.location.replace("/scores?id="+$globalChallenge.ChallengeID);
            } else {
                alert("Failed to submit guess?!");
            }
        });
    }

    // The leaflet minimap!
    function createMinimap() {
        leafletMap = L.map("leaflet-map");

        // Load marker if it was previously stored (see reload button)
        let oldMarker = null;
        try {
            oldMarker = JSON.parse(sessionStorage.getItem("lastMarker"));
        } finally {
            if (oldMarker != null && oldMarker.gameID == $globalResult.ChallengeResultID && oldMarker.roundNumber == $globalResult.Guesses.length) {
                marker = L.marker(L.latLng(oldMarker.lat, oldMarker.lng));
                marker.addTo(leafletMap);
                guessButton.className = guessButton.className.replace("disabled", "");
            }
            sessionStorage.removeItem("lastMarker");
        }

        // only show text labels on minimap if the user wishes so
        L.tileLayer(tileServerURL, {
            attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OSM</a> contributors, <a href="https://wikitech.wikimedia.org/wiki/Wikitech:Cloud_Services_Terms_of_use">Wikimedia Cloud Services</a>'
        }).addTo(leafletMap);

        function onMapClick(event) {
            if (marker != null) {
                leafletMap.removeControl(marker);
            }
            marker = L.marker(event.latlng);
            marker.addTo(leafletMap);
            guessButton.className = guessButton.className.replace("disabled", "");
        }

        leafletMap.on("click", onMapClick);

        leafletMapPolyGroup = L.layerGroup().addTo(leafletMap);
        let map_poly = showPolygonOnMap(leafletMapPolyGroup, $globalMap.Polygon);

        // Zoom out map
        setTimeout(function() {
            leafletMap.invalidateSize();
            if (oldMarker) {
                leafletMap.setView(oldMarker.mapCenter, oldMarker.mapZoom);
            } else if ($globalMap.Polygon) {
                leafletMap.fitBounds(map_poly.getBounds());
            } else {
                // Would use fitWorld() here, but the map is in its unfocused state, so it would just use the small map dimensions
                leafletMap.setView([0.0, 0.0], mapSizes[storedMapSize].zoom);
            }
        }, 100);

        // TODO: can we move the compass without doing this?
        // Move the compass from inside the google code to the compass container.
        let compassContainer = document.getElementById("compass-container");
        let compass = document.getElementById("compass");
        compass.parentNode.removeChild(compass);
        compassContainer.appendChild(compass);
        
        // score, round number, and timer
        // TODO: can we use an absolute timer instead of this interval?
        if ($globalMap.TimeLimit > 0) {
            timerInterval = setInterval(function() {
                timeRemaining -= 1;
                sessionStorage.setItem("timeRemaining", timeRemaining);
                if (timeRemaining <= 0) {
                    if (marker == null) {
                        makeGuess(L.latLng(0, 0));
                    } else {
                        makeGuess(marker.getLatLng());
                    }
                    clearInterval(timeRemaining);
                }
            }, 1000);
        }
    }

    function focusMap() {
        mapFocused = true;
        curMapSize = storedMapSize;
        clearInterval(unfocusMapInterval);
    }

    function releaseMap() {
        unfocusMapInterval = setInterval(() => {
            if (shrinkMap) {
                curMapSize = 1;
                mapFocused = false;
            }
        }, 800);
    }

    function setPolygonVisibility(show) {
        if (leafletMap && leafletMapPolyGroup) {
            if (show) {
                leafletMap.addLayer(leafletMapPolyGroup);
            } else {
                leafletMap.removeLayer(leafletMapPolyGroup);
            }
        }
    }

    // Prevent the escape key from being used to exit street view to a map
    window.addEventListener('keydown', (e) => {
        if (e.key === "Escape") {
            e.stopImmediatePropagation();
        }
    }, { capture: true });
</script>

<svelte:window bind:innerWidth={innerWidth} bind:innerHeight={innerHeight}/>

<style>
    #leaflet-map {
        border-radius: 0.25rem;
        width: 100%;
        height: 100%;
    }

    #leaflet-container {
        opacity: 50%;
        width: 300px;
        height: 225px;
        max-height: calc(100vh - 300px);
        max-width: calc(100vw - 60px);
        align-self: flex-end;
        margin-left: auto;
        margin-bottom: 0.5rem;
    }

    #leaflet-container * {
        pointer-events: all;
    }

    #navigation-bar {
        height: 30px;
    }

    #leaflet-container.focused {
        opacity: 100%;
    }

    #settings-button {
        margin-left: 0.5rem;
        height: 30px;
    }

    #compass-container {
        width: 200px;
        height: 100px;
        position: absolute;
        left: 30px;
        top: 30px;
        text-align: left;
    }

    .round-info-span {
        color: white;
        font-size: 25px;
        text-shadow: 2px 2px black;
        font-weight: bold;
    }

    #right-bar {
        position: absolute;
        height: calc(100% - 60px);
        min-width: 300px;
        margin: 30px;
        right: 0;
        pointer-events: none;
        display: flex;
        flex-direction: column;
    }

    #right-bar > * {
        pointer-events: all;
    }

    #top {
        flex: 0 0;
        user-select: none;
        max-width: 300px;
        margin-left: auto;
        text-align: center;
        float: right;
        color: white;
        background-color: rgba(0, 0, 0, 0.5);
        border-radius: 0.25rem;
        padding: 0.5rem;
    }

    #top button {
        margin: 0.25rem;
    }

    #middle {
        position: relative;
        flex: 10 1;
        width: 300px;
        margin-top: .5rem;
        margin-bottom: .5rem;
        margin-left: auto;
        margin-right: 0;
        border-radius: 0.25rem;
        padding: 0.5rem;
        float: right;
        clear: both;
        background-color: white;
        overflow: auto;
    }

    #middle h3 {
        text-align: center;
    }

    #close-button {
        position: absolute;
        top: 0.5rem;
        right: 1rem;
    }

    .hidden {
        visibility: hidden;
    }

    #bottom {
        flex: 1 1;
        float: right;
        clear: both;
    }

    #controls {
        width: 300px;
        margin-left: auto;
        float: right;
    }
</style>

<div id="right-bar">
    <div id="top">
        <div class="container">
            <div class="row">
                <div class="col">
                    Round
                </div>
                <div class="col">
                    {$globalResult && $globalMap ? ($globalResult.Guesses.length + 1) + " of " + $globalMap.NumRounds : "loading..."}
                </div>
            </div>
            <div class="row">
                <div class="col text-nowrap">
                    Total Points
                </div>
                <div class="col">
                    {totalScore}
                </div>
            </div>
            {#if timeRemaining > 0}
                <div class="row">
                    <div class="col">
                        Time
                    </div>
                    <div class="col">
                        {Math.floor(timeRemaining / 60).toString()}:{Math.floor(timeRemaining % 60).toString().padStart(2, '0')}
                    </div>
                </div>
            {/if}
            <div class="row justify-content-center">
                <button 
                    class="btn btn-light btn-sm"
                    on:click={() => {
                        if (marker != null) {
                            // Put marker into session storage
                            /* console.log(marker.getLatLng()); */
                            sessionStorage.setItem("lastMarker", JSON.stringify({
                                "lat": marker.getLatLng().lat,
                                "lng": marker.getLatLng().lng,
                                "roundNumber": $globalResult.Guesses.length,
                                "gameID": $globalResult.ChallengeResultID,
                                "mapCenter": leafletMap.getCenter(),
                                "mapZoom": leafletMap.getZoom(),
                            }));
                        }
                        // https://www.phpied.com/files/location-location/location-location.html
                        location = location;
                    }}
                >
                    return to start
                </button>
            </div>
        </div>
    </div>
    <div id="middle" class={showSettings ? "" : "hidden"}>
        <h3>Settings</h3>
        <a on:click={() => {showSettings = false;}} id="close-button" title="Close">×</a>
        <hr/>
        <div class="form-group">
            <div class="form-check">
                <input type="checkbox" class="form-check-input" id="show-polygon" bind:checked={showPolygon}>
                <label class="form-check-label" for="show-polygon">Show polygon on map</label>
            </div>
            <div class="form-check">
                <input type="checkbox" class="form-check-input" id="shrink-map" bind:checked={shrinkMap}>
                <label class="form-check-label" for="shrink-polygon">Shrink map when not in use</label>
            </div>
        </div>
    </div>
    <div 
        id="bottom"
        on:mouseenter={focusMap} 
        on:mouseleave={releaseMap}
    >
        <div 
            bind:this={floatingContainer} 
            id="leaflet-container"  
            class={mapFocused ? "focused" : ""}
        >
            <div id="leaflet-map"></div>
        </div>
        <div id="controls">
            <div id="settings-button" class="btn-group btn-group-sm float-right">
                <button class="btn btn-light" on:click={() => {showSettings = !showSettings;}}>
                    <img alt="Map Settings" src="/public/assets/settings.svg" style="width: auto; height: 1rem;"/>
                </button>
            </div>
            <div id="navigation-bar" class="btn-group btn-group-sm float-right">
                <button 
                    class="btn btn-light" 
                    on:click={() => {if (storedMapSize < mapSizes.length - 1) {storedMapSize++;}}}
                    disabled={storedMapSize == mapSizes.length - 1}
                >⬉</button>
                <button 
                    class="btn btn-light" 
                    on:click={() => {if (storedMapSize > 0) {storedMapSize--;}}} 
                    disabled={storedMapSize == 0}
                >⬊</button>
            </div>
            <button 
                bind:this={guessButton}
                class="btn btn-primary btn-sm float-left disabled" 
                on:click={() => {
                    if (marker == null) {
                        alert("You have to add a marker first! Do this by clicking the map.");
                        return;
                    }
                    makeGuess(marker.getLatLng());
                }}>
                Guess!
            </button>
        </div>
    </div>
</div>

<div id="compass-container"></div>

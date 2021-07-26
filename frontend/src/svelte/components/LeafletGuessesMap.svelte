<script>
    import {onMount} from 'svelte';
    import { ewapi, globalMap, globalChallenge } from '../../js/stores.js';
    import { showGuessOnMap, showPolygonOnMap } from '../../js/earthwalker';
    import L from 'leaflet';

    export let displayedResults, showAll, curRound;

    let tileServer;

    let mapDiv;

    let lMap;
    let polyGroup;
    let guessGroup;

    $: if (guessGroup && displayedResults) {
        guessGroup.clearLayers();
        if (showAll) {
            for (let displayedResult of displayedResults) {
                showGuesses(guessGroup, displayedResult.Guesses, displayedResult.Nickname, displayedResult.Icon);
            }
        } else {
            for (let displayedResult of displayedResults) {
                showGuesses(guessGroup, [displayedResult.Guesses[curRound]], displayedResult.Nickname, displayedResult.Icon);
            }
        }
        lMap.fitBounds(guessGroup.getBounds());
    };

    onMount(async () => {
        lMap = new L.Map(mapDiv);
        lMap.setView([0.0, 0.0], 1);

        tileServer = (await $ewapi.getTileServer()).tileserver;
        L.tileLayer(tileServer, {
            attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OSM</a> contributors, <a href="https://wikitech.wikimedia.org/wiki/Wikitech:Cloud_Services_Terms_of_use">Wikimedia Cloud Servides</a>'
        }).addTo(lMap);

        polyGroup = L.layerGroup().addTo(lMap);
        if ($globalMap.Polygon) {
            showPolygonOnMap(polyGroup, $globalMap.Polygon);
        }

        guessGroup = L.featureGroup().addTo(lMap);
    });

    function showGuesses(layer, guesses, nickname, icon) {
        guesses.forEach(guess => {
            showGuessOnMap(layer, guess, $globalChallenge.Places[guess.RoundNum], guess.RoundNum, nickname, icon);
        });
    }
</script>

<style>
    div {
        width: 100%;
        height: 50vh;
    }
</style>

<div bind:this={mapDiv}></div>

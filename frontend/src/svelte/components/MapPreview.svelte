<script>
    import { onMount } from "svelte";
    import { ewapi } from "../../js/stores";
    import { showPolygonOnMap } from '../../js/earthwalker';

    export let map;

    let tileServer;
    let mapDiv;

    let lMap;
    let polyGroup;

    onMount(async () => {
        lMap = new L.Map(mapDiv, {
            attributionControl: false,
            zoomControl: false
        });
        lMap.setView([0, 0], 0);

        tileServer = (await $ewapi.getTileServer()).tileserver;
        L.tileLayer(tileServer).addTo(lMap);

        polyGroup = L.featureGroup().addTo(lMap);
        if (map.Polygon) {
            showPolygonOnMap(polyGroup, map.Polygon);
            lMap.fitBounds(polyGroup.getBounds());
        }
    });
</script>

<div class="card mt-4 mx-1">
    <div 
        bind:this={mapDiv}
        id="map"
        style="min-height: 5rem; width: 100%;"
    ></div>
    <div class="card-body">
        <h5 class="card-title">{map.Name}</h5>
        <p class="card-text">
            Rounds: {map.NumRounds}<br>
            {map.TimeLimit > 0 ? `Time limit: ${Math.floor(map.TimeLimit / 60)}:${Math.floor(map.TimeLimit % 60).toString().padStart(2, '0')}` : 'No Time Limit'}<br>
        </p>
        <a href="/createchallenge?mapid={map.MapID}" class="btn btn-primary">Use Map</a>
        <a href="/editmap?mapid={map.MapID}" class="btn btn-secondary">Edit Map</a>
    </div>
</div>
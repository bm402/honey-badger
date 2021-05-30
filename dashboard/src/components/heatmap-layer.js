import { useEffect } from 'react';
import { useMap } from 'react-leaflet'
import L from 'leaflet'
import 'leaflet.heat'

const HeatmapLayer = () => {
    const map = useMap();

    useEffect(() => {
        fetch('https://omf1aavgfc.execute-api.eu-west-2.amazonaws.com/prod/v1/heatmap-data')
            .then(res => res.json())
            .then(rawData => rawData.heatmap_data_points)
            .then(rawDataPoints => {

                let sortedDataPoints = rawDataPoints
                    .map(point => [point.lat, point.lon, point.count])
                    .sort(function(a, b){ return a[2]-b[2] });

                let flattenedCount = 1
                let lastCount = sortedDataPoints[0][2]
                for (let i = 0; i < sortedDataPoints.length; i++) {
                    if (sortedDataPoints[i][2] > lastCount) {
                        flattenedCount++
                    }
                    lastCount = sortedDataPoints[i][2]
                    sortedDataPoints[i][2] = flattenedCount
                }

                const maxCount = sortedDataPoints[sortedDataPoints.length-1][2]
                const normalisedDataPoints = sortedDataPoints
                    .map(point => [point[0], point[1], point[2]/maxCount]);

                const heatmapConfig = {
                    minOpacity: 0.4,
                    maxZoom: 3,
                    radius: 20,
                    blur: 15,
                    gradient: {
                        0.4: 'blue',
                        0.6: 'lime',
                        0.8: 'yellow',
                        0.9: 'orange',
                        1.0: 'red',
                    },
                };
                
                L.heatLayer(normalisedDataPoints, heatmapConfig).addTo(map);

                let markerBounds = L.latLngBounds([]);
                normalisedDataPoints.forEach(normalisedDataPoint => {
                    markerBounds.extend([normalisedDataPoint[0], normalisedDataPoint[1]]);
                });

                map.fitBounds(markerBounds);
            })
            .catch(console.log);

    }, [map]);

    return null
}

export default HeatmapLayer

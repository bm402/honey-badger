import { useEffect } from 'react';
import { latLngBounds } from 'leaflet'
import { useMap } from 'react-leaflet'

const StatsMapLayer = props => {
    const map = useMap();

    useEffect(() => {
        if (!props.mapData || props.mapData.length === 0) {
            return;
        }

        let markerBounds = latLngBounds([]);
        props.mapData.forEach(mapDataItem => {
            markerBounds.extend([mapDataItem.lat, mapDataItem.lon]);
        });

        map.fitBounds(markerBounds, {
            maxZoom: 10,
        });
    });

    return null;
};

export default StatsMapLayer;

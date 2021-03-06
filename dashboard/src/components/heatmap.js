import { MapContainer, TileLayer } from 'react-leaflet'
import HeatmapLayer from './heatmap-layer'

import './heatmap.css';

const Heatmap = () => {
    return (
        <MapContainer className="heatmap">
            <HeatmapLayer />
            <TileLayer
                attribution='&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'
                url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
            />
        </MapContainer>
    );
};

export default Heatmap;

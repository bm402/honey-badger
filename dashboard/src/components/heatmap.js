import { MapContainer, TileLayer } from 'react-leaflet'
import './heatmap.css';

const Heatmap = () => {
    return (
        <MapContainer className="heatmap" center={[0, 0]} zoom={2}>
            <TileLayer
                attribution='&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'
                url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
            />
        </MapContainer>
    );
};

export default Heatmap;

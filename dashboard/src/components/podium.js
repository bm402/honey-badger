import { useState } from 'react'
import Card from 'react-bootstrap/Card'
import CardDeck from 'react-bootstrap/CardDeck';
import Modal from 'react-bootstrap/Modal'
import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet'
import StatsMapLayer from './stats-map-layer';

import GoldMedal from '../images/gold.jpg'
import SilverMedal from '../images/silver.jpg' 
import BronzeMedal from '../images/bronze.jpg'

import './podium.css'
import './loading.css'

const Podium = props => {

    const [showDetail, setShowDetail] = useState(false);
    const [detailPosition, setDetailPosition] = useState(0);

    const handleCloseDetail = () => setShowDetail(false);
    const handleShowDetail = position => {
        setDetailPosition(position);
        setShowDetail(true);
    };

    const createModalMetadata = metadata => {
        if (!metadata) {
            return null;
        }

        return Object.keys(metadata).map(key => (
            <p className="modal-metadata-item" key={key}>
                {metadata[key].title}: {formatModalMetadataItemValue(metadata[key].value)}
            </p>
        ));
    };

    const formatModalMetadataItemValue = value => {
        if (Array.isArray(value)) {
            return value
                .sort((a, b) => a.localeCompare(b))
                .join(', ');
        }
        return value;
    };

    const createModalMap = mapData => {
        if (!mapData || mapData.length === 0) {
            return null;
        }

        return (
            <MapContainer className="statsmap" maxZoom="18">
                <StatsMapLayer mapData={mapData} />
                <TileLayer
                    attribution='&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'
                    url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                />
                {mapData.map((mapDataItem, key) => (
                    <Marker position={[mapDataItem.lat, mapDataItem.lon]} key={key}>
                        {mapDataItem.metadata &&
                            <Popup>
                                {createModalMetadata(mapDataItem.metadata)}
                            </Popup>
                        }
                    </Marker>
                ))}
            </MapContainer>
        );
    }

    return (
        <div className="podium">
            <CardDeck>
                <Card className="podium-title-card">
                    <Card.Body className="podium-title-text">
                        <Card.Text>{props.title}</Card.Text>
                    </Card.Body>
                </Card>
                <Card className="podium-data-card" onClick={() => handleShowDetail(0)}>
                    <Card.Body>
                        <Card.Img className="podium-medal-image" src={GoldMedal} />
                        <Card.Text className={`${props.isDataLoaded ? "" : "loading"} ${"podium-text-" + props.type}`}>{props.data[0].value}</Card.Text>
                    </Card.Body>
                    <Card.Footer className="podium-gold-footer" />
                </Card>
                <Card className="podium-data-card" onClick={() => handleShowDetail(1)}>
                    <Card.Body>
                        <Card.Img className="podium-medal-image" src={SilverMedal} />
                        <Card.Text className={`${props.isDataLoaded ? "" : "loading"} ${"podium-text-" + props.type}`}>{props.data[1].value}</Card.Text>
                    </Card.Body>
                    <Card.Footer className="podium-silver-footer" />
                </Card>
                <Card className="podium-data-card" onClick={() => handleShowDetail(2)}>
                    <Card.Body>
                        <Card.Img className="podium-medal-image" src={BronzeMedal} />
                        <Card.Text className={`${props.isDataLoaded ? "" : "loading"} ${"podium-text-" + props.type}`}>{props.data[2].value}</Card.Text>
                    </Card.Body>
                    <Card.Footer className="podium-bronze-footer" />
                </Card>
            </CardDeck>
            <Modal show={showDetail} onHide={handleCloseDetail}>
                <Modal.Header closeButton>
                    <Modal.Title>{props.title}: {props.data[detailPosition].value}</Modal.Title>
                </Modal.Header>
                <Modal.Body>
                    {createModalMetadata(props.data[detailPosition].metadata)}
                    {createModalMap(props.data[detailPosition].map_data)}
                </Modal.Body>
            </Modal>
        </div>
    );
};

export default Podium;

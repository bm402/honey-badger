import CardDeck from 'react-bootstrap/CardDeck';
import Card from 'react-bootstrap/Card'

import GoldMedal from '../images/gold.jpg'
import SilverMedal from '../images/silver.jpg' 
import BronzeMedal from '../images/bronze.jpg'

import './podium.css'
import './loading.css'

const Podium = props => {

    return (
        <CardDeck className="podium">
            <Card className="podium-title-card">
                <Card.Body className="podium-title-text">
                    <Card.Text>{props.title}</Card.Text>
                </Card.Body>
            </Card>
            <Card>
                <Card.Body>
                    <Card.Img className="podium-medal-image" src={GoldMedal} />
                    <Card.Text className={`${props.isDataLoaded ? "" : "loading"} ${"podium-text-" + props.type}`}>{props.data[0].value}</Card.Text>
                </Card.Body>
                <Card.Footer className="podium-gold-footer" />
            </Card>
            <Card>
                <Card.Body>
                    <Card.Img className="podium-medal-image" src={SilverMedal} />
                    <Card.Text className={`${props.isDataLoaded ? "" : "loading"} ${"podium-text-" + props.type}`}>{props.data[1].value}</Card.Text>
                </Card.Body>
                <Card.Footer className="podium-silver-footer" />
            </Card>
            <Card>
                <Card.Body>
                    <Card.Img className="podium-medal-image" src={BronzeMedal} />
                    <Card.Text className={`${props.isDataLoaded ? "" : "loading"} ${"podium-text-" + props.type}`}>{props.data[2].value}</Card.Text>
                </Card.Body>
                <Card.Footer className="podium-bronze-footer" />
            </Card>
        </CardDeck>
    );
};

export default Podium;

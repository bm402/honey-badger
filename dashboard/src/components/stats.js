import { useEffect, useState } from 'react';

import Podium from './podium'

import './stats.css'

const Stats = () => {
    const [isDataLoaded, setIsDataLoaded] = useState(false)
    const [statsData, setStatsData] = useState({
        "most_connections": [{},{},{}],
        "most_active_cities": [{},{},{}],
        "most_active_countries": [{},{},{}],
        "most_ip_addresses": [{},{},{}],
        "most_ingress_ports": [{},{},{}],
    });

    useEffect(() => {
        fetch('https://omf1aavgfc.execute-api.eu-west-2.amazonaws.com/prod/v1/stats-data')
            .then(res => res.json())
            .then(data => {
                setStatsData(data);
                setIsDataLoaded(true);
            })
            .catch(console.log);
    }, []);

    return (
        <div className="stats-page">
            <Podium title="Most connection attempts" data={statsData.most_connections} type="number" isDataLoaded={isDataLoaded} />
            <Podium title="Most active city" data={statsData.most_active_cities} type="string" isDataLoaded={isDataLoaded} />
            <Podium title="Most active country" data={statsData.most_active_countries} type="string" isDataLoaded={isDataLoaded} />
            <Podium title="Most IP addresses used" data={statsData.most_ip_addresses} type="number" isDataLoaded={isDataLoaded} />
            <Podium title="Most ingress ports tried" data={statsData.most_ingress_ports} type="number" isDataLoaded={isDataLoaded} />
        </div>
    );
};

export default Stats

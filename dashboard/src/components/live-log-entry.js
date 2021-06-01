import { useState } from 'react'

import './live-logs.css'

const LiveLogEntry = props => {
    const [isOpen, setIsOpen] = useState(false);

    return (
        <tr className={`live-log-entry ${isOpen ? "open" : "closed"}`}>
            <td className="live-log-entry-field toggle" onClick={() => {isOpen ? setIsOpen(false) : setIsOpen(true)}}>{isOpen ? '\u25bc' : '\u25b6'}</td>
            <td className="live-log-entry-field timestamp">{new Date(props.logEntry.timestamp).toLocaleString()}</td>
            <td className="live-log-entry-field ingress-port">{props.logEntry.ingress_port}</td>
            <td className="live-log-entry-field ip-address">{props.logEntry.ip_address}</td>
            <td className="live-log-entry-field location">{props.logEntry.city + ', ' + props.logEntry.country}</td>
            <td className="live-log-entry-field input">{props.logEntry.input}</td>
        </tr>
    );
};

export default LiveLogEntry;

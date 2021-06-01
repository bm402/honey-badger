import { useState } from 'react'
import Button from 'react-bootstrap/Button'
import Table from 'react-bootstrap/Table'
import LiveLogEntry from './live-log-entry';

import './live-logs.css'

const LiveLogs = () => {
    const [websocket, setWebsocket] = useState(null);
    const [isListening, setIsListening] = useState(false);
    const [status, setStatus] = useState("Disconnected");
    const [logEntries, setLogEntries] = useState([]);

    const connectWebsocket = () => {
        const websocket = new WebSocket('wss://rgs4h7oyra.execute-api.eu-west-2.amazonaws.com/prod');
        setWebsocket(websocket);
        setIsListening(true);
        setStatus("Connecting...");

        websocket.onopen = () => {
            setStatus("Connected");
        };
    
        websocket.onmessage = (message) => {
            setLogEntries(logEntries => [JSON.parse(message.data), ...logEntries]);
        };

        websocket.onclose = () => {
            setStatus("Disconnected");
        };
    };

    const disconnectWebsocket = () => {
        websocket.close();
        setIsListening(false);
    };

    return (
        <div className="live-logs-page">
            <div className="live-logs-control-panel">
                <Button className="live-logs-control-panel-button" variant="success" disabled={isListening} onClick={!isListening ? connectWebsocket : null}>{'\u25b6'}</Button>
                <Button className="live-logs-control-panel-button" variant="danger" disabled={!isListening} onClick={isListening ? disconnectWebsocket : null}>{'\u25a0'}</Button>
                <div className="live-logs-status">Status: {status}</div>
            </div>
            <div className="live-logs-log-section">
                <Table striped bordered hover size="sm">
                    <thead>
                        <tr>
                            <th className="live-log-entry-field toggle"></th>
                            <th className="live-log-entry-field timestamp">Time</th>
                            <th className="live-log-entry-field ingress-port">Port</th>
                            <th className="live-log-entry-field ip-address">IP</th>
                            <th className="live-log-entry-field location">Location</th>
                            <th className="live-log-entry-field input">Input</th>
                        </tr>
                    </thead>
                    <tbody>
                        {logEntries.map((logEntry, key) => (
                            <LiveLogEntry logEntry={logEntry} key={key} />
                        ))}
                    </tbody>
                </Table>
            </div>
        </div>
    );
};

export default LiveLogs;

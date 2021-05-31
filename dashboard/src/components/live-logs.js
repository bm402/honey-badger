import { w3cwebsocket as W3CWebSocket } from 'websocket'

const LiveLogs = () => {
    const client = new W3CWebSocket('wss://rgs4h7oyra.execute-api.eu-west-2.amazonaws.com/prod');

    client.onopen = () => {
        console.log('WebSocket Client Connected');
    };

    client.onmessage = (message) => {
        console.log(message);
    };

    return (
        <div>Testing websockets</div>
    );
};

export default LiveLogs;

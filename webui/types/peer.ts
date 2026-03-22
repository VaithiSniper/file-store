type Peer = {
    id: string;
    name: string;
    address: string;
    apiUrl: string;
    status: 'syncing' | 'connected' | 'disconnected' | 'unknown';
    uptime: string;
    latency: string;
    dataIn: string;
    dataOut: string;
}

enum PEER_API_HEALTH_STATUS {
    OK = 'ok',
    ERROR = 'error',
}

export type { Peer }
export { PEER_API_HEALTH_STATUS }
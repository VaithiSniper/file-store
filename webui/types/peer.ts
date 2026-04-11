enum PEER_API_HEALTH_STATUS {
    HEALTHY = 'healthy',
    SYNCING = 'syncing',
    DISCONNECTED = 'disconnected',
    ERROR = 'error',
    UNKNOWN = 'unknown'
}
type Peer = {
    id: string;
    name: string;
    listen_address: string;
    api_url: string;
    status: PEER_API_HEALTH_STATUS;
    uptime: string;
    latency: string;
    data_in: string;
    data_out: string;
    path_transform_func: string;
    message_format: string;
    base_storage_location: string;
    bootstrap_nodes: string[];
}

type PeerStatus = {
    status: PEER_API_HEALTH_STATUS;
    message: string;
}
type PeerListEntry = {
    listen_address: string;
    api_url: string;
}

export type { Peer, PeerListEntry, PeerStatus }
export { PEER_API_HEALTH_STATUS }
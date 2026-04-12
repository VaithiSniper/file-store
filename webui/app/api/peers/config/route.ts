import type { NextRequest } from 'next/server'
import { Peer, PEER_API_HEALTH_STATUS } from "@/types/peer";

export async function getPeerConfig(apiUrl: string) {
    const configAPIEndpoint = `${apiUrl}/api/config`;
    let peerConfig: Peer | null = {
        id: "",
        name: "",
        listen_address: "",
        api_url: apiUrl,
        status: PEER_API_HEALTH_STATUS.UNKNOWN,
        uptime: "",
        latency: "",
        data_in: "",
        data_out: "",
        path_transform_func: "",
        message_format: "",
        base_storage_location: "",
        bootstrap_nodes: []
    }
    try {
        const resp = await fetch(configAPIEndpoint);
        const result: Peer = await resp.json();
        if (resp.ok && result) {
            peerConfig = result;
        } else {
            console.log("Failed to fetch peer config from", apiUrl);
            peerConfig = null;
        }
    } catch (e) {
        console.log("No peer found at", apiUrl);
        peerConfig = null;
    }

    return peerConfig;
}

export async function GET(
    request: NextRequest,
) {
    const { searchParams } = request.nextUrl;
    const peer_api_url = searchParams.get("peer_api_url") || "";
    const peerConfig = await getPeerConfig(peer_api_url);

    return Response.json(peerConfig);
}
import type { NextRequest } from 'next/server'
import { PEER_API_HEALTH_STATUS, PeerStatus } from "@/types/peer";

export async function pingPeerStatus(apiUrl: string) {
    const statusAPIEndpoint = `${apiUrl}/api/status`;
    let peerStatus: PeerStatus = {
        status: PEER_API_HEALTH_STATUS.UNKNOWN,
        message: "unreachable"
    };
    try {
        const resp = await fetch(statusAPIEndpoint);
        const result = await resp.json();
        if (resp.ok && result) {
            switch (result.status) {
                case PEER_API_HEALTH_STATUS.HEALTHY:
                    peerStatus.message = `Peer discovered at ${apiUrl} is healthy and reachable.`;
                    break;
                case PEER_API_HEALTH_STATUS.SYNCING:
                    peerStatus.message = `Peer discovered at ${apiUrl} is currently syncing.`;
                    break;
                case PEER_API_HEALTH_STATUS.DISCONNECTED:
                    peerStatus.message = `Peer discovered at ${apiUrl} is currently disconnected.`;
                    break;
                case PEER_API_HEALTH_STATUS.ERROR:
                    peerStatus.message = `Peer discovered at ${apiUrl} is reachable but reported an error status.`;
                    break;
                default:
                    peerStatus.message = `Peer discovered at ${apiUrl} is reachable but returned an unknown status.`;
                    peerStatus.status = PEER_API_HEALTH_STATUS.UNKNOWN;
            }
            peerStatus.status = result.status;
        }
    } catch (e) {
        console.log("No peer found at", apiUrl);
    }

    return peerStatus;
}

export async function GET(
    request: NextRequest,
) {
    const { searchParams } = request.nextUrl;
    const peer_api_url = searchParams.get("peer_api_url") || "";
    const peerStatus = await pingPeerStatus(peer_api_url);

    return Response.json(peerStatus);
}

import type { NextRequest } from 'next/server'
import { PEER_API_HEALTH_STATUS, PeerListEntry } from "@/types/peer";

export async function getPeerList(apiUrl: string) {
    const peerListAPIEndpoint = `${apiUrl}/api/peers`;
    let peerList: PeerListEntry[] = [];
    try {
        const resp = await fetch(peerListAPIEndpoint);
        const result = await resp.json();
        if (resp.ok && result) {
            peerList = result.peers.map((peerAddr: string) => (
                {
                    listen_address: peerAddr,
                    api_url: "",
                }
            ));
        }
    } catch (e) {
        console.log("No peer list found for", apiUrl);
    }

    return peerList;
}

export async function GET(
    request: NextRequest,
) {
    const { searchParams } = request.nextUrl;
    const peer_api_url = searchParams.get("peer_api_url") || "";
    const peerList = await getPeerList(peer_api_url);

    return Response.json(peerList);
}

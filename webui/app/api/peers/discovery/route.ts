import { Peer, PEER_API_HEALTH_STATUS } from "@/types/peer";

function validateUrl(address: string) {
    try {
        new URL(address);
        return true;
    } catch (e) {
        return false;
    }
}

async function pingPeerStatus(apiUrl: string) {
    let peer: Peer = {
        id: "",
        name: "",
        address: "",
        apiUrl: "",
        status: "unknown",
        uptime: "",
        latency: "",
        dataIn: "",
        dataOut: ""
    };

    try {
        const resp = await fetch(apiUrl);
        const result = await resp.json();
        if (resp.ok && result && result.status === PEER_API_HEALTH_STATUS.OK) {
            console.log("Peer auto-discovered at", apiUrl);
            // Hit other endpoints to get the peer info and fill in the details.
            peer.apiUrl = apiUrl;
        }
    } catch (e) {
        console.log("No peer found at", apiUrl);
    }

    return peer;
}

async function autoDiscoverPeers(apiUrl: string) {
    let discoveredPeers: Peer[] = [];
    if (!apiUrl || apiUrl.trim() === "" || !validateUrl(apiUrl)) {
        console.log("No API URL provided for peer discovery.");
        return discoveredPeers;
    }
    apiUrl = apiUrl.trim() + "/api/status";

    const peer = await pingPeerStatus(apiUrl);
    if (peer.apiUrl !== "") {
        // Add first known peer
        discoveredPeers.push(peer);
    }

    // Now, get the list of known peers from the first discovered peer and ping them as well.
    if (discoveredPeers.length > 0) {
        try {
            const resp = await fetch(`${discoveredPeers[0].apiUrl}/api/peers`);
            const result = await resp.json();
            if (resp.ok && result && Array.isArray(result.peers)) {
                for (const peerInfo of result.peers) {
                    const peerStatus = await pingPeerStatus(peerInfo.apiUrl);
                    if (peerStatus.apiUrl !== "") {
                        discoveredPeers.push(peerStatus);
                    }
                }
            }
        } catch (e) {
            console.log("Failed to fetch peers from the first discovered peer.");
        }
    }

    // If we can't find any peer, we return an empty list.
    return discoveredPeers;
}

export async function POST(request: Request) {
    const { apiUrl } = await request.json();

    const discoveredPeers = await autoDiscoverPeers(apiUrl);

    return Response.json({ peers: discoveredPeers });
}
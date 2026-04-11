import { Peer, PEER_API_HEALTH_STATUS, PeerListEntry, PeerStatus } from "@/types/peer";
import { pingPeerStatus } from "../status/route";
import { getPeerConfig } from "../config/route";
import { getPeerList } from "../peer_list/route";

function validateUrl(address: string) {
    try {
        new URL(address);
        return true;
    } catch (e) {
        return false;
    }
}

function isValidPeerStatus(peerStatus: PeerStatus) {
    return peerStatus && peerStatus.status !== PEER_API_HEALTH_STATUS.UNKNOWN && peerStatus.message !== "unreachable";
}

async function checkPeerHealthAndGetConfig(apiUrl: string) {
    // First, we try to ping the provided API URL to see if it's valid and healthy.
    let peer: Peer | null = null;
    const peerStatus = await pingPeerStatus(apiUrl);
    if (!isValidPeerStatus(peerStatus)) {
        console.log("No peer found at the provided API URL.");
        return peer;
    }
    console.log(`Peer discovered at provided API URL=${apiUrl} is ${peerStatus.status}`);

    // Get peer details
    peer = await getPeerConfig(apiUrl);
    if (peer && peer.id === "" && peer.name === "") {
        console.log("Invalid peer config received.");
        return peer;
    }

    console.log("Peer config fetched successfully:", peer);
    return peer;
}

async function autoDiscoverPeers(apiUrl: string) {
    let discoveredPeers: Peer[] = [];
    let peer: Peer | null = null;
    let peerList: PeerListEntry[] = [];

    apiUrl = apiUrl.trim();
    if (!apiUrl || apiUrl === "" || !validateUrl(apiUrl)) {
        console.log("No API URL provided for peer discovery.");
        return discoveredPeers;
    }

    peer = await checkPeerHealthAndGetConfig(apiUrl);
    if (!peer) {
        console.log("Failed to discover a valid peer at the provided API URL.");
        return discoveredPeers;
    }
    // We found a peer at the provided API URL, add it to the list.
    discoveredPeers.push(peer);

    // Get the peerlist from this peer.
    peerList = await getPeerList(apiUrl);

    for (const peerInfo of peerList) {
        if (!peerInfo.api_url || peerInfo.api_url === "" || !validateUrl(peerInfo.api_url)) {
            console.log("Invalid API URL in peer list, skipping:", peerInfo.api_url);
            continue;
        }
        peer = await checkPeerHealthAndGetConfig(peerInfo.api_url);
        if (!peer) {
            console.log(`Failed to get peer %s from peer list:`, peerInfo.api_url);
            continue;
        }
        discoveredPeers.push(peer);
    }

    console.log(`Peer discovery completed. Total peers discovered: ${discoveredPeers.length}.`);
    return discoveredPeers;
}

export async function POST(request: Request) {
    const { apiUrl } = await request.json();

    const discoveredPeers = await autoDiscoverPeers(apiUrl);

    return Response.json({ peers: discoveredPeers });
}
import type { Peer } from "@/types/peer"

export async function GET() {
    const baseUrl = process.env.PEER_2_API_SERVER_URL || 'http://localhost:9001';
    let peerList: Peer[] = [];


    const resp = await fetch(`${baseUrl}/api/peers`);
    const result = await resp.json();
    if (!resp.ok || !result || !result.peers) {
        console.error("Failed to fetch peers from peer 2:", result);
        return new Response("Failed to fetch peers", { status: 500 });
    }

    if (result.peers.length >= 0) {
        peerList = result.peers.map((peerAddr: any) => ({
            address: peerAddr,
            id: '',
            name: '',
            status: 'connected',
            uptime: '',
            latency: '',
            dataIn: '',
            dataOut: '',
        }));

    }
    console.log("Peer list fetched from peer 2:", peerList);

    return Response.json(peerList);
}
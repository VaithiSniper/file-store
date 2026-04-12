import { Peer } from '@/types/peer';
import { Stat } from '@/types/stats';
import { create } from 'zustand';

type PeersStoreState = {
    peerList: Peer[];
    peersCount: number;
    discoveryDone: boolean;
};

type PeersStoreActions = {
    setPeers: (peers: Peer[]) => void;
    addPeerToList: (peer: Peer) => void;
    updatePeerNameById: (peerId: string, newPeerName: string) => void;
    removePeerById: (peerId: string) => void;
    clearPeerList: () => void;
    mergePeerLists: (newPeers: Peer[]) => void;
    setDiscoveryDone: (done: boolean) => void;
};

type PeersStore = PeersStoreState & PeersStoreActions;

const usePeerStore = create<PeersStore>((set) => ({
    peerList: [],
    peersCount: 0,
    discoveryDone: false,
    setPeers: (peers) => set({ peerList: peers, peersCount: peers.length }),
    addPeerToList: (peer) => set((state) => ({ peerList: [...state.peerList, peer], peersCount: state.peersCount + 1 })),
    updatePeerNameById: (peerId: string, newPeerName: string) => set((state) => {
        const peer = state.peerList.find((p) => p.id === peerId);
        if (!peer) return state;
        return {
            peerList: state.peerList.map((p) => (p.id === peerId ? { ...peer, name: newPeerName } : p)),
            peersCount: state.peersCount,
            discoveryDone: state.discoveryDone
        }
    }),
    removePeerById: (peerId: string) => set((state) => ({ peerList: state.peerList.filter((p) => p.id !== peerId), peersCount: state.peersCount - 1 })),
    clearPeerList: () => set({ peerList: [], peersCount: 0, discoveryDone: false }),
    mergePeerLists: (newPeers: Peer[]) => set((state) => {
        const existingPeerIds = new Set(state.peerList.map((p) => p.id));
        const mergedPeers = [...state.peerList];
        newPeers.forEach((peer) => {
            if (!existingPeerIds.has(peer.id)) {
                mergedPeers.push(peer);
            }
        });
        return { peerList: mergedPeers, peersCount: mergedPeers.length, discoveryDone: state.discoveryDone };
    }),
    setDiscoveryDone: (done: boolean) => set({ discoveryDone: done }),
}));

type HyperstoreStats = {
    stats: Stat[]
}

type HyperstoreStatsActions = {
    setStats: (stats: Stat[]) => void;
    updateStatByLabel: (stat: Stat) => void;
    clearStats: () => void;
};

type HyperstoreStatsStore = HyperstoreStats & HyperstoreStatsActions;

const useHyperstoreStatsStore = create<HyperstoreStatsStore>()((set) => ({
    stats: [],
    setStats: (stats) => set({ stats }),
    updateStatByLabel: (stat) => set((state) => ({
        stats: state.stats.map((s) => (s.label === stat.label ? { ...s, ...stat } : s)),
    })),
    clearStats: () => set({ stats: [] }),
}));

export { usePeerStore, useHyperstoreStatsStore };
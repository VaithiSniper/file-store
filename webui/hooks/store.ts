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
    updatePeerById: (peer: Peer) => void;
    removePeerFromList: (peer: Peer) => void;
    clearPeerList: () => void;
    setDiscoveryDone: (done: boolean) => void;
};

type PeersStore = PeersStoreState & PeersStoreActions;

const usePeerStore = create<PeersStore>()((set) => ({
    peerList: [],
    peersCount: 0,
    discoveryDone: false,
    setPeers: (peers) => set({ peerList: peers, peersCount: peers.length, discoveryDone: true }),
    addPeerToList: (peer) => set((state) => ({ peerList: [...state.peerList, peer], peersCount: state.peersCount + 1, discoveryDone: state.discoveryDone })),
    updatePeerById: (peer) => set((state) => ({
        peerList: state.peerList.map((p) => (p.id === peer.id ? { ...p, ...peer } : p)),
        peersCount: state.peersCount,
        discoveryDone: state.discoveryDone
    })),
    removePeerFromList: (peer) => set((state) => ({ peerList: state.peerList.filter((p) => p !== peer), peersCount: state.peersCount - 1, discoveryDone: state.discoveryDone })),
    clearPeerList: () => set({ peerList: [], peersCount: 0, discoveryDone: false }),
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
'use client'

import { StatCardsRow } from "@/components/dashboard/stat-cards-row"
import { NetworkChart } from "@/components/dashboard/network-chart"
import { StorageChart } from "@/components/dashboard/storage-chart"
import { PeersList } from "@/components/dashboard/peers-list"
import { FilesTable } from "@/components/dashboard/files-table"
import { NodeInfo } from "@/components/dashboard/node-info"
import { ActivityFeed } from "@/components/dashboard/activity-feed"
import { Button } from "@/components/ui/button"
import { useState } from "react"
import { Peer } from "@/types/peer"
import { useFilesStore, usePeerStore } from "@/hooks/store"
import { FileType } from "@/types/file"



export default function DashboardPage() {
  const [apiUrl, setApiUrl] = useState("http://localhost:9000");
  const peers = usePeerStore((state) => state.peerList);
  const setPeers = usePeerStore((state) => state.setPeers);
  const discoveryDone = usePeerStore((state) => state.discoveryDone);
  const setDiscoveryDone = usePeerStore((state) => state.setDiscoveryDone);
  const files = useFilesStore((state) => state.fileList);
  const setFiles = useFilesStore((state) => state.setFiles);

  async function discoverPeers(apiUrl: string) {
    const baseUrl = process.env.NEXTJS_API_BASE_URL || '';

    const resp: Response = await fetch(`${baseUrl}/api/peers/discovery`,
      {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ apiUrl }),
      }
    );
    const peers: { peers: Peer[] } = await resp.json();
    console.log("Discovered peers:", peers.peers);
    setPeers(peers.peers);
    const respFiles: Response = await fetch(`${baseUrl}/api/peers/files?peer_api_url=${apiUrl}`);
    const files: { files: FileType[] } = await respFiles.json();
    console.log("Discovered files:", files.files);
    setFiles(files.files);
  }

  return (
    <main className="flex-1 overflow-auto">
      <div className="p-6 space-y-6">

        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <h1 className="text-2xl font-semibold tracking-tight">Overview</h1>
            {
              discoveryDone &&
              <p className="text-sm text-muted-foreground">
                Monitor your distributed file storage network
              </p>
            }
          </div>
          <div className="flex items-center gap-2 text-xs text-muted-foreground">
            {
              discoveryDone &&
              <>
                <span className="h-2 w-2 rounded-full bg-success animate-pulse" />
                All systems operational
              </>
            }
          </div>
        </div>
        {
          !discoveryDone && peers.length === 0 ?
            <>
              <div className="flex flex-col items-center justify-center gap-10 py-20 max-w-lg mx-auto">
                <p className="text-4xl text-muted-foreground">No peers discovered yet</p>
                <form className="w-full border border-gray-300 rounded-md p-10 space-y-8 flex flex-col" onSubmit={(e) => {
                  e.preventDefault();
                  discoverPeers(apiUrl);
                  setDiscoveryDone(true);
                }}>
                  <label className="block text-sm font-medium text-muted-foreground mb-2" htmlFor="apiUrl">
                    API URL of a peer
                  </label>
                  <input
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    type="text"
                    id="apiUrl"
                    placeholder="http://localhost:9000"
                    value={apiUrl}
                    onChange={(e) => setApiUrl(e.target.value)}
                  />
                  <Button variant="default" size="8xl" type="submit">
                    Start Discovery
                  </Button>
                </form>
              </div>
            </>
            :
            <>
              <StatCardsRow />

              <div className="grid gap-6 lg:grid-cols-2">
                <NetworkChart />
                <StorageChart />
              </div>

              <PeersList peers={peers} />

              <div className="grid gap-6 xl:grid-cols-3">
                <div className="xl:col-span-2">
                  <FilesTable fileList={files} />
                </div>
                <div className="space-y-6">
                  <NodeInfo />
                  <ActivityFeed />
                </div>
              </div>
            </>
        }
      </div>
    </main >
  )
}

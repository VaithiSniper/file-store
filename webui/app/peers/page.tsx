'use client'

import { PeersList } from "@/components/dashboard/peers-list"
import { usePeerStore } from "@/hooks/store";

export default function PeersPage() {
  const peers = usePeerStore((state) => state.peerList)

  return (
    <main className="flex-1 overflow-auto">
      <div className="p-6 space-y-6">
        {/* Page Header */}
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <h1 className="text-2xl font-semibold tracking-tight">Peers</h1>
            <p className="text-sm text-muted-foreground">
              View and manage connected peers in your distributed file storage network
            </p>
          </div>
        </div>
        <PeersList peers={peers} />
      </div>
    </main>
  )
}

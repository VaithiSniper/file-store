"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { MoreHorizontal, ArrowUpRight, ArrowDownLeft, Server, EditIcon, DeleteIcon, Edit, Delete, Trash, LogsIcon, PlusIcon } from "lucide-react"
import { Peer } from "@/types/peer"
import { useState } from "react"
import { EditPeerModal } from "../peers/edit-peer-modal"
import { usePeerStore } from "@/hooks/store"
import { DeleteConfirmationModal } from "../ui/delete-confirm-modal"
import { AddPeerModal } from "../peers/add-peer-modal"

export function PeersList({ peers }: { peers: Peer[] }) {
  const [selectedPeer, setSelectedPeer] = useState<Peer | null>(null);
  const [showAddPeerModal, setShowAddPeerModal] = useState(false);
  const [showEditPeerModal, setShowEditPeerModal] = useState(false);
  const [showDeletePeerModal, setShowDeletePeerModal] = useState(false);
  const mergePeerLists = usePeerStore((state) => state.mergePeerLists);
  const updatePeerNameById = usePeerStore((state) => state.updatePeerNameById);
  const removePeerById = usePeerStore((state) => state.removePeerById);

  function handleEditPeer(peer: Peer) {
    setSelectedPeer(peer);
    setShowEditPeerModal(true);
  }

  function handleDeletePeer(peer: Peer) {
    setSelectedPeer(peer);
    setShowDeletePeerModal(true);
  }

  function handlePeerLogs(peer: Peer) {
    alert(`Showing logs for peer ${peer.name}`);
  }

  async function handleAddPeerEvent(peerAPIURL: string) {
    const resp = await fetch('/api/peers/discovery', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ apiUrl: peerAPIURL }),
    });

    if (!resp.ok) {
      const errorData = await resp.json();
      throw new Error(errorData.message || 'Failed to add peer');
    }

    const result = await resp.json();
    if (result && result.peers && Array.isArray(result.peers) && result.peers.length > 0) {
      // Since some of these peers might already exist, let's merge them
      mergePeerLists(result.peers);
    }
  }

  return (
    <Card className="border-border bg-card">
      <CardHeader className="flex flex-row items-center justify-between pb-4">
        <div className="flex items-center gap-2">
          <Server className="h-4 w-4 text-muted-foreground" />
          <CardTitle className="text-sm font-medium">Connected Peers</CardTitle>
        </div>
        <div className="flex gap-4">
          <Button variant="ghost" size="sm" className="text-xs text-muted-foreground" onClick={() => setShowAddPeerModal(true)}>
            Add nodes <PlusIcon className="h-3 w-3" />
          </Button>
        </div>
      </CardHeader>
      <CardContent className="pt-0">
        <div className="space-y-3">
          {
            peers.length == 0 ?
              <div className="flex flex-col items-center justify-center gap-2 py-10">
                <Server className="h-6 w-6 text-muted-foreground" />
                <p className="text-sm text-muted-foreground">No peers connected</p>
              </div>
              :
              <table className="w-full text-left">
                <thead className="text-xs uppercase text-muted-foreground border-b-2 border-muted-foreground/50 px-2 py-4">
                  <tr>
                    <th className="px-4 py-4">Status</th>
                    <th className="px-4 py-4">Peer ID</th>
                    <th className="px-4 py-4">Addresses</th>
                    <th className="px-4 py-4">Latency</th>
                    <th className="px-4 py-4">Data In/Out</th>
                    <th className="px-4 py-4">Actions</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-muted-foreground/50">
                  {
                    peers.map((peer) => (
                      <tr
                        key={peer.id}
                        className="hover:bg-muted-foreground/10 transition-colors rounded-md border-transparent hover:border-muted-foreground/50"
                      >
                        <td className="px-4 py-4">
                          <div className="flex items-center gap-4">
                            <div className="flex h-8 w-8 items-center justify-center rounded-full bg-primary/10">
                              <div
                                className={`h-2 w-2 rounded-full ${peer.status === "healthy"
                                  ? "bg-success animate-pulse"
                                  : peer.status === "syncing"
                                    ? "bg-warning animate-pulse"
                                    : "bg-muted-foreground"
                                  }`}
                              />
                            </div>
                            <Badge
                              variant={
                                peer.status === "healthy"
                                  ? "default"
                                  : peer.status === "syncing"
                                    ? "secondary"
                                    : "outline"
                              }
                              className={
                                peer.status === "healthy"
                                  ? "bg-success/10 text-success hover:bg-success/20 border-success/30"
                                  : peer.status === "syncing"
                                    ? "bg-warning/10 text-warning hover:bg-warning/20 border-warning/30"
                                    : ""
                              }
                            >
                              {peer.status}
                            </Badge>
                          </div>
                        </td>
                        <td className="px-4 py-4">
                          <p className="text-sm font-medium">{peer.name}</p>
                          <p className="text-xs text-muted-foreground font-mono">
                            {
                              peer.id ? `${peer.id.substring(0, 12)}...` : `ID not available`
                            }
                          </p>
                        </td>
                        <td className="px-4 py-4">
                          <p className="text-sm font-medium">
                            API: {peer.listen_address}
                          </p>
                          <p className="text-xs text-muted-foreground font-mono">
                            Peering: {peer.api_url}
                          </p>
                        </td>
                        <td className="px-4 py-4">
                          <p className="font-medium font-mono">{peer.latency}</p>
                        </td>
                        <td className="px-4 py-4">
                          <div className="flex items-center gap-1 text-success">
                            <ArrowDownLeft className="h-3 w-3" />
                            <span className="font-mono">{peer.data_in}</span>
                          </div>
                          <div className="flex items-center gap-1 text-chart-2">
                            <ArrowUpRight className="h-3 w-3" />
                            <span className="font-mono">{peer.data_out}</span>
                          </div>
                        </td>
                        <td className="px-4 py-4">
                          <div className="flex w-fit items-center gap-4 bg-[#131212] rounded-md p-2">
                            <Button variant="outline" size="icon" className="h-7 w-7" onClick={() => handlePeerLogs(peer)}>
                              <LogsIcon className="h-4 w-4" />
                            </Button>
                            <Button variant="outline" size="icon" className="h-7 w-7" onClick={() => handleEditPeer(peer)}>
                              <Edit className="h-4 w-4" />
                            </Button>
                            <Button variant="destructive" size="icon" className="h-7 w-7" onClick={() => handleDeletePeer(peer)}>
                              <Trash className="h-4 w-4" />
                            </Button>
                          </div>
                        </td>
                      </tr>
                    ))
                  }
                </tbody>
              </table>
          }
          {
            showAddPeerModal && (
              <AddPeerModal
                open={showAddPeerModal}
                onOpenChange={setShowAddPeerModal}
                onSubmit={handleAddPeerEvent}
              />
            )
          }
          {
            selectedPeer && showEditPeerModal && (
              <EditPeerModal
                peerId={selectedPeer.id}
                peerAddress={selectedPeer.listen_address}
                peerBandwidth={selectedPeer.data_out}
                peerLatency={Number(selectedPeer.latency)}
                peerName={selectedPeer.name}
                open={showEditPeerModal}
                onOpenChange={setShowEditPeerModal}
                onSubmit={updatePeerNameById}
              />
            )
          }
          {
            selectedPeer && showDeletePeerModal && (
              <DeleteConfirmationModal
                itemId={selectedPeer.id}
                itemName={selectedPeer.name}
                itemType="peer"
                title="Delete Peer"
                open={showDeletePeerModal}
                onOpenChange={setShowDeletePeerModal}
                onSubmit={removePeerById}
              />
            )
          }
        </div>
      </CardContent>
    </Card >
  )
}

"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { MoreHorizontal, ArrowUpRight, ArrowDownLeft, Server } from "lucide-react"
import { Peer } from "@/types/peer"

export function PeersList({ peers }: { peers: Peer[] }) {
  return (
    <Card className="border-border bg-card">
      <CardHeader className="flex flex-row items-center justify-between pb-4">
        <div className="flex items-center gap-2">
          <Server className="h-4 w-4 text-muted-foreground" />
          <CardTitle className="text-sm font-medium">Connected Peers</CardTitle>
        </div>
        <Button variant="ghost" size="sm" className="text-xs text-muted-foreground">
          View All
        </Button>
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
              peers.map((peer) => (
                <div
                  key={peer.id}
                  className="flex items-center justify-between rounded-lg border border-border bg-muted/30 p-3"
                >
                  <div className="flex items-center gap-3">
                    <div className="flex h-8 w-8 items-center justify-center rounded-full bg-primary/10">
                      <div
                        className={`h-2 w-2 rounded-full ${peer.status === "connected"
                          ? "bg-success"
                          : peer.status === "syncing"
                            ? "bg-warning animate-pulse"
                            : "bg-muted-foreground"
                          }`}
                      />
                    </div>
                    <div>
                      <p className="text-sm font-medium">{peer.name}</p>
                      <p className="text-xs text-muted-foreground font-mono">
                        {peer.address}
                      </p>
                    </div>
                  </div>
                  <div className="hidden sm:flex items-center gap-6 text-xs">
                    <div className="text-right">
                      <p className="text-muted-foreground">Latency</p>
                      <p className="font-medium font-mono">{peer.latency}</p>
                    </div>
                    <div className="flex items-center gap-2">
                      <div className="flex items-center gap-1 text-success">
                        <ArrowDownLeft className="h-3 w-3" />
                        <span className="font-mono">{peer.dataIn}</span>
                      </div>
                      <div className="flex items-center gap-1 text-chart-2">
                        <ArrowUpRight className="h-3 w-3" />
                        <span className="font-mono">{peer.dataOut}</span>
                      </div>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    <Badge
                      variant={
                        peer.status === "connected"
                          ? "default"
                          : peer.status === "syncing"
                            ? "secondary"
                            : "outline"
                      }
                      className={
                        peer.status === "connected"
                          ? "bg-success/10 text-success hover:bg-success/20 border-success/30"
                          : peer.status === "syncing"
                            ? "bg-warning/10 text-warning hover:bg-warning/20 border-warning/30"
                            : ""
                      }
                    >
                      {peer.status}
                    </Badge>
                    <Button variant="ghost" size="icon" className="h-7 w-7">
                      <MoreHorizontal className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              ))
          }
        </div>
      </CardContent>
    </Card>
  )
}

"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Server, Cpu, HardDrive, Wifi, Clock, Globe } from "lucide-react"

export function NodeInfo() {
  return (
    <Card className="border-border bg-card">
      <CardHeader className="flex flex-row items-center justify-between pb-4">
        <div className="flex items-center gap-2">
          <Server className="h-4 w-4 text-muted-foreground" />
          <CardTitle className="text-sm font-medium">Local Node</CardTitle>
        </div>
        <Badge className="bg-success/10 text-success hover:bg-success/20 border-success/30">
          Online
        </Badge>
      </CardHeader>
      <CardContent className="pt-0 space-y-4">
        <div className="space-y-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <Globe className="h-3.5 w-3.5" />
              <span>Node ID</span>
            </div>
            <span className="text-xs font-mono text-foreground">
              QmYw...x8Kf
            </span>
          </div>
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <Wifi className="h-3.5 w-3.5" />
              <span>IP Address</span>
            </div>
            <span className="text-xs font-mono text-foreground">
              192.168.1.100
            </span>
          </div>
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <Clock className="h-3.5 w-3.5" />
              <span>Uptime</span>
            </div>
            <span className="text-xs font-mono text-foreground">
              14d 7h 23m
            </span>
          </div>
        </div>
        
        <div className="h-px bg-border" />
        
        <div className="space-y-3">
          <div>
            <div className="flex items-center justify-between mb-1.5">
              <div className="flex items-center gap-2 text-xs text-muted-foreground">
                <Cpu className="h-3 w-3" />
                <span>CPU Usage</span>
              </div>
              <span className="text-xs font-medium">23%</span>
            </div>
            <div className="h-1.5 rounded-full bg-muted overflow-hidden">
              <div className="h-full w-[23%] rounded-full bg-chart-1" />
            </div>
          </div>
          
          <div>
            <div className="flex items-center justify-between mb-1.5">
              <div className="flex items-center gap-2 text-xs text-muted-foreground">
                <HardDrive className="h-3 w-3" />
                <span>Memory</span>
              </div>
              <span className="text-xs font-medium">2.1 / 8 GB</span>
            </div>
            <div className="h-1.5 rounded-full bg-muted overflow-hidden">
              <div className="h-full w-[26%] rounded-full bg-chart-3" />
            </div>
          </div>
          
          <div>
            <div className="flex items-center justify-between mb-1.5">
              <div className="flex items-center gap-2 text-xs text-muted-foreground">
                <HardDrive className="h-3 w-3" />
                <span>Disk I/O</span>
              </div>
              <span className="text-xs font-medium">45 MB/s</span>
            </div>
            <div className="h-1.5 rounded-full bg-muted overflow-hidden">
              <div className="h-full w-[45%] rounded-full bg-chart-2" />
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

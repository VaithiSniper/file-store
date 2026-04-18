'use client'

import { DashboardHeader } from "@/components/dashboard/header"
import { NetworkChart } from "@/components/dashboard/network-chart"
import { StorageChart } from "@/components/dashboard/storage-chart"
import { PeersList } from "@/components/dashboard/peers-list"
import { FilesTable } from "@/components/dashboard/files-table"
import { NodeInfo } from "@/components/dashboard/node-info"
import { ActivityFeed } from "@/components/dashboard/activity-feed"
import { useFilesStore } from "@/hooks/store"

export default function DashboardPage() {
  const files = useFilesStore((state) => state.fileList);
  const setFiles = useFilesStore((state) => state.setFiles);

  return (
    <main className="flex-1 overflow-auto">
      <div className="p-6 space-y-6">
        {/* Page Header */}
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <h1 className="text-2xl font-semibold tracking-tight">Files</h1>
            <p className="text-sm text-muted-foreground">
              View and manage files in your distributed file storage network
            </p>
          </div>
        </div>
        <StorageChart />
        <FilesTable fileList={files} />
      </div>
    </main>
  )
}

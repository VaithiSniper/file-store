"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import {
  FolderOpen,
  Search,
  FileText,
  FileImage,
  FileArchive,
  FileCode,
  Film,
  MoreHorizontal,
} from "lucide-react"
import { FileSyncStatus, FileType } from "@/types/file"

const fileIcons: Record<string, React.ComponentType<{ className?: string }>> = {
  document: FileText,
  image: FileImage,
  archive: FileArchive,
  code: FileCode,
  video: Film,
}

export function FilesTable(props: { fileList: FileType[] }) {
  const { fileList } = props;
  const maxReplicaCount = 5; // TODO: Get quorum from the store 
  const replicaCount = 4; // TODO: Get actual replica count from the store 

  return (
    <Card className="border-border bg-card">
      <CardHeader className="flex flex-row items-center justify-between pb-4">
        <div className="flex items-center gap-2">
          <FolderOpen className="h-4 w-4 text-muted-foreground" />
          <CardTitle className="text-sm font-medium">Recent Files</CardTitle>
        </div>
        <div className="flex items-center gap-2">
          <div className="relative">
            <Search className="absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-muted-foreground" />
            <Input
              placeholder="Search files..."
              className="h-8 w-50 pl-8 text-xs bg-muted/50 border-border"
            />
          </div>
          <Button size="sm" className="h-8 text-xs">
            Upload
          </Button>
        </div>
      </CardHeader>
      <CardContent className="pt-0">
        <div className="rounded-md border border-border overflow-hidden">
          <Table>
            <TableHeader>
              <TableRow className="hover:bg-transparent border-border">
                <TableHead className="text-xs font-medium text-muted-foreground">
                  Name
                </TableHead>
                <TableHead className="text-xs font-medium text-muted-foreground">
                  Size
                </TableHead>
                <TableHead className="text-xs font-medium text-muted-foreground">
                  Replicas
                </TableHead>
                <TableHead className="text-xs font-medium text-muted-foreground">
                  Status
                </TableHead>
                <TableHead className="text-xs font-medium text-muted-foreground">
                  Modified
                </TableHead>
                <TableHead className="w-10"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {fileList.map((file) => {
                const Icon = FileText
                return (
                  <TableRow key={file.fullPath} className="border-border hover:bg-muted/30">
                    <TableCell className="font-medium">
                      <div className="flex items-center gap-2">
                        <Icon className="h-4 w-4 text-muted-foreground" />
                        <span className="text-sm truncate max-w-50">
                          {file.keyPath}
                        </span>
                      </div>
                    </TableCell>
                    <TableCell className="text-sm text-muted-foreground font-mono">
                      {file.size}
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-1">
                        {Array.from({ length: maxReplicaCount }).map((_, i) => (
                          <div
                            key={i}
                            className={`h-1.5 w-1.5 rounded-full ${i < replicaCount ? "bg-primary" : "bg-muted"
                              }`}
                          />
                        ))}
                        <span className="ml-1 text-xs text-muted-foreground">
                          {replicaCount}
                        </span>
                      </div>
                    </TableCell>
                    <TableCell>
                      <Badge
                        variant="outline"
                        className={
                          file.syncStatus === FileSyncStatus.Synced
                            ? "border-success/30 bg-success/10 text-success"
                            : file.syncStatus === FileSyncStatus.SyncingWithDataMessage || file.syncStatus === FileSyncStatus.SyncingWithStreaming
                              ? "border-warning/30 bg-warning/10 text-warning"
                              : "border-destructive/30 bg-destructive/10 text-destructive"
                        }
                      >
                        {file.syncStatus}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-sm text-muted-foreground">
                      {file.updatedAt}
                    </TableCell>
                    <TableCell>
                      <Button variant="ghost" size="icon" className="h-7 w-7">
                        <MoreHorizontal className="h-4 w-4" />
                      </Button>
                    </TableCell>
                  </TableRow>
                )
              })}
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>
  )
}

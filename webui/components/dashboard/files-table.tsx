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

const files = [
  {
    id: "file-1",
    name: "project-backup-2024.tar.gz",
    type: "archive",
    size: "2.4 GB",
    replicas: 4,
    status: "healthy",
    modified: "2 hours ago",
  },
  {
    id: "file-2",
    name: "annual-report.pdf",
    type: "document",
    size: "15.2 MB",
    replicas: 3,
    status: "healthy",
    modified: "5 hours ago",
  },
  {
    id: "file-3",
    name: "hero-banner.png",
    type: "image",
    size: "4.8 MB",
    replicas: 2,
    status: "syncing",
    modified: "1 day ago",
  },
  {
    id: "file-4",
    name: "app-source-v2.zip",
    type: "code",
    size: "890 MB",
    replicas: 5,
    status: "healthy",
    modified: "2 days ago",
  },
  {
    id: "file-5",
    name: "training-video.mp4",
    type: "video",
    size: "1.2 GB",
    replicas: 1,
    status: "at-risk",
    modified: "3 days ago",
  },
  {
    id: "file-6",
    name: "database-dump.sql",
    type: "code",
    size: "340 MB",
    replicas: 4,
    status: "healthy",
    modified: "4 days ago",
  },
]

const fileIcons: Record<string, React.ComponentType<{ className?: string }>> = {
  document: FileText,
  image: FileImage,
  archive: FileArchive,
  code: FileCode,
  video: Film,
}

export function FilesTable() {
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
              className="h-8 w-[200px] pl-8 text-xs bg-muted/50 border-border"
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
                <TableHead className="w-[40px]"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {files.map((file) => {
                const Icon = fileIcons[file.type] || FileText
                return (
                  <TableRow key={file.id} className="border-border hover:bg-muted/30">
                    <TableCell className="font-medium">
                      <div className="flex items-center gap-2">
                        <Icon className="h-4 w-4 text-muted-foreground" />
                        <span className="text-sm truncate max-w-[200px]">
                          {file.name}
                        </span>
                      </div>
                    </TableCell>
                    <TableCell className="text-sm text-muted-foreground font-mono">
                      {file.size}
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-1">
                        {Array.from({ length: 5 }).map((_, i) => (
                          <div
                            key={i}
                            className={`h-1.5 w-1.5 rounded-full ${
                              i < file.replicas ? "bg-primary" : "bg-muted"
                            }`}
                          />
                        ))}
                        <span className="ml-1 text-xs text-muted-foreground">
                          {file.replicas}
                        </span>
                      </div>
                    </TableCell>
                    <TableCell>
                      <Badge
                        variant="outline"
                        className={
                          file.status === "healthy"
                            ? "border-success/30 bg-success/10 text-success"
                            : file.status === "syncing"
                            ? "border-warning/30 bg-warning/10 text-warning"
                            : "border-destructive/30 bg-destructive/10 text-destructive"
                        }
                      >
                        {file.status}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-sm text-muted-foreground">
                      {file.modified}
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

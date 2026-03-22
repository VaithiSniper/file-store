"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Activity, Upload, Download, UserPlus, AlertTriangle, CheckCircle } from "lucide-react"

const activities = [
  {
    id: 1,
    type: "upload",
    message: "project-backup-2024.tar.gz uploaded",
    time: "2 min ago",
    icon: Upload,
    color: "text-chart-1",
  },
  {
    id: 2,
    type: "peer",
    message: "node-delta-04 connected",
    time: "15 min ago",
    icon: UserPlus,
    color: "text-success",
  },
  {
    id: 3,
    type: "download",
    message: "database-dump.sql replicated to 3 peers",
    time: "1 hour ago",
    icon: Download,
    color: "text-chart-2",
  },
  {
    id: 4,
    type: "warning",
    message: "training-video.mp4 has low replica count",
    time: "2 hours ago",
    icon: AlertTriangle,
    color: "text-warning",
  },
  {
    id: 5,
    type: "success",
    message: "Storage optimization completed",
    time: "4 hours ago",
    icon: CheckCircle,
    color: "text-success",
  },
]

export function ActivityFeed() {
  return (
    <Card className="border-border bg-card">
      <CardHeader className="flex flex-row items-center justify-between pb-4">
        <div className="flex items-center gap-2">
          <Activity className="h-4 w-4 text-muted-foreground" />
          <CardTitle className="text-sm font-medium">Recent Activity</CardTitle>
        </div>
      </CardHeader>
      <CardContent className="pt-0">
        <div className="space-y-4">
          {activities.map((activity, index) => (
            <div key={activity.id} className="flex gap-3">
              <div className="relative flex flex-col items-center">
                <div className={`rounded-full p-1.5 bg-muted`}>
                  <activity.icon className={`h-3 w-3 ${activity.color}`} />
                </div>
                {index < activities.length - 1 && (
                  <div className="flex-1 w-px bg-border mt-2" />
                )}
              </div>
              <div className="flex-1 pb-4">
                <p className="text-sm">{activity.message}</p>
                <p className="text-xs text-muted-foreground mt-0.5">
                  {activity.time}
                </p>
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  )
}

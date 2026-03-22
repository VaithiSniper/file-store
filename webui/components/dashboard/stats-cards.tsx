"use client"

import { Card, CardContent } from "@/components/ui/card"
import { Database, HardDrive, Users, Activity, ArrowUp, ArrowDown } from "lucide-react"

const stats = [
  {
    label: "Total Storage",
    value: "2.4 TB",
    subtext: "of 4 TB",
    icon: Database,
    change: "+12%",
    trend: "up",
    color: "text-chart-1",
    bgColor: "bg-chart-1/10",
  },
  {
    label: "Files Stored",
    value: "14,328",
    subtext: "across network",
    icon: HardDrive,
    change: "+284",
    trend: "up",
    color: "text-chart-2",
    bgColor: "bg-chart-2/10",
  },
  {
    label: "Connected Peers",
    value: "7",
    subtext: "of 12 known",
    icon: Users,
    change: "+2",
    trend: "up",
    color: "text-chart-3",
    bgColor: "bg-chart-3/10",
  },
  {
    label: "Network Traffic",
    value: "847 MB",
    subtext: "last 24h",
    icon: Activity,
    change: "-5%",
    trend: "down",
    color: "text-chart-4",
    bgColor: "bg-chart-4/10",
  },
]

export function StatsCards() {
  return (
    <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
      {stats.map((stat) => (
        <Card key={stat.label} className="border-border bg-card">
          <CardContent className="p-5">
            <div className="flex items-start justify-between">
              <div className="flex flex-col gap-1">
                <span className="text-xs font-medium text-muted-foreground">
                  {stat.label}
                </span>
                <div className="flex items-baseline gap-2">
                  <span className="text-2xl font-semibold tracking-tight">
                    {stat.value}
                  </span>
                  <span className="text-xs text-muted-foreground">
                    {stat.subtext}
                  </span>
                </div>
              </div>
              <div className={`rounded-lg p-2 ${stat.bgColor}`}>
                <stat.icon className={`h-4 w-4 ${stat.color}`} />
              </div>
            </div>
            <div className="mt-3 flex items-center gap-1">
              {stat.trend === "up" ? (
                <ArrowUp className="h-3 w-3 text-success" />
              ) : (
                <ArrowDown className="h-3 w-3 text-destructive" />
              )}
              <span
                className={`text-xs font-medium ${
                  stat.trend === "up" ? "text-success" : "text-destructive"
                }`}
              >
                {stat.change}
              </span>
              <span className="text-xs text-muted-foreground">vs last week</span>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  )
}

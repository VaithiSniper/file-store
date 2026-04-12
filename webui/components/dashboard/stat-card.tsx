"use client"

import { Card, CardContent } from "@/components/ui/card"
import { Database, HardDrive, Users, Activity, ArrowUp, ArrowDown } from "lucide-react"
import { Stat } from "@/types/stats"

export function StatCard({ stat }: { stat: Stat }) {
  return (
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
            className={`text-xs font-medium ${stat.trend === "up" ? "text-success" : "text-destructive"
              }`}
          >
            {stat.change}
          </span>
          <span className="text-xs text-muted-foreground">vs last week</span>
        </div>
      </CardContent>
    </Card>
  )
}

"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { ChartContainer, ChartTooltip, ChartTooltipContent } from "@/components/ui/chart"
import { Area, AreaChart, CartesianGrid, XAxis, YAxis } from "recharts"
import { ArrowRightLeft } from "lucide-react"

const networkData = [
  { time: "00:00", incoming: 120, outgoing: 85 },
  { time: "04:00", incoming: 180, outgoing: 120 },
  { time: "08:00", incoming: 320, outgoing: 280 },
  { time: "12:00", incoming: 450, outgoing: 380 },
  { time: "16:00", incoming: 380, outgoing: 420 },
  { time: "20:00", incoming: 280, outgoing: 240 },
  { time: "Now", incoming: 340, outgoing: 290 },
]

const chartConfig = {
  incoming: {
    label: "Incoming",
    color: "var(--chart-1)",
  },
  outgoing: {
    label: "Outgoing",
    color: "var(--chart-2)",
  },
}

export function NetworkChart() {
  return (
    <Card className="border-border bg-card">
      <CardHeader className="flex flex-row items-center justify-between pb-2">
        <div className="flex items-center gap-2">
          <ArrowRightLeft className="h-4 w-4 text-muted-foreground" />
          <CardTitle className="text-sm font-medium">Data Transfer</CardTitle>
        </div>
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2">
            <span className="text-xs text-muted-foreground">Outgoing</span>
            <span className="text-sm font-semibold text-chart-1">496 GB</span>
          </div>
          <div className="flex items-center gap-2">
            <span className="text-xs text-muted-foreground">Incoming</span>
            <span className="text-sm font-semibold text-chart-2">381 GB</span>
          </div>
        </div>
      </CardHeader>
      <CardContent className="pt-0">
        <ChartContainer config={chartConfig} className="h-[200px] w-full">
          <AreaChart data={networkData} margin={{ top: 10, right: 10, left: 0, bottom: 0 }}>
            <defs>
              <linearGradient id="incomingGradient" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="var(--chart-1)" stopOpacity={0.3} />
                <stop offset="95%" stopColor="var(--chart-1)" stopOpacity={0} />
              </linearGradient>
              <linearGradient id="outgoingGradient" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="var(--chart-2)" stopOpacity={0.3} />
                <stop offset="95%" stopColor="var(--chart-2)" stopOpacity={0} />
              </linearGradient>
            </defs>
            <CartesianGrid strokeDasharray="3 3" stroke="var(--border)" vertical={false} />
            <XAxis
              dataKey="time"
              axisLine={false}
              tickLine={false}
              tick={{ fill: "var(--muted-foreground)", fontSize: 10 }}
            />
            <YAxis
              axisLine={false}
              tickLine={false}
              tick={{ fill: "var(--muted-foreground)", fontSize: 10 }}
              tickFormatter={(value) => `${value}MB`}
            />
            <ChartTooltip content={<ChartTooltipContent />} />
            <Area
              type="monotone"
              dataKey="incoming"
              stroke="var(--chart-1)"
              strokeWidth={2}
              fill="url(#incomingGradient)"
            />
            <Area
              type="monotone"
              dataKey="outgoing"
              stroke="var(--chart-2)"
              strokeWidth={2}
              fill="url(#outgoingGradient)"
            />
          </AreaChart>
        </ChartContainer>
      </CardContent>
    </Card>
  )
}

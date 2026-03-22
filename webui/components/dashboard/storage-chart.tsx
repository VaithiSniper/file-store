"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { ChartContainer, ChartTooltip, ChartTooltipContent } from "@/components/ui/chart"
import { Bar, BarChart, CartesianGrid, XAxis, YAxis } from "recharts"
import { Database } from "lucide-react"
import { Progress } from "@/components/ui/progress"

const storageData = [
  { category: "Documents", size: 420, color: "var(--chart-1)" },
  { category: "Media", size: 890, color: "var(--chart-2)" },
  { category: "Archives", size: 560, color: "var(--chart-3)" },
  { category: "Code", size: 320, color: "var(--chart-4)" },
  { category: "Other", size: 210, color: "var(--chart-5)" },
]

const chartConfig = {
  size: {
    label: "Size (GB)",
    color: "var(--chart-1)",
  },
}

export function StorageChart() {
  const totalUsed = 2400
  const totalCapacity = 4000
  const usagePercent = (totalUsed / totalCapacity) * 100

  return (
    <Card className="border-border bg-card">
      <CardHeader className="flex flex-row items-center justify-between pb-2">
        <div className="flex items-center gap-2">
          <Database className="h-4 w-4 text-muted-foreground" />
          <CardTitle className="text-sm font-medium">Storage Breakdown</CardTitle>
        </div>
      </CardHeader>
      <CardContent className="pt-0">
        <div className="mb-4">
          <div className="flex items-center justify-between mb-2">
            <span className="text-xs text-muted-foreground">Total Usage</span>
            <span className="text-xs font-medium">
              {totalUsed} GB / {totalCapacity} GB
            </span>
          </div>
          <Progress value={usagePercent} className="h-2" />
        </div>
        <ChartContainer config={chartConfig} className="h-[160px] w-full">
          <BarChart data={storageData} layout="vertical" margin={{ top: 0, right: 10, left: 0, bottom: 0 }}>
            <CartesianGrid strokeDasharray="3 3" stroke="var(--border)" horizontal={false} />
            <XAxis
              type="number"
              axisLine={false}
              tickLine={false}
              tick={{ fill: "var(--muted-foreground)", fontSize: 10 }}
              tickFormatter={(value) => `${value}GB`}
            />
            <YAxis
              type="category"
              dataKey="category"
              axisLine={false}
              tickLine={false}
              tick={{ fill: "var(--muted-foreground)", fontSize: 10 }}
              width={70}
            />
            <ChartTooltip content={<ChartTooltipContent />} />
            <Bar
              dataKey="size"
              fill="var(--chart-1)"
              radius={[0, 4, 4, 0]}
            />
          </BarChart>
        </ChartContainer>
      </CardContent>
    </Card>
  )
}

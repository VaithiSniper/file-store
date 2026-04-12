"use client"

import { Card, CardContent } from "@/components/ui/card"
import { useHyperstoreStatsStore } from "@/hooks/store"
import { Database, HardDrive, Users, Activity, ArrowUp, ArrowDown } from "lucide-react"
import { StatCard } from "./stat-card";

export function StatCardsRow() {
    const stats = useHyperstoreStatsStore((state) => state.stats);

    return (
        <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
            {
                stats.length === 0 ?
                    <div className="border p-4 col-span-full text-center text-muted-foreground">
                        No stats available. Please connect to a node to see real-time statistics.
                    </div>
                    :
                    stats.map((stat) => (
                        <StatCard key={stat.label} stat={stat} />
                    ))
            }
        </div>
    )
}

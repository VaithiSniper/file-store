"use client"

import { HardDrive, Settings, Bell, RefreshCw } from "lucide-react"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import Link from "next/link"

export function DashboardHeader() {
  return (
    <header className="sticky top-0 z-50 border-b border-border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="flex h-14 items-center justify-between px-6">
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2">
            <div className="flex h-8 w-8 items-center justify-center rounded-md bg-transparent">
              <HardDrive className="h-4 w-4 text-primary-foreground" />
            </div>
            <span className="text-lg font-semibold tracking-tight">HyperStore</span>
          </div>
          <nav className="hidden md:flex items-center gap-1">
            <Link href="/" passHref>
              <Button variant="ghost" size="sm" className="text-foreground">
                Overview
              </Button>
            </Link>
            <Link href="/files" passHref>
              <Button variant="ghost" size="sm" className="text-muted-foreground">
                Files
              </Button>
            </Link>
            <Link href="/peers" passHref>
              <Button variant="ghost" size="sm" className="text-muted-foreground">
                Peers
              </Button>
            </Link>
          </nav>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="ghost" size="icon" className="h-8 w-8">
            <RefreshCw className="h-4 w-4" />
          </Button>
          <Button variant="ghost" size="icon" className="h-8 w-8 relative">
            <Bell className="h-4 w-4" />
            <span className="absolute top-1 right-1 h-2 w-2 rounded-full bg-success" />
          </Button>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="h-8 w-8">
                <Settings className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem>Node Settings</DropdownMenuItem>
              <DropdownMenuItem>Network Config</DropdownMenuItem>
              <DropdownMenuItem>Storage Limits</DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>
    </header >
  )
}

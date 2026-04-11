'use client'

import { useState } from 'react'
import { Info, Network } from 'lucide-react'
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Peer } from '@/types/peer'
import { validateUrl } from '@/lib/utils'
import { tr } from 'date-fns/locale'

interface AddPeerModalProps {
    open: boolean
    onOpenChange: (open: boolean) => void
    onSubmit: (peerAPIURL: string) => void
}

export function AddPeerModal({ open, onOpenChange, onSubmit }: AddPeerModalProps) {
    const [peerAPIURL, setPeerAPIURL] = useState('')
    const [isAdding, setIsAdding] = useState(false)
    const [error, setError] = useState('')

    const handleAdd = async () => {
        setError('')
        let trimmedAPIURL = peerAPIURL.trim();
        if (!trimmedAPIURL.startsWith('http://') && !trimmedAPIURL.startsWith('https://')) {
            trimmedAPIURL = 'http://' + trimmedAPIURL;
        }
        if (!trimmedAPIURL || trimmedAPIURL === 'http://' || !validateUrl(trimmedAPIURL)) {
            setError('Please enter a valid peer API URL')
            return
        }


        setIsAdding(true)
        try {
            await onSubmit(trimmedAPIURL)
            setPeerAPIURL('')
            onOpenChange(false)
        } catch (err) {
            setError(err instanceof Error ? err.message : 'An unexpected error occurred while adding the peer')
            console.error('Error adding peer:', err)
        } finally {
            setIsAdding(false)
        }
    }

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="max-w-2xl">
                <DialogHeader>
                    <DialogTitle className="flex items-center gap-2">
                        <Network className="size-5" />
                        Discover remote peers
                    </DialogTitle>
                    <DialogDescription>
                        Discover peers via API URL to manage all your nodes in one place.
                    </DialogDescription>
                </DialogHeader>

                <div className="space-y-4">
                    <div>
                        <Label htmlFor="peer-address">Peer API URL</Label>
                        <Input
                            id="peer-address"
                            placeholder="example.com:8080 or 192.168.1.100:8080"
                            value={peerAPIURL}
                            onChange={(e) => setPeerAPIURL(e.target.value)}
                            className="mt-1"
                        />
                    </div>

                    {/* Info Box */}
                    <div className="rounded-lg border border-border bg-secondary p-3">
                        <div className="flex gap-2">
                            <Info className="size-4 text-muted-foreground mt-0.5 shrink-0" />
                            <div className="text-xs text-muted-foreground">
                                <p className="font-medium text-foreground">Connection Verification</p>
                                <p className="mt-1">
                                    Peer must be online and reachable for a successful connection.
                                </p>
                            </div>
                        </div>
                    </div>

                    {/* Error Message */}
                    {
                        error && (
                            <div className="rounded-lg border border-destructive/20 bg-destructive/5 p-3">
                                <p className="text-sm text-destructive">{error}</p>
                            </div>
                        )
                    }
                </div>

                <DialogFooter>
                    <Button variant="outline" onClick={() => onOpenChange(false)} disabled={isAdding}>
                        Cancel
                    </Button>
                    <Button onClick={async () => (await handleAdd())} disabled={isAdding} className="gap-2">
                        {
                            isAdding && (
                                <div className="size-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
                            )
                        }
                        {isAdding ? 'Connecting...' : 'Add Peer'}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    )
}

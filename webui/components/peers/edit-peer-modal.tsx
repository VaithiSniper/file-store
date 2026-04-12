'use client'

import { useState } from 'react'
import { AlertCircle, Globe, Zap } from 'lucide-react'
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
import { Badge } from '@/components/ui/badge'
import { Peer } from '@/types/peer'

interface EditPeerModalProps {
    open: boolean
    onOpenChange: (open: boolean) => void
    peerId: string
    peerName: string
    peerAddress: string
    peerLatency: number
    peerBandwidth: string
    onSubmit?: (peerId: string, newPeerName: string) => void
}

export interface PeerEditData {
    peerId: string
    peerName: string
    customNotes?: string
}

export function EditPeerModal({
    open,
    onOpenChange,
    peerId,
    peerName,
    peerAddress,
    peerLatency,
    peerBandwidth,
    onSubmit,
}: EditPeerModalProps) {
    const [name, setName] = useState(peerName)
    const [notes, setNotes] = useState('')
    const [isSaving, setIsSaving] = useState(false)

    const handleSave = () => {
        setIsSaving(true)
        try {
            onSubmit?.(peerId, name)
            onOpenChange(false)
        } catch (error) {
            alert('Failed to save peer information. Please try again.')
        } finally {
            setIsSaving(false)
        }
    }

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Edit Peer Information</DialogTitle>
                    <DialogDescription>
                        Update peer name, notes, and custom settings for this remote node
                    </DialogDescription>
                </DialogHeader>

                <div className="space-y-4">
                    <div className="space-y-4">
                        <label className="text-sm font-medium text-foreground">Peer ID</label>
                        <div className="flex items-center gap-2 rounded-lg bg-muted p-3">
                            <code className="flex-1 font-mono text-xs text-muted-foreground">
                                {peerId}
                            </code>
                            <Badge variant="outline" className="whitespace-nowrap">
                                Read-only
                            </Badge>
                        </div>
                    </div>
                    <div className="space-y-4">
                        <label className="text-sm font-medium text-foreground">Address</label>
                        <div className="flex items-center gap-2 rounded-lg bg-muted p-3 text-sm">
                            <Globe className="size-4 text-muted-foreground" />
                            <code className="flex-1 font-mono text-xs text-muted-foreground">
                                {peerAddress}
                            </code>
                            <Badge variant="outline" className="whitespace-nowrap">
                                Read-only
                            </Badge>
                        </div>
                    </div>


                    <div className="space-y-4">
                        <label htmlFor="peer-name" className="text-sm font-medium text-foreground">
                            Peer Name
                        </label>
                        <Input
                            id="peer-name"
                            placeholder="Enter a friendly name for this peer"
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            className="bg-muted"
                        />
                    </div>

                </div>

                <DialogFooter>
                    <Button variant="outline" onClick={() => onOpenChange(false)}>
                        Cancel
                    </Button>
                    <Button onClick={handleSave} disabled={isSaving || !name.trim()}>
                        {isSaving ? 'Saving...' : 'Save Changes'}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    )
}

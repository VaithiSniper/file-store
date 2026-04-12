'use client'

import { AlertTriangle } from 'lucide-react'
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'

interface DeleteConfirmationModalProps {
    open: boolean
    onOpenChange: (open: boolean) => void
    title: string
    description?: string
    itemId: string
    itemName?: string
    itemType: 'peer' | 'file' | 'folder'
    isDeleting?: boolean
    onSubmit?: (peerId: string) => void
    variant?: 'destructive' | 'warning'
}

export function DeleteConfirmationModal({
    open,
    onOpenChange,
    title,
    description,
    itemId,
    itemName,
    itemType,
    isDeleting = false,
    onSubmit,
    variant = 'destructive',
}: DeleteConfirmationModalProps) {
    const handleSubmit = () => {
        try {
            onSubmit?.(itemId)
            onOpenChange(false)
        } catch (error) {
            alert(`Failed to delete ${itemType}. Please try again.`)
        }
    }

    const getTypeIcon = () => {
        switch (itemType) {
            case 'peer':
                return 'peer node'
            case 'file':
                return 'file'
            case 'folder':
                return 'folder'
            default:
                return 'item'
        }
    }

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent>
                <DialogHeader>
                    <div className="flex items-center gap-3">
                        <div className="flex size-10 items-center justify-center rounded-full bg-destructive/10">
                            <AlertTriangle className="size-5 text-destructive" />
                        </div>
                        <div>
                            <DialogTitle>{title}</DialogTitle>
                        </div>
                    </div>
                    <DialogDescription className="mt-2">{description}</DialogDescription>
                </DialogHeader>

                <div className="space-y-3">
                    <div className="rounded-lg border border-destructive/20 bg-destructive/5 p-3">
                        <p className="text-sm text-foreground">
                            You&apos;re about to permanently delete this {getTypeIcon()}:
                        </p>
                        <p className="mt-2 font-mono text-sm font-medium text-destructive">{itemName || itemId}</p>
                    </div>

                    {itemType === 'peer' && (
                        <div className="rounded-lg bg-secondary p-3 text-xs text-muted-foreground">
                            <p className="font-medium text-foreground">This action will:</p>
                            <ul className="mt-2 list-inside space-y-1">
                                <li>• Disconnect from this peer node</li>
                                <li>• Remove the connection from your network</li>
                                <li>• Stop syncing files with this peer</li>
                                <li>• This action cannot be undone</li>
                            </ul>
                        </div>
                    )}

                    {itemType === 'file' && (
                        <div className="rounded-lg bg-secondary p-3 text-xs text-muted-foreground">
                            <p className="font-medium text-foreground">This action will:</p>
                            <ul className="mt-2 list-inside space-y-1">
                                <li>• Remove the file from all connected peers</li>
                                <li>• Free up storage space across the network</li>
                                <li>• This action cannot be undone</li>
                            </ul>
                        </div>
                    )}
                </div>

                <DialogFooter>
                    <Button variant="outline" onClick={() => onOpenChange(false)}>
                        Cancel
                    </Button>
                    <Button
                        variant="destructive"
                        onClick={handleSubmit}
                        disabled={isDeleting}
                        className="gap-2"
                    >
                        {isDeleting && (
                            <div className="size-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
                        )}
                        {isDeleting ? 'Deleting...' : 'Delete'}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    )
}

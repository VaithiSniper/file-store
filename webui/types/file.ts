enum FileSyncStatus {
    Created = "created",
    SyncingWithDataMessage = "syncing",
    SyncingWithStreaming = "streaming",
    Synced = "synced",
    Failed = "failed",
    DeleteSyncing = "delete_syncing",
    Deleting = "deleting",
}
type FileType = {
    basePath: string;
    keyPath: string;
    fullPath: string;
    permissions: string;
    size: string;
    createdAt: string;
    updatedAt: string;
    syncStatus: FileSyncStatus;
}

export type { FileType }
export { FileSyncStatus }
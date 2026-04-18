import type { NextRequest } from 'next/server'
import { PEER_API_HEALTH_STATUS, PeerListEntry } from "@/types/peer";
import { FileSyncStatus, FileType } from "@/types/file";

type FileReponseType = {
    base_path: string;
    key_path: string;
    full_path: string;
    file_mode: string;
    size: string;
    created_at: string;
    updated_at: string;
    sync_status: string;
}
type FilesAPIResponseType = {
    files: {
        file_map: Record<string, FileReponseType>;
        total_files: number;
    }
}

export async function getFilesFromPeer(apiUrl: string) {
    const peerListAPIEndpoint = `${apiUrl}/api/files`;
    let fileList: FileType[] = [];
    try {
        const resp = await fetch(peerListAPIEndpoint);
        const result: FilesAPIResponseType = await resp.json();
        if (resp.ok && result) {
            console.log(`Files fetched successfully from peer at ${apiUrl}:`, Object.entries(result.files.file_map));
            Object.entries(result.files.file_map).forEach(([_, fileMeta]) => {
                console.log("Processing file metadata:", fileMeta);
                fileList.push({
                    basePath: fileMeta.base_path,
                    keyPath: fileMeta.key_path,
                    fullPath: fileMeta.full_path,
                    permissions: fileMeta.file_mode,
                    size: fileMeta.size,
                    createdAt: fileMeta.created_at,
                    updatedAt: fileMeta.updated_at,
                    syncStatus: fileMeta.sync_status as FileSyncStatus,
                });
            });
            if (fileList.length != result.files.total_files) {
                console.error("File count mismatch");
            }
        }
    } catch (e) {
        console.log("No files found for", apiUrl);
    }

    return fileList;
}

export async function GET(
    request: NextRequest,
) {
    const { searchParams } = request.nextUrl;
    const peer_api_url = searchParams.get("peer_api_url") || "";
    const fileList = await getFilesFromPeer(peer_api_url);

    return Response.json({ files: fileList });
}

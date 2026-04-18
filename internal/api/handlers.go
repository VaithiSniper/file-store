package api

import (
	"fmt"

	"file-store/internal/storage"
	"file-store/internal/util"
)

// /api/status
func getStatus() StatusResponse {
	return StatusResponse{
		API:   "healthy",
		Store: "healthy",
		DB:    "healthy",
	}
}

// /api/config
func getSelfConfig() SelfConfigResponse {
	serverIdentity := util.GetIdentity()

	storeOpts := storage.GetStoreInstance().StoreOpts
	messageFormat := "nil"
	if storeOpts.MessageFormat != nil {
		messageFormat = fmt.Sprintf("%T", storeOpts.MessageFormat)
	}

	apiServerOpts := GetAPIServerInstance().ApiServerOpts
	apiUrl := fmt.Sprintf("http://%s/api", apiServerOpts.APIServerListenAddress)

	status := "healthy"
	uptime := util.GetUptime().String()

	return SelfConfigResponse{
		Id:                  serverIdentity.Id,
		Name:                serverIdentity.Name,
		ListenAddress:       storeOpts.ListenAddress,
		ApiUrl:              apiUrl,
		Status:              status,
		Uptime:              uptime,
		Latency:             "10ms",
		DataIn:              "1.5GB",
		DataOut:             "3.2GB",
		PathTransformFunc:   storage.PathTransformFuncName(storeOpts.PathTransformFunc),
		MessageFormat:       messageFormat,
		BaseStorageLocation: storeOpts.BaseStorageLocation,
		BootstrapNodes:      storeOpts.BootstrapNodes,
	}
}

// /api/peers
func getPeerList() []string {
	var peerList []string
	storageInstance := storage.GetStoreInstance()
	for peer := range storageInstance.PeerMap {
		peerList = append(peerList, peer)
	}
	return peerList
}

// /api/files
func getFiles() FileInfoResponse {
	storageInstance := storage.GetStoreInstance()
	var fileMap FileTypeMap = make(FileTypeMap, 0)
	var totalFiles int = 0

	for fileKey, file := range storageInstance.FileMap {
		fileMap[fileKey] = FileType{
			BasePath:    file.BasePath,
			KeyPath:     file.KeyPath,
			FullPath:    file.BasePath + "/" + file.KeyPath,
			Permissions: file.FileMode.Perm().String(),
			Size:        util.ConvertFileSizeToReadableFormat(file.FileSize),
			CreatedAt:   file.CreatedAt,
			UpdatedAt:   file.UpdatedAt,
			SyncStatus:  file.SyncStatus,
		}
		totalFiles++
	}
	return FileInfoResponse{
		FileMap:    fileMap,
		TotalFiles: totalFiles,
	}
}

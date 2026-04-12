package api

import (
	"fmt"

	"file-store/internal/file"
	"file-store/internal/storage"
	"file-store/internal/util"
)

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

// /api/peer/:peerId/files
func getFiles() map[string]file.File {
	storageInstance := storage.GetStoreInstance()
	return storageInstance.FileMap
}

package api

import (
	"fmt"

	"file-store/internal/file"
	"file-store/internal/storage"
)

func getSelfConfig() SelfConfigResponse {
	opts := storage.GetStoreInstance().StoreOpts
	messageFormat := "nil"
	if opts.MessageFormat != nil {
		messageFormat = fmt.Sprintf("%T", opts.MessageFormat)
	}
	return SelfConfigResponse{
		Id:                  "1asda12313ased-1231-1231-1231-asdas123123",
		Name:                opts.Name,
		ListenAddress:       opts.ListenAddress,
		ApiUrl:              "http://localhost:8080/api",
		Status:              "healthy",
		Uptime:              "72h3m4s",
		Latency:             "10ms",
		DataIn:              "1.5GB",
		DataOut:             "3.2GB",
		PathTransformFunc:   storage.PathTransformFuncName(opts.PathTransformFunc),
		MessageFormat:       messageFormat,
		BaseStorageLocation: opts.BaseStorageLocation,
		BootstrapNodes:      opts.BootstrapNodes,
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

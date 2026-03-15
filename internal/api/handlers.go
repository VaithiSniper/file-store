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
		ListenAddress:       opts.ListenAddress,
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

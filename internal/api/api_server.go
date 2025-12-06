package api

import (
	"file-store/internal/storage"
	"log"

	"github.com/gin-gonic/gin"
)

func getPeerList() []string {
	var peerList []string
	storageInstance := storage.GlobalStore
	if storageInstance == nil {
		return peerList
	}

	for peer, _ := range storageInstance.PeerMap {
		peerList = append(peerList, peer)
	}
	return peerList
}

func RegisterAPIRoutes(router *gin.Engine) {
	apiGroup := router.Group("/api")
	{
		apiGroup.GET(
			"/status", func(c *gin.Context) {
				c.JSON(200, gin.H{"status": "ok"})
			},
		)
		apiGroup.GET(
			"/peers", func(c *gin.Context) {
				c.JSON(200, gin.H{"peers": getPeerList()})
			},
		)
		// Additional API routes can be registered here
	}
}

func StartAPIServer(apiServerListenAddress string) {
	router := gin.Default()
	RegisterAPIRoutes(router)
	err := router.Run(apiServerListenAddress)
	if err != nil {
		log.Fatalf("Error starting API server: %v.", err)
	}
}

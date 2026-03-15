package api

import (
	"log"

	"github.com/gin-gonic/gin"
)

func StartAPIServer(apiServerListenAddress string) {
	router := gin.Default()
	RegisterAPIRoutes(router)
	err := router.Run(apiServerListenAddress)
	if err != nil {
		log.Fatalf("Error starting API server: %v.", err)
	}
}

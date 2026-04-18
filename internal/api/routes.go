package api

import "github.com/gin-gonic/gin"

func RegisterAPIRoutes(router *gin.Engine) {
	apiGroup := router.Group("/api")
	{
		apiGroup.GET(
			"/status", func(c *gin.Context) {
				c.JSON(200, getStatus())
			},
		)
		apiGroup.GET(
			"/config", func(c *gin.Context) {
				c.JSON(200, getSelfConfig())
			},
		)
		apiGroup.GET(
			"/peers", func(c *gin.Context) {
				c.JSON(200, gin.H{"peers": getPeerList()})
			},
		)
		apiGroup.GET(
			"/files", func(c *gin.Context) {
				c.JSON(200, gin.H{"files": getFiles()})
			},
		)
		// Additional API routes can be registered here
	}
}

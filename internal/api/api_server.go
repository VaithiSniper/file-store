package api

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"file-store/internal/logger"
	"file-store/internal/util"
)

const moduleName = "API"

type APIServerOpts struct {
	APIServerListenAddress string
	LogLevel               logger.LogLevel
}

type APIServer struct {
	ApiServerOpts APIServerOpts
	router        *gin.Engine
}

var GlobalAPIServer *APIServer

// GetAPIServerInstance returns a singleton instance of APIServer.
func GetAPIServerInstance() *APIServer {
	return GlobalAPIServer
}

// GinLoggerMiddleware creates a Gin middleware that uses our custom logger
func GinLoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Map Gin log levels to our custom logger levels
		var level logger.LogLevel
		if param.StatusCode >= 500 {
			level = logger.LevelError
		} else if param.StatusCode >= 400 {
			level = logger.LevelWarning
		} else {
			level = logger.LevelInfo
		}

		// Format the log using our fmt
		message := fmt.Sprintf("%s %s %d %s %s",
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
		)

		// Use logger as per level
		switch level {
		case logger.LevelError:
			logger.LogError(moduleName, message)
		case logger.LevelWarning:
			logger.LogWarning(moduleName, message)
		default:
			logger.LogInfo(moduleName, message)
		}

		return ""
	})
}

// GinRecoveryMiddleware creates a Gin recovery middleware that uses our custom logger
func GinRecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			logger.LogError(moduleName, "Panic recovered: %s", err)
		} else {
			logger.LogError(moduleName, "Panic recovered: %v", recovered)
		}
		c.AbortWithStatus(500)
	})
}

func CreateAPIServerWithUserOptions(opts APIServerOpts) *APIServer {
	router := gin.New()
	router.Use(GinLoggerMiddleware())
	router.Use(GinRecoveryMiddleware())
	RegisterAPIRoutes(router)
	return &APIServer{
		ApiServerOpts: opts,
		router:        router,
	}
}

func (apiServer *APIServer) SetupAPIServer() {
	GlobalAPIServer = apiServer

	err := apiServer.router.Run(apiServer.ApiServerOpts.APIServerListenAddress)
	if err != nil {
		logger.LogCritical(moduleName, "Error starting API server: %v.", err)
		os.Exit(1)
	}
}

// StartAPIServer initializes and starts the API server with command line arguments
func StartAPIServer(commandLineArgs util.CommandLineArgs) {
	opts := APIServerOpts{
		APIServerListenAddress: commandLineArgs.ApiServerListenAddress,
		LogLevel:               commandLineArgs.LogLevel,
	}

	apiServer := CreateAPIServerWithUserOptions(opts)
	apiServer.SetupAPIServer()
}

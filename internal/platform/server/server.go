package server

import (
	"github.com/gin-gonic/gin"
	"storePrices/internal/platform/conf"
)

// New creates the Gin Engine
func New() *gin.Engine {
	r := gin.Default()
	// Add global middlewares here (CORS, Auth, etc)
	return r
}

// Start launches the server (Invoked by Fx)
func Start(cfg *conf.Config, r *gin.Engine) {
	// In a real app, use Fx Lifecycle to handle graceful shutdown
	go func() {
		if err := r.Run(":" + cfg.Server.Port); err != nil {
			panic(err)
		}
	}()
}

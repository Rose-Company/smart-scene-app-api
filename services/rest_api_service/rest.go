package rest_api_service

import (
	"os"
	"smart-scene-app-api/common"
	"smart-scene-app-api/internal/handlers"
	"smart-scene-app-api/middleware"
	"smart-scene-app-api/server"

	"github.com/gin-contrib/requestid"

	"github.com/gin-contrib/cors"
	_ "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RestHandler(sc server.ServerContext) func() *gin.Engine {
	return func() *gin.Engine {
		mode, ok := os.LookupEnv(common.ENV_GIN_DEBUG)
		if !ok {
			mode = "debug"
		}
		router := gin.New()
		gin.SetMode(mode)
		router.Use(requestid.New())
		router.Use(middleware.Logger(sc), middleware.Recovery(sc))
		router.Use(cors.New(cors.Config{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{"*"},
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		}))

		// Swagger documentation
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		// sc.InitAuthorizationData()

		health := router.Group("/health")
		{
			health.GET("/status", handlers.Check(sc))
		}

		// Handler
		handler := handlers.NewHandler(sc)
		handler.RegisterRouter(router)

		return router
	}
}

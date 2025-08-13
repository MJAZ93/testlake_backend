package app

import (
	"os"

	"testlake/middleware"
	"testlake/utils"

	"github.com/gin-gonic/gin"
)

// @title TestLake API
// @version 1.0
// @description Test Data Management Platform API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func ServeApplication() {
	router := gin.Default()

	router.Use(middleware.DefaultAuthMiddleware())

	Swagger(router)

	baseRoute := router.Group("/api/v1")

	publicRoutes := baseRoute.Group("")
	PublicRoutes(publicRoutes)

	privateRoutes := baseRoute.Group("")
	privateRoutes.Use(middleware.JWTAuthMiddleware())
	PrivateRoutes(privateRoutes)

	router.NoRoute(utils.HandleNoRoute())

	ip := os.Getenv("IP")
	port := os.Getenv("PORT")
	router.Run(ip + ":" + port)
}

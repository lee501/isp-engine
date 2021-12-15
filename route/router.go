package route

import (
	"github.com/gin-gonic/gin"
	v1 "isp-engine/api/v1"
	"isp-engine/middleware"
)

func InitRouter() *gin.Engine {
	route := gin.New()
	route.Use(gin.Logger(), gin.Recovery())
	api := route.Group("/api")
	//api.Use(middleware.Auth(), middleware.Cors())
	api.Use(middleware.Cors())
	api.Use(middleware.JSONAppErrorReporter())
	api.POST("/change_pass", v1.ChangePassword)
	api.POST("/open_port", v1.OpenPort)
	api.GET("/ws", v1.OpenSsh)

	globalRoutes := middleware.GetglobalRoutes()
	for _, r := range globalRoutes {
		route.RouterGroup.Handle(r.HttpMethod, r.RelativePath, r.Handlers...)
	}
	return route
}

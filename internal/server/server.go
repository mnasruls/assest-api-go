package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type RestApi struct {
	router *gin.Engine
}

func (r *RestApi) Serve(address string) error {
	return r.router.Run(address)
}

func NewRestApi(db *gorm.DB) *RestApi {

	router := gin.Default()
	api := &RestApi{
		router,
	}

	// health check
	router.GET("/", HealthCheck)

	// Swagger
	url := ginSwagger.URL("/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	Build(router, db)
	return api
}

func HealthCheck(c *gin.Context) {
	res := map[string]interface{}{
		"data": "Server is up and running",
	}

	c.JSON(http.StatusOK, res)
}

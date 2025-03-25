package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

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

	Build(router, db)
	return api
}

func HealthCheck(c *gin.Context) {
	res := map[string]interface{}{
		"message": "Server is up and running",
	}

	c.JSON(http.StatusOK, res)
}

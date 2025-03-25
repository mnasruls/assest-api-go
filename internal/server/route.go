package server

import (
	"assets-api-go/docs"
	"assets-api-go/internal/handlers"
	"assets-api-go/internal/repositories"
	"assets-api-go/internal/services"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Build(route *gin.Engine, db *gorm.DB) {
	assetRepo := repositories.NewAssetRepository(db)
	assetServie := services.NewAssetService(assetRepo)
	assetHandler := handlers.NewAssetHandler(assetServie)

	path := "api/v1"
	// Swagger
	url := ginSwagger.URL("/swagger/doc.json")
	docs.SwaggerInfo.BasePath = path
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	route.POST(path+"/assets", assetHandler.CreateAsset)
	route.GET(path+"/assets", assetHandler.GetAssets)
	route.GET(path+"/assets/:id", assetHandler.GetAssetById)
	route.PUT(path+"/assets/:id", assetHandler.UpdateAsset)
	route.DELETE(path+"/assets/:id", assetHandler.DeleteAsset)
}

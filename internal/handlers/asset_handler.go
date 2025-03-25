package handlers

import (
	"assets-api-go/internal/common"
	"assets-api-go/internal/dto"
	"assets-api-go/internal/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AssetHandlerInterface interface {
	CreateAsset(c *gin.Context)
	UpdateAsset(c *gin.Context)
	GetAssetById(c *gin.Context)
	GetAssets(c *gin.Context)
	DeleteAsset(c *gin.Context)
}

type assetHandler struct {
	service services.AssetServiceInterface
}

func NewAssetHandler(service services.AssetServiceInterface) AssetHandlerInterface {
	return &assetHandler{service: service}
}

func (h *assetHandler) CreateAsset(c *gin.Context) {
	request := new(dto.AssetInputDto)
	err := c.Bind(&request)
	if err != nil {
		log.Println("[assetHandler][CreateAsset] error binding request :", err)
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Error:            common.BadRequest,
			ErrorDescription: "invalid request",
		})
	}

	c.JSON(h.service.CreateAsset(request))
}

func (h *assetHandler) UpdateAsset(c *gin.Context) {
	request := new(dto.AssetInputDto)
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Error:            common.BadRequest,
			ErrorDescription: "invalid request",
		})
	}
	err := c.Bind(&request)
	if err != nil {
		log.Println("[assetHandler][UpdateAsset] error binding request :", err)
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Error:            common.BadRequest,
			ErrorDescription: "invalid request",
		})
	}

	c.JSON(h.service.UpdateAsset(id, request))

}

func (h *assetHandler) GetAssetById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Error:            common.BadRequest,
			ErrorDescription: "invalid request",
		})
	}
	c.JSON(h.service.GetAssetById(id))
}

func (h *assetHandler) GetAssets(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		log.Println("[assetHandler][GetAssets] error binding request :", err)
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Error:            common.BadRequest,
			ErrorDescription: "invalid request",
		})
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		log.Println("[assetHandler][GetAssets] error binding request :", err)
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Error:            common.BadRequest,
			ErrorDescription: "invalid request",
		})
	}
	pagination := &dto.MetaPagination{
		Page:   page,
		Limit:  limit,
		Order:  c.Query("order"),
		SortBy: c.Query("sort_by"),
	}

	pagination = pagination.ParsePagination()
	c.JSON(h.service.GetAssets(pagination))
}

func (h *assetHandler) DeleteAsset(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Error:            common.BadRequest,
			ErrorDescription: "invalid request",
		})
	}
	c.JSON(h.service.DeleteAsset(id))
}

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

// CreateAsset creates a new asset
//
//	@Summary      Create a new asset
//	@Description  Takes an asset JSON and store in DB. Return saved JSON.
//	@Tags         assets
//	@Accept       json
//	@Produce      json
//	@Param        asset  body      dto.AssetInputDto  true  "Asset JSON"
//	@Success      200    {object}  dto.BaseResponse{data=dto.AssetOutputDto}
//	@Failure      400    {object}  dto.BaseResponse{data=nil}
//	@Failure      404    {object}  dto.BaseResponse{data=nil}
//	@Failure      500    {object}  dto.BaseResponse{data=nil}
//	@Router       /assets [post]
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

// UpdateAsset updates an asset
//
//	@Summary      Update an asset
//	@Description  Takes an asset JSON and update in DB. Return updated JSON.
//	@Tags         assets
//	@Accept       json
//	@Produce      json
//	@Param        id   path      string  true  "Asset ID"
//	@Param        asset  body      dto.AssetInputDto  true  "Asset JSON"
//	@Success      200    {object}  dto.BaseResponse{data=dto.AssetOutputDto}
//	@Failure      400    {object}  dto.BaseResponse{data=nil}
//	@Failure      404    {object}  dto.BaseResponse{data=nil}
//	@Failure      500    {object}  dto.BaseResponse{data=nil}
//	@Router       /assets/{id} [put]
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

// GetAssetById returns an asset
//
//	@Summary      Get an asset
//	@Description  Returns an asset JSON.
//	@Tags         assets
//	@Accept       json
//	@Produce      json
//	@Param        id   path      string  true  "Asset ID"
//	@Success      200    {object}  dto.BaseResponse{data=dto.AssetOutputDto}
//	@Failure      400    {object}  dto.BaseResponse{data=nil}
//	@Failure      404    {object}  dto.BaseResponse{data=nil}
//	@Failure      500    {object}  dto.BaseResponse{data=nil}
//	@Router       /assets/{id} [get]
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

// GetAssets returns a list of assets
//
//	@Summary      List assets
//	@Description  Returns a list of assets JSON.
//	@Tags         assets
//	@Accept       json
//	@Produce      json
//	@Param        page   query      int  false  "Page number"
//	@Param        limit   query      int  false  "Limit number"
//	@Param        order   query      string  false  "Order"
//	@Param        sort_by   query      string  false  "Sort by"
//	@Success      200    {object}  dto.MetaPagination{data=[]dto.AssetOutputDto}
//	@Failure      400    {object}  dto.MetaPagination{data=nil}
//	@Failure      500    {object}  dto.MetaPagination{data=nil}
//	@Router       /assets [get]
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

// DeleteAsset deletes an asset
//
//	@Summary      Delete an asset
//	@Description  Delete an asset.
//	@Tags         assets
//	@Accept       json
//	@Produce      json
//	@Param        id   path      string  true  "Asset ID"
//	@Success      200    {object}  dto.BaseResponse{data=nil,}
//	@Failure      400    {object}  dto.BaseResponse{data=nil}
//	@Failure      404    {object}  dto.BaseResponse{data=nil}
//	@Failure      500    {object}  dto.BaseResponse{data=nil}
//	@Router       /assets/{id} [delete]
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

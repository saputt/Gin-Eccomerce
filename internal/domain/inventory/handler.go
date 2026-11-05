package inventory

import (
	"net/http"
	"orders/internal/common/response"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	service *InventoryService
}

func NewInventoryHandler(service *InventoryService) *InventoryHandler {
	return &InventoryHandler{
		service: service,
	}
}

func (h *InventoryHandler) RegisterRoutes(rg *gin.RouterGroup) {
	inventory := rg.Group("/inventory")
	{
		inventory.PATCH("/:productId/stock", h.UpdateStock)
		inventory.GET("/:productId", h.FindByProduct)
	}
}

func (h *InventoryHandler) UpdateStock(c *gin.Context) {
	var req UpdateStockRequest

	productId := c.Param("productId")

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		response.ErrorRes(c, http.StatusBadRequest, err.Error())
		return
	}

	inv, err := h.service.UpdateStock(c, productId, req.Quantity)
	if err != nil {
		switch err {
		case ErrInventoryNotFound:
			response.ErrorRes(c, http.StatusNotFound, err.Error())
		default:
			response.ErrorRes(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.SuccessRes(c, http.StatusOK, "Update stock success", inv)
}

func (h *InventoryHandler) FindByProduct(c *gin.Context) {
	productId := c.Param("productId")

	inv, err := h.service.FindByProduct(c, productId)
	if err != nil {
		switch err {
		case ErrInventoryNotFound:
			response.ErrorRes(c, http.StatusNotFound, err.Error())
		default:
			response.ErrorRes(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.SuccessRes(c, http.StatusOK, "Update stock success", inv)
}

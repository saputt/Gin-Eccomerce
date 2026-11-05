package order

import (
	"net/http"
	"orders/internal/common/middleware"
	"orders/internal/common/response"
	"orders/internal/domain/auth"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service *OrderService
}

func NewOrderHandler(service *OrderService, authService *auth.AuthService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (h *OrderHandler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	orders := rg.Group("/orders")
	orders.Use(authMiddleware)
	{
		orders.POST("", h.CreateOrder)
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest

	userId, _ := c.Get(middleware.UserIdKey)

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		response.ErrorRes(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.CreateOrder(c, req, userId.(string))
	if err != nil {
		switch err {
		case ErrProductNotFound:
			response.ErrorRes(c, http.StatusNotFound, err.Error())
		default:
			response.ErrorRes(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.SuccessRes(c, http.StatusOK, "Create order success", res)
}

// func (h *OrderHandler) CancelOrder(c *gin.Context) {
// 	h.service.CancelOrder()
// }

// func (h *OrderHandler) FindOrder(c *gin.Context) {}

// func (h *OrderHandler) FindMySuborderForBuyer(c *gin.Context) {}

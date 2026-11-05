package product

import (
	"net/http"
	"orders/internal/common/middleware"
	"orders/internal/common/response"
	"orders/internal/domain/auth"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service     *ProductService
	authService *auth.AuthService
}

func NewProductHandler(service *ProductService, authService *auth.AuthService) *ProductHandler {
	return &ProductHandler{
		service:     service,
		authService: authService,
	}
}

func (h *ProductHandler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	products := rg.Group("/products")
	{
		products.GET("", h.GetProducts)
		products.GET(":productId", h.GetProduct)
	}

	products.Use(authMiddleware, middleware.RequireRoles("SELLER"))
	{
		products.POST("", h.CreateProduct)
		products.PATCH(":productId", h.UpdateProducts)
		products.DELETE(":productId", h.DeleteProduct)
	}
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	products, err := h.service.GetProducts(c)
	if err != nil {
		response.ErrorRes(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessRes(c, http.StatusOK, "Get products success", products)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	productId := c.Param("userId")

	product, err := h.service.GetProduct(c, productId)
	if err != nil {
		response.ErrorRes(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessRes(c, http.StatusOK, "Get product success", product)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	userId, _ := c.Get(middleware.UserIdKey)

	var req CreateProductRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		response.ErrorRes(c, http.StatusBadRequest, err.Error())
		return
	}

	products, err := h.service.CreateProduct(c, req, userId.(string))
	if err != nil {
		response.ErrorRes(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessRes(c, http.StatusOK, "Create product success", products)
}

func (h *ProductHandler) UpdateProducts(c *gin.Context) {
	userId, _ := c.Get(middleware.UserIdKey)
	productId := c.Param("productId")

	var req UpdateProductRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		response.ErrorRes(c, http.StatusBadRequest, err.Error())
		return
	}

	products, err := h.service.UpdateProduct(c, &req, productId, userId.(string))
	if err != nil {
		switch err {
		case ErrUnauthorized:
			response.ErrorRes(c, http.StatusUnauthorized, err.Error())
		default:
			response.ErrorRes(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.SuccessRes(c, http.StatusOK, "Get products success", products)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	userId, _ := c.Get(middleware.UserIdKey)
	productId := c.Param("productId")

	err := h.service.DeleteProduct(c, productId, userId.(string))
	if err != nil {
		switch err {
		case ErrProductNotFound:
			response.ErrorRes(c, http.StatusNotFound, err.Error())
		case ErrUnauthorized:
			response.ErrorRes(c, http.StatusUnauthorized, err.Error())
		default:
			response.ErrorRes(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.SuccessRes(c, http.StatusOK, "Delete product success", nil)
}

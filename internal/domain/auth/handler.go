package auth

import (
	"fmt"
	"net/http"
	"orders/internal/common/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *AuthService
}

func NewAuthHandler(service *AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/register", h.Register)
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	fmt.Println("masuk")

	var req LoginRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		response.ErrorRes(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.Login(c, &req)
	if err != nil {
		switch err {
		case ErrInvalidCredential:
			response.ErrorRes(c, http.StatusUnauthorized, err.Error())
			return
		default:
			response.ErrorRes(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	response.SuccessRes(c, http.StatusOK, "Login successfully", result)
}

func (h *AuthHandler) Register(c *gin.Context) {
	fmt.Println("masuk")

	var req RegisterRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		response.ErrorRes(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.Register(c, &req)
	if err != nil {
		switch err {
		case ErrInvalidCredential:
			response.ErrorRes(c, http.StatusUnauthorized, err.Error())
			return
		default:
			response.ErrorRes(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	response.SuccessRes(c, http.StatusOK, "Login successfully", result)
}

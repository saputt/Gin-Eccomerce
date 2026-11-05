package middleware

import (
	"fmt"
	"net/http"
	"orders/internal/common/response"
	"orders/internal/domain/auth"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	UserIdKey   = "userId"
	UserRoleKey = "role"
)

type TokenValidator interface {
	ValidateToken(token string) (*auth.JWTClaims, error)
}

func AuthMiddleware(tokenValidator TokenValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ErrorRes(c, http.StatusUnauthorized, "Authorization header is required")
			c.Abort()
			return
		}

		part := strings.Split(authHeader, " ")
		if len(part) != 2 || part[0] != "Bearer" {
			response.ErrorRes(c, http.StatusUnauthorized, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := part[1]
		claims, err := tokenValidator.ValidateToken(tokenString)
		if err != nil {
			response.ErrorRes(c, http.StatusUnauthorized, "Invalid token or expired")
			c.Abort()
			return
		}

		c.Set(UserIdKey, claims.Id)
		c.Set(UserRoleKey, claims.Role)

		c.Next()
	}
}

func RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get(UserRoleKey)
		if !exists {
			response.ErrorRes(c, http.StatusUnauthorized, "Role note found in context")
			return
		}

		userRole := strings.TrimSpace(fmt.Sprintf("%v", val))

		fmt.Println("UserRole in context: ", userRole)

		allowed := false
		for _, role := range allowedRoles {
			fmt.Println("UserRole in allowed: ", role)
			if strings.TrimSpace(role) == userRole {
				fmt.Println("UserRole in req: ", role)
				allowed = true
				break
			}
		}

		fmt.Println("Status allowed: ", allowed)

		if !allowed {
			response.ErrorRes(c, http.StatusUnauthorized, "You do not have permission to access this resource")
			return
		}

		c.Next()
	}
}

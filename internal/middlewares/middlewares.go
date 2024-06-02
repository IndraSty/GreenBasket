package middlewares

import (
	"net/http"
	"strings"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	tokenSvc domain.TokenService
}

func NewMiddleware(tokenSvc domain.TokenService) *Middleware {
	return &Middleware{
		tokenSvc: tokenSvc,
	}
}

func (m *Middleware) UserAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		clientToken := strings.Split(authHeader, " ")[1]

		if clientToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorization"})
			c.Abort()
			return
		}

		claims, err := m.tokenSvc.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("uid", claims.Uid)

		c.Next()
	}
}

func (m *Middleware) SellerAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		headerParts := strings.Split(authHeader, " ")

		if len(headerParts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
			c.Abort()
			return
		}

		clientToken := headerParts[1]
		if clientToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorization"})
			c.Abort()
			return
		}

		claims, err := m.tokenSvc.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("uid", claims.Uid)

		c.Next()
	}
}

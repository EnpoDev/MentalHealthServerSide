package middleware

import (
	"mental-health-companion/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func sendError(c *gin.Context, err *utils.APIError) {
	c.JSON(err.HTTPStatus, err)
	c.Abort()
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			sendError(c, utils.NewMissingTokenError())
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			sendError(c, utils.NewInvalidTokenFormatError())
			return
		}

		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			sendError(c, utils.NewInvalidTokenError())
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}

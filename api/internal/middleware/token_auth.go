package middleware

import (
	"net/http"
	"strings"

	"github.com/ElshadHu/vulnly/api/internal/repository"
	"github.com/gin-gonic/gin"
)

// TokenAuth validates API tokens (alternative to JWT)

func TokenAuth(repo *repository.DynamoDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next() // Let JWT middleware handle it
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}
		token := parts[1]
		// Check if it is an API token (starts with vly_)
		if !strings.HasPrefix(token, "vly_") {
			c.Next() // let JWT handle
			return
		}
		apiToken, err := repo.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid API token",
			})
			return
		}
		c.Set("user_id", apiToken.UserID)
		c.Set("token_id", "api_token")
		c.Next()
	}
}

package handler

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/ElshadHu/vulnly/api/internal/repository"
	"github.com/gin-gonic/gin"
)

type CreateTokenRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateTokenResponse struct {
	TokenID   string `json:"token_id"`
	Token     string `json:"token"` // Plaintext, shown only once
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type TokenResponse struct {
	TokenID    string `json:"token_id"`
	Name       string `json:"name"`
	CreatedAt  string `json:"created_at"`
	LastUsedAt string `json:"last_used_at,omitempty"`
}

func (h *API) CreateToken(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	if len(req.Name) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name must be 100 characters or less"})
		return
	}

	// Check token limit
	const maxTokensPerUser = 10
	existingTokens, err := h.repo.ListTokensByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check token limit"})
		return
	}
	if len(existingTokens) >= maxTokensPerUser {
		c.JSON(http.StatusBadRequest, gin.H{"error": "maximum 10 tokens allowed per user"})
		return
	}

	// Generate token
	plaintext, bcryptHash, sha256Lookup, err := repository.GenerateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// Store token
	token, err := h.repo.CreateToken(c.Request.Context(), userID, req.Name, bcryptHash, sha256Lookup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}

	// Return plaintext only once
	c.JSON(http.StatusCreated, CreateTokenResponse{
		TokenID:   token.TokenID,
		Token:     plaintext,
		Name:      token.Name,
		CreatedAt: token.CreatedAt.Format(time.RFC3339),
	})
}

// ListTokens returns all tokens for the authenticated user
func (h *API) ListTokens(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	tokens, err := h.repo.ListTokensByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list tokens"})
		return
	}

	response := make([]TokenResponse, len(tokens))
	for i, t := range tokens {
		response[i] = TokenResponse{
			TokenID:   t.TokenID,
			Name:      t.Name,
			CreatedAt: t.CreatedAt.Format(time.RFC3339),
		}
		if !t.LastUsedAt.IsZero() {
			response[i].LastUsedAt = t.LastUsedAt.Format(time.RFC3339)
		}
	}

	c.JSON(http.StatusOK, gin.H{"tokens": response})
}

// DeleteToken deletes a token by ID
func (h *API) DeleteToken(c *gin.Context) {
	userID := c.GetString("user_id")
	tokenID := c.Param("token_id")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.repo.DeleteToken(c.Request.Context(), userID, tokenID); err != nil {
		if errors.Is(err, repository.ErrTokenNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "token not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete token"})
		return
	}

	c.Status(http.StatusNoContent)
}

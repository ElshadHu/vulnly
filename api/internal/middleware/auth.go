package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// CognitoClaims represents the claims in a Cognito JWT
type CognitoClaims struct {
	jwt.RegisteredClaims
	TokenUse string `json:"token_use"` // "access" or "id"
	ClientID string `json:"client_id"` // App client ID
	Username string `json:"username"`  // Cognito username
	Sub      string `json:"sub"`       // User's unique ID
}

var (
	ErrMissingConfig = errors.New("missing Cognito configuration")
	ErrMissingToken  = errors.New("missing authorization header")
	ErrInvalidFormat = errors.New("invalid authorization header format")
	ErrInvalidToken  = errors.New("invalid token")
	ErrTokenClaims   = errors.New("invalid token claims")
	ErrInvalidType   = errors.New("invalid token type")
)

// Auth holds the JWKS keyfunc for token validation
type Auth struct {
	jwks     keyfunc.Keyfunc
	issuer   string
	clientID string
}

// NewAuth creates a new Auth middleware with JWKS from Cognito
func NewAuth(ctx context.Context) (*Auth, error) {
	region := os.Getenv("AWS_REGION")
	userPoolID := os.Getenv("COGNITO_USER_POOL_ID")
	clientID := os.Getenv("COGNITO_CLIENT_ID")

	if region == "" || userPoolID == "" || clientID == "" {
		return nil, ErrMissingConfig
	}

	jwksURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", region, userPoolID)

	k, err := keyfunc.NewDefaultCtx(ctx, []string{jwksURL})
	if err != nil {
		return nil, fmt.Errorf("failed to create JWKS keyfunc: %w", err)
	}

	return &Auth{
		jwks:     k,
		issuer:   fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", region, userPoolID),
		clientID: clientID,
	}, nil
}

// Middleware returns a Gin middleware that validates JWT tokens
func (a *Auth) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip if already authenticated by token middleware
		if c.GetString("user_id") != "" {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": ErrMissingToken.Error(),
			})
			return
		}

		// Bearer token format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": ErrInvalidFormat.Error(),
			})
			return
		}

		tokenString := parts[1]

		// Parse and validate the token using JWKS
		token, err := jwt.ParseWithClaims(
			tokenString,
			&CognitoClaims{},
			a.jwks.Keyfunc,
			jwt.WithIssuer(a.issuer),
			jwt.WithValidMethods([]string{"RS256"}), // Cognito uses RS256
		)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": ErrInvalidToken.Error(),
			})
			return
		}

		claims, ok := token.Claims.(*CognitoClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": ErrTokenClaims.Error(),
			})
			return
		}

		// Verify it's an access token (not id token)
		if claims.TokenUse != "access" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": ErrInvalidType.Error(),
			})
			return
		}

		// Store user ID in context for handlers to use
		c.Set("user_id", claims.Sub)
		c.Set("username", claims.Username)
		c.Next()
	}
}

package repository

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrTokenNotFound = errors.New("token not found")
	ErrMarshalToken  = errors.New("failed to marshal token")
	ErrQueryTokens   = errors.New("failed to query tokens")
	ErrInvalidToken  = errors.New("invalid token")
	ErrDeleteToken   = errors.New("failed to delete token")
)

const (
	// TokenPrefix identifies Vulnly API tokens in logs and code scanning
	TokenPrefix = "vly_"
	TokenLength = 32
	// bcryptCost balances security vs performance
	bcryptCost = 10
)

// GenerateToken creates a new token with prefix
func GenerateToken() (plaintext, bcryptHash, sha256Lookup string, err error) {
	bytes := make([]byte, TokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	plaintext = TokenPrefix + base64.RawURLEncoding.EncodeToString(bytes)

	// SHA-256 for fast GSI lookup
	sha256Sum := sha256.Sum256([]byte(plaintext))
	sha256Lookup = base64.RawURLEncoding.EncodeToString(sha256Sum[:])

	// bcrypt for secure storage
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcryptCost)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to hash token: %w", err)
	}

	return plaintext, string(hashBytes), sha256Lookup, nil
}

// CreateToken stores a new API token for a user
// Returns the token metadata
func (d *DynamoDB) CreateToken(ctx context.Context, userID, name, tokenHash, tokenLookup string) (*APIToken, error) {
	token := &APIToken{
		UserID:      userID,
		TokenID:     uuid.New().String(),
		TokenHash:   tokenHash,
		TokenLookUp: tokenLookup,
		Name:        name,
		CreatedAt:   time.Now().UTC(),
	}

	item, err := attributevalue.MarshalMap(token)
	if err != nil {
		return nil, fmt.Errorf("%w:%w", ErrMarshalToken, err)
	}
	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(d.apiTokensTable),
		Item:      item,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	return token, nil
}

func (d *DynamoDB) ValidateToken(ctx context.Context, plaintext string) (*APIToken, error) {
	// Compute SHA-256 lookup key
	sha256Sum := sha256.Sum256([]byte(plaintext))
	lookup := base64.RawURLEncoding.EncodeToString(sha256Sum[:])

	// Query GSI by lookup key
	result, err := d.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(d.apiTokensTable),
		IndexName:              aws.String("TokenLookupIndex"),
		KeyConditionExpression: aws.String("tokenLookup = :lookup"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":lookup": &types.AttributeValueMemberS{Value: lookup},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrQueryTokens, err)
	}

	if len(result.Items) == 0 {
		return nil, ErrTokenNotFound
	}

	var token APIToken
	if err := attributevalue.UnmarshalMap(result.Items[0], &token); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrQueryTokens, err)
	}

	// Verify bcrypt hash (double-check)
	if err := bcrypt.CompareHashAndPassword([]byte(token.TokenHash), []byte(plaintext)); err != nil {
		return nil, ErrInvalidToken
	}

	return &token, nil
}

// ListTokensByUser returns all tokens for a user (without hashes).
func (d *DynamoDB) ListTokensByUser(ctx context.Context, userID string) ([]APIToken, error) {
	result, err := d.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(d.apiTokensTable),
		KeyConditionExpression: aws.String("userId = :uid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uid": &types.AttributeValueMemberS{Value: userID},
		},
		ProjectionExpression: aws.String("userId, tokenId, #n, createdAt, lastUsedAt"),
		ExpressionAttributeNames: map[string]string{
			"#n": "name",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrQueryTokens, err)
	}

	var tokens []APIToken
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &tokens); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrQueryTokens, err)
	}

	return tokens, nil
}

// DeleteToken deletes a token by ID
func (d *DynamoDB) DeleteToken(ctx context.Context, userID, tokenID string) error {
	_, err := d.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(d.apiTokensTable),
		Key: map[string]types.AttributeValue{
			"userId":  &types.AttributeValueMemberS{Value: userID},
			"tokenId": &types.AttributeValueMemberS{Value: tokenID},
		},
	})
	if err != nil {
		return fmt.Errorf("%w:%w", ErrDeleteToken, err)
	}
	return nil
}

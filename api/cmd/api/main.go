package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/ElshadHu/vulnly/api/internal/handler"
	"github.com/ElshadHu/vulnly/api/internal/repository"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create a root context
	ctx := context.Background()

	repo, err := repository.NewDynamoDB(ctx)
	if err != nil {
		log.Fatalf("failed to create repository: %v", err)
	}
	r := setupRouter(repo)

	// Detect environment and run accordingly
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		adapter := ginadapter.New(r)
		lambda.Start(func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			return adapter.ProxyWithContext(ctx, req)
		})
		return
	}
	log.Println("Running locally on :8080")
	r.Run(":8080")
}

func setupRouter(repo *repository.DynamoDB) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
		h := handler.New(repo)
		api.POST("/ingest", h.Ingest)
		api.GET("/projects", h.ListProjects)
		api.GET("/projects/:project_id", h.GetProject)
		api.GET("/projects/:project_id/scans", h.ListScans)
	}
	return r
}

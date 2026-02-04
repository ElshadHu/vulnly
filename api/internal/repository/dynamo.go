// Package repository provides data access layer for DynamoDB operations.

package repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

var (
	ErrLoadConfig               = errors.New("failed to load AWS config")
	ErrMarshalProject           = errors.New("failed to marshal project")
	ErrMarshalScan              = errors.New("failed to marshal scan")
	ErrMarshalVulnerability     = errors.New("failed to marshal vulnerability")
	ErrPutProject               = errors.New("failed to put project")
	ErrQueryProjects            = errors.New("failed to query projects")
	ErrBatchWrite               = errors.New("failed to batch write")
	ErrQueryVulnerabilities     = errors.New("failed to query vulnerabilities")
	ErrUnmarshalVulnerabilities = errors.New("failed to unmarshal vulnerabilities")
	ErrQueryScans               = errors.New("failed to query scans")
	ErrUnmarshalScans           = errors.New("failed to unmarshal scans")
)

type DynamoDB struct {
	client               *dynamodb.Client
	usersTable           string
	projectsTable        string
	scansTable           string
	vulnerabilitiesTable string
	apiTokensTable       string
	trendsTable          string
}

// NewDynamoDB creates a new DynamoDB  repository with connections to all required tables
func NewDynamoDB(ctx context.Context) (*DynamoDB, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrLoadConfig, err)
	}
	return &DynamoDB{
		client:               dynamodb.NewFromConfig(cfg),
		usersTable:           os.Getenv("USERS_TABLE"),
		projectsTable:        os.Getenv("PROJECTS_TABLE"),
		scansTable:           os.Getenv("SCANS_TABLE"),
		vulnerabilitiesTable: os.Getenv("VULNERABILITIES_TABLE"),
		apiTokensTable:       os.Getenv("API_TOKENS_TABLE"),
		trendsTable:          os.Getenv("TRENDS_TABLE"),
	}, nil
}

// GetOrCreateProject retrieves an existing project by name or creates a new one if it doesn't exist
func (d *DynamoDB) GetOrCreateProject(ctx context.Context, userID, projectName string) (*Project, error) {
	projectID := uuid.New().String()
	project := &Project{
		UserID:    userID,
		ProjectID: projectID,
		Name:      projectName,
		CreatedAt: time.Now().UTC(),
	}
	item, err := attributevalue.MarshalMap(project)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrMarshalProject, err)
	}
	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(d.projectsTable),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(userId)"),
	})
	if err != nil {
		var cfe *types.ConditionalCheckFailedException
		if !errors.As(err, &cfe) {
			return nil, fmt.Errorf("%w: %w", ErrPutProject, err)
		}
	}
	return project, nil
}

// CreateScan creates a new scan for a project
func (d *DynamoDB) CreateScan(ctx context.Context, scan *Scan) error {
	item, err := attributevalue.MarshalMap(scan)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrMarshalScan, err)
	}
	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(d.scansTable),
		Item:      item,
	})
	return err
}

// CreateVulnerabilities batch writes vulnerabilities in chunk of 25 (DynamoDB limit)
func (d *DynamoDB) CreateVulnerabilities(ctx context.Context, vulns []Vulnerability) error {
	if len(vulns) == 0 {
		return nil
	}
	var requests []types.WriteRequest
	for _, v := range vulns {
		item, err := attributevalue.MarshalMap(v)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrMarshalVulnerability, err)
		}
		requests = append(requests, types.WriteRequest{
			PutRequest: &types.PutRequest{Item: item},
		})
	}
	for i := 0; i < len(requests); i += 25 {
		end := i + 25
		if end > len(requests) {
			end = len(requests)
		}
		_, err := d.client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				d.vulnerabilitiesTable: requests[i:end],
			},
		})
		if err != nil {
			return fmt.Errorf("%w: %w", ErrBatchWrite, err)
		}
	}
	return nil
}

// ListProjectsByUser returns all projects owned by a user
func (d *DynamoDB) ListProjectsByUser(ctx context.Context, userID string) ([]Project, error) {
	result, err := d.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(d.projectsTable),
		KeyConditionExpression: aws.String("userId = :uid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uid": &types.AttributeValueMemberS{Value: userID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrQueryProjects, err)
	}
	var projects []Project
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &projects); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrQueryProjects, err)
	}
	return projects, nil
}

// ListScansByProject return the most 50 recent scans for a project, newest first
func (d *DynamoDB) ListScansByProject(ctx context.Context, projectID string) ([]Scan, error) {
	result, err := d.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(d.scansTable),
		KeyConditionExpression: aws.String("projectId = :pid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pid": &types.AttributeValueMemberS{Value: projectID},
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(50),
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrQueryProjects, err)
	}
	var scans []Scan
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &scans); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrQueryProjects, err)
	}
	return scans, nil
}

func (d *DynamoDB) GetProjectByName(ctx context.Context, userID, name string) (*Project, error) {
	result, err := d.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(d.projectsTable),
		KeyConditionExpression: aws.String("userId = :uid"),
		FilterExpression:       aws.String("#n = :name"),
		ExpressionAttributeNames: map[string]string{
			"#n": "name", // name is a reserved word in DynamoDB
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uid":  &types.AttributeValueMemberS{Value: userID},
			":name": &types.AttributeValueMemberS{Value: name},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrQueryProjects, err)
	}
	if len(result.Items) == 0 {
		return nil, nil
	}

	var project Project
	if err := attributevalue.UnmarshalMap(result.Items[0], &project); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrQueryProjects, err)
	}
	return &project, nil
}

// GetProjectByID retrieves a project by userId and projectId
func (d *DynamoDB) GetProjectByID(ctx context.Context, userID, projectID string) (*Project, error) {
	result, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(d.projectsTable),
		Key: map[string]types.AttributeValue{
			"userId":    &types.AttributeValueMemberS{Value: userID},
			"projectId": &types.AttributeValueMemberS{Value: projectID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrQueryProjects, err)
	}
	if result.Item == nil {
		return nil, nil
	}

	var project Project
	if err := attributevalue.UnmarshalMap(result.Item, &project); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrQueryProjects, err)
	}
	return &project, nil
}

// GetScanByID retrieves a scan by scanId using GSI
func (d *DynamoDB) GetScanByID(ctx context.Context, scanID string) (*Scan, error) {
	result, err := d.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(d.scansTable),
		IndexName:              aws.String("ScanIdIndex"),
		KeyConditionExpression: aws.String("scanId = :sid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":sid": &types.AttributeValueMemberS{Value: scanID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrQueryScans, err)
	}
	if len(result.Items) == 0 {
		return nil, nil
	}

	var scan Scan
	if err := attributevalue.UnmarshalMap(result.Items[0], &scan); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnmarshalScans, err)
	}
	return &scan, nil
}

func (d *DynamoDB) ListVulnerabilitiesByScan(ctx context.Context, scanID, severity, packageName string, limit int) ([]Vulnerability, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(d.vulnerabilitiesTable),
		KeyConditionExpression: aws.String("scanId = :sid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":sid": &types.AttributeValueMemberS{Value: scanID},
		},
		Limit: aws.Int32(int32(limit)),
	}

	// Build filter expression for optional filters
	var filterParts []string
	if severity != "" {
		filterParts = append(filterParts, "severity = :sev")
		input.ExpressionAttributeValues[":sev"] = &types.AttributeValueMemberS{Value: severity}
	}
	if packageName != "" {
		filterParts = append(filterParts, "contains(packageName, :pkg)")
		input.ExpressionAttributeValues[":pkg"] = &types.AttributeValueMemberS{Value: packageName}
	}
	if len(filterParts) > 0 {
		input.FilterExpression = aws.String(strings.Join(filterParts, " AND "))
	}

	result, err := d.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrQueryVulnerabilities, err)
	}

	var vulns []Vulnerability
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &vulns); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnmarshalVulnerabilities, err)
	}
	return vulns, nil
}

// GetRecentScans returns scans from the last N days for a project, ordered oldest first.
// Used by the trends API to generate time-series vulnerability data for charts.
func (d *DynamoDB) GetRecentScans(ctx context.Context, projectID string, days int) ([]Scan, error) {
	startTime := time.Now().UTC().AddDate(0, 0, -days)

	result, err := d.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(d.scansTable),
		KeyConditionExpression: aws.String("projectId = :pid"),
		FilterExpression:       aws.String("createdAt >= :start"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pid":   &types.AttributeValueMemberS{Value: projectID},
			":start": &types.AttributeValueMemberS{Value: startTime.Format(time.RFC3339)},
		},
		ScanIndexForward: aws.Bool(true),
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrQueryScans, err)
	}

	var scans []Scan
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &scans); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnmarshalScans, err)
	}

	return scans, nil
}

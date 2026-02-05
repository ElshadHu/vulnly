package repository

import "time"

type User struct {
	UserID     string    `dynamodbav:"userId"`
	Email      string    `dynamodbav:"email"`
	CognitoSub string    `dynamodbav:"cognitoSub"`
	CreatedAt  time.Time `dynamodbav:"createdAt"`
}
type Project struct {
	UserID     string    `dynamodbav:"userId" json:"userId"`
	ProjectID  string    `dynamodbav:"projectId" json:"projectId"`
	Name       string    `dynamodbav:"name" json:"name"`
	CreatedAt  time.Time `dynamodbav:"createdAt" json:"createdAt"`
	LastScanAt time.Time `dynamodbav:"lastScanAt,omitempty" json:"lastScanAt,omitempty"`
}

type VulnSummary struct {
	Critical int `dynamodbav:"critical" json:"critical"`
	High     int `dynamodbav:"high" json:"high"`
	Medium   int `dynamodbav:"medium" json:"medium"`
	Low      int `dynamodbav:"low" json:"low"`
}

type Scan struct {
	ProjectID string      `dynamodbav:"projectId" json:"projectId"`
	ScanID    string      `dynamodbav:"scanId" json:"scanId"`
	Commit    string      `dynamodbav:"commit,omitempty" json:"commit,omitempty"`
	Branch    string      `dynamodbav:"branch,omitempty" json:"branch,omitempty"`
	Ecosystem string      `dynamodbav:"ecosystem" json:"ecosystem"`
	TotalDeps int         `dynamodbav:"totalDeps" json:"totalDeps"`
	Summary   VulnSummary `dynamodbav:"summary" json:"summary"`
	CreatedAt time.Time   `dynamodbav:"createdAt" json:"createdAt"`
}

type Vulnerability struct {
	ScanID         string `dynamodbav:"scanId" json:"scanId"`
	VulnID         string `dynamodbav:"vulnId" json:"vulnId"`
	PackageName    string `dynamodbav:"packageName" json:"packageName"`
	PackageVersion string `dynamodbav:"packageVersion" json:"packageVersion"`
	Severity       string `dynamodbav:"severity" json:"severity"`
	FixedVersion   string `dynamodbav:"fixedVersion,omitempty" json:"fixedVersion,omitempty"`
	Description    string `dynamodbav:"description,omitempty" json:"description,omitempty"`
}

// APIToken represents a long-lived token for CLI authentication
type APIToken struct {
	UserID      string    `dynamodbav:"userId"`
	TokenID     string    `dynamodbav:"tokenId"`
	TokenHash   string    `dynamodbav:"tokenHash"`
	TokenLookUp string    `dynamodbav:"tokenLookup"` //SHA-256 GSI lookup
	Name        string    `dynamodbav:"name"`
	CreatedAt   time.Time `dynamodbav:"createdAt"`
	LastUsedAt  time.Time `dynamodbav:"lastUsedAt,omitempty"`
}

// TrendPoint stores pre-aggregated vulnerability counts per day.
type TrendPoint struct {
	ProjectID string      `dynamodbav:"projectId"`
	Date      string      `dynamodbav:"date"` // YYYY-MM-DD
	Summary   VulnSummary `dynamodbav:"summary"`
	ScanCount int         `dynamodbav:"scanCount"`
	TTL       int64       `dynamodbav:"ttl,omitempty"`
}

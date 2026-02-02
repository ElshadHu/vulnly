package repository

import "time"

type User struct {
	UserID     string    `dynamodbav:"userId"`
	Email      string    `dynamodbav:"email"`
	CognitoSub string    `dynamodbav:"cognitoSub"`
	CreatedAt  time.Time `dynamodbav:"createdAt"`
}
type Project struct {
	UserID     string    `dynamodbav:"userId"`
	ProjectID  string    `dynamodbav:"projectId"`
	Name       string    `dynamodbav:"name"`
	CreatedAt  time.Time `dynamodbav:"createdAt"`
	LastScanAt time.Time `dynamodbav:"lastScanAt,omitempty"`
}

type VulnSummary struct {
	Critical int `dynamodbav:"critical" json:"critical"`
	High     int `dynamodbav:"high" json:"high"`
	Medium   int `dynamodbav:"medium" json:"medium"`
	Low      int `dynamodbav:"low" json:"low"`
}

type Scan struct {
	ProjectID string      `dynamodbav:"projectId"`
	ScanID    string      `dynamodbav:"scanId"`
	Commit    string      `dynamodbav:"commit,omitempty"`
	Branch    string      `dynamodbav:"branch,omitempty"`
	Ecosystem string      `dynamodbav:"ecosystem"`
	TotalDeps int         `dynamodbav:"totalDeps"`
	Summary   VulnSummary `dynamodbav:"summary"`
	CreatedAt time.Time   `dynamodbav:"createdAt"`
}

type Vulnerability struct {
	ScanID         string `dynamodbav:"scanId"`
	VulnID         string `dynamodbav:"vulnId"`
	PackageName    string `dynamodbav:"packageName"`
	PackageVersion string `dynamodbav:"packageVersion"`
	Severity       string `dynamodbav:"severity"`
	FixedVersion   string `dynamodbav:"fixedVersion,omitempty"`
	Description    string `dynamodbav:"description,omitempty"`
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

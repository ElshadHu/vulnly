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

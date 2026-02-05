// Frontend types for API responses
// These types match the JSON structure returned by backend handlers
// See: api/internal/handler/*.go for response formats

export type Severity = "CRITICAL" | "HIGH" | "MEDIUM" | "LOW" | "UNKNOWN";

// Domain Types (match DynamoDB models returned via handlers)
// Source: api/internal/repository/models.go

// GET /api/projects returns Project[] directly from DynamoDB
export type Project = {
  projectId: string;
  userId: string;
  name: string;
  createdAt: string;
  lastScanAt?: string;
};

// Embedded in Scan, used for vulnerability counts
export type VulnSummary = {
  critical: number;
  high: number;
  medium: number;
  low: number;
};

// Returned in GetProjectResponse and ListScansResponse
export type Scan = {
  scanId: string;
  projectId: string;
  commit?: string;
  branch?: string;
  ecosystem: string;
  totalDeps: number;
  summary: VulnSummary;
  createdAt: string;
};

// GET /api/vulnerabilities?scanId=xxx returns Vulnerability[]
export type Vulnerability = {
  vulnId: string;
  scanId: string;
  packageName: string;
  packageVersion: string;
  severity: Severity;
  fixedVersion?: string;
  description?: string;
};

// Handler Response Types (explicit JSON structures from handlers)
// Source: api/internal/handler/*.go

// GET /api/trends?projectId=xxx returns TrendDataPoint[]
// Source: handler/trends.go - TrendDataPoint struct
export type TrendDataPoint = {
  date: string;
  critical: number;
  high: number;
  medium: number;
  low: number;
  total: number;
  scanId: string;
};

// GET /api/tokens returns APIToken[] (snake_case from handler.TokenResponse)
// Source: handler/tokens.go - TokenResponse struct
export type APIToken = {
  token_id: string;
  name: string;
  created_at: string;
  last_used_at?: string;
};

// Request Types

// POST /api/tokens body
export type CreateTokenRequest = {
  name: string;
};

// Response Wrapper Types (match gin.H{} responses in handlers)

// POST /api/tokens - returns token once (snake_case)
// Source: handler/tokens.go - CreateTokenResponse struct
export type CreateTokenResponse = {
  token_id: string;
  token: string; // plaintext, shown once only
  name: string;
  created_at: string;
};

// GET /api/tokens - { "tokens": [...] }
export type ListTokensResponse = {
  tokens: APIToken[];
};

// GET /api/projects - { "projects": [...] }
export type ListProjectsResponse = {
  projects: Project[];
};

// GET /api/projects/:project_id - { "project_id": "...", "name": "...", "scans": [...] }
// Source: handler/projects.go line 44-48
export type GetProjectResponse = {
  project_id: string;
  name: string;
  scans: Scan[];
};

// GET /api/projects/:project_id/scans - same structure as GetProjectResponse
// Source: handler/projects.go line 72-76
export type ListScansResponse = GetProjectResponse;

// GET /api/trends?projectId=xxx - { "dataPoints": [...] }
// Source: handler/trends.go line 69
export type TrendsResponse = {
  dataPoints: TrendDataPoint[];
};

// GET /api/vulnerabilities?scanId=xxx - { "vulnerabilities": [...], "count": N }
// Source: handler/vulnerabilities.go line 70-73
export type VulnerabilitiesResponse = {
  vulnerabilities: Vulnerability[];
  count: number;
};

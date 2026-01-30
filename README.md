# Vulnly

Dependency vulnerability scanner for CI/CD pipelines.

## What it will do

- Scans npm,go and pip dependencies (other language dependencies can be added later If I feel like it and fall in love with projet)
- Checks against OSV.dev vulnerability database
- Fails CI builds on critical vulnerabilities
- Dashboard to track trends over time

## Architecture

```
CLI (vulnly scan .)
    |
    | parses lockfiles, queries OSV API
    |
    v
POST /api/ingest (planned)
    |
    | sends scan results
    |
    v
Backend API (Lambda + API Gateway + DynamoDB)
    |
    | stores scan history, projects, vulnerabilities
    |
    v
Dashboard (Next.js) -> planned
    |
    | fetches from API
    v
Shows projects, scan history, trend charts
```

### How it works

1. CLI scans locally: vulnly scan . parses lockfiles and queries OSV.dev
2. CLI uploads results: After scanning, CLI sends results to backend via --upload flag (in progress)
3. Backend stores data: Lambda receives scan data and stores in DynamoDB (Planned version let's see)
4. Dashboard displays trends: Frontend fetches from API and shows vulnerability history (Planned version)


## Features (That is what I have planned so far)

- [x] Parse npm, Go, and Python lockfiles
- [x] OSV.dev API integration with hydration
- [x] File-based caching with TTL
- [x] CVSS severity parsing
- [x] CLI scanner with vulnerability output
- [ ] AWS backend (Lambda + DynamoDB) - in progress
- [ ] Dashboard to track trends



## Installation

```bash
# Build locally
go build -o vulnly ./cmd/vulnly

# Or install to $GOPATH/bin
go install ./cmd/vulnly
```

## Usage

```bash
# Scan current directory (use ./vulnly if not installed)
vulnly scan .

# Scan specific directory
vulnly scan ./path/to/project

# JSON output
vulnly scan . --format json

# Fail on severity threshold (for CI)
vulnly scan . --fail-on-severity HIGH
vulnly scan . --fail-on-severity CRITICAL

# Check version
vulnly version
```

## Dependencies

| Package | Why |
|---------|-----|
| `github.com/spf13/cobra` | CLI framework for commands and flags |
| `github.com/olekukonko/tablewriter` | Pretty table output in terminal |
| `github.com/goark/go-cvss` | Parse CVSS vectors to get accurate severity scores |

### Backend Dependencies (api/)

| Package | Why |
|---------|-----|
| `github.com/gin-gonic/gin` | HTTP router for Lambda |
| `github.com/aws/aws-sdk-go-v2` | DynamoDB client |
| `github.com/aws/aws-lambda-go` | Lambda runtime |
| `github.com/awslabs/aws-lambda-go-api-proxy` | Gin adapter for Lambda |

## Roadmap

It could seem like roadmap is broken but that is just the way it is, I need to go around more in order to figure out more :)
During the process of development, I will update the roadmap

- [x] CLI scanner (npm, Go, pip)
- [x] Lockfile parsers (npm, Go, pip)
- [x] OSV API client with batch queries
- [x] Severity parsing with CVSS
- [x] End-to-end scan with vulnerability output
- [x] OSV.dev integration complete
- [ ] AWS SAM infrastructure - in progress
- [ ] Backend API handlers - in progress
- [ ] CLI --upload flag
- [ ] Dashboard with trends
- [ ] Email alerts
- [ ] Slack integration (v2)



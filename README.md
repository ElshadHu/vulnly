# Vulnly

Dependency vulnerability scanner for CI/CD pipelines.

## What it will do

- Scans npm,go and pip dependencies (other language dependencies can be added later If I feel like it and fall in love with projet)
- Checks against OSV.dev vulnerability database
- Fails CI builds on critical vulnerabilities
- Dashboard to track trends over time

## Features (That is what I have planned so far)

- [+] Parse npm, Go, and Python lockfiles
- [+] OSV.dev API integration with hydration
- [+] File-based caching with TTL
- [+] CVSS severity parsing
- [ ] CLI scanner with vulnerability output
- [ ] AWS backend (Lambda + DynamoDB)
- [ ] Dashboard to track trends



## Usage

```bash
# Scan current directory
vulnly scan .

# JSON output
vulnly scan . --format json

# Fail on high severity vulnerabilities
vulnly scan . --fail-on-severity HIGH

# Output as table (default)
vulnly scan . --format table
```

## Architecture

TODO: Will be added with proper diagrams

## Dependencies

| Package | Why |
|---------|-----|
| `github.com/spf13/cobra` | CLI framework for commands and flags |
| `github.com/olekukonko/tablewriter` | Pretty table output in terminal |
| `github.com/goark/go-cvss` | Parse CVSS vectors to get accurate severity scores |

## Roadmap

It could seem like roadmap is broken but that is just the way it is, I need to go around more in order to figure out more :)
During the process of development, I will update the roadmap

- [x] CLI scanner (npm, Go, I will see how it goes)
- [x] Lockfile parsers (npm, Go, pip)
- [x] OSV API client with batch queries (currenntly I am doing this)
- [x] Severity parsing with CVSS
- [ ] End-to-end scan with vulnerability output
- [ ] OSV.dev integration complete
- [ ] AWS SAM infrastructure
- [ ] Dashboard with trends
- [ ] Email alerts
- [ ] Slack integration (v2)



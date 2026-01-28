package osv

// Package represents a package in an OSV query
type Package struct {
	Name      string `json:"name"`
	Ecosystem string `json:"ecosystem"`
}

// Query represents a single query to OSV
type Query struct {
	Package Package `json:"package"`
	Version string  `json:"version"`
}

// BatchedQuery represents multiple queries
type BatchedQuery struct {
	Queries []*Query `json:"queries"`
}

// MinimalVulnerability is returned by /v1/querybatch
type MinimalVulnerability struct {
	ID       string `json:"id"`
	Modified string `json:"modified"`
}

// MinimalResponse contains minimal vulns for a single query
type MinimalResponse struct {
	Vulns []MinimalVulnerability `json:"vulns"`
}

// BatchedResponse is the response from /v1/querybatch
type BatchedResponse struct {
	Results []MinimalResponse `json:"results"`
}

type SeverityEntry struct {
	Type  string `json:"type"`
	Score string `json:"score"`
}

type Affected struct {
	Package  Package  `json:"package"`
	Ranges   []Range  `json:"ranges,omitempty"`
	Versions []string `json:"versions,omitempty"`
}

type Range struct {
	Type   string  `json:"type"`
	Events []Event `json:"events"`
}
type Event struct {
	Introduced string `json:"introduced,omitempty"`
	Fixed      string `json:"fixed,omitempty"`
}

type Reference struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

// Vulnerability is the full vulnerability from /v1/vulns/{id}
type Vulnerability struct {
	ID         string          `json:"id"`
	Summary    string          `json:"summary"`
	Details    string          `json:"details"`
	Severity   []SeverityEntry `json:"severity,omitempty"`
	Affected   []Affected      `json:"affected"`
	References []Reference     `json:"references,omitempty"`
}

func (a *Affected) MatchesPackage(ecosystem, name string) bool {
	return a.Package.Name == name && a.Package.Ecosystem == ecosystem
}

func (r *Range) FixedVersion() string {
	for _, e := range r.Events {
		if e.Fixed != "" {
			return e.Fixed
		}
	}
	return ""
}

func (a *Affected) FixedVersion() string {
	for _, r := range a.Ranges {
		if v := r.FixedVersion(); v != "" {
			return v
		}
	}
	return ""
}

func (v *Vulnerability) GetFixedVersion(ecosystem, name string) string {
	for _, aff := range v.Affected {
		if aff.MatchesPackage(ecosystem, name) {
			return aff.FixedVersion()
		}
	}
	return ""
}

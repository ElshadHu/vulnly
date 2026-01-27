package osv

import (
	"strings"

	cvss3 "github.com/goark/go-cvss/v3/metric"
)

type Severity string

const (
	SeverityCritical Severity = "CRITICAL"
	SeverityHigh     Severity = "HIGH"
	SeverityMedium   Severity = "MEDIUM"
	SeverityLow      Severity = "LOW"
	SeverityUnknown  Severity = "UNKNOWN"
)

func (s Severity) Priority() int {
	switch s {
	case SeverityCritical:
		return 4
	case SeverityHigh:
		return 3
	case SeverityMedium:
		return 2
	case SeverityLow:
		return 1
	default:
		return 0
	}
}

func ParseSeverityFromScore(score float64) Severity {
	switch {
	case score >= 9.0:
		return SeverityCritical
	case score >= 7.0:
		return SeverityHigh
	case score >= 4.0:
		return SeverityMedium
	case score > 0:
		return SeverityLow
	default:
		return SeverityUnknown
	}
}

// ParseCVSSVector parses a CVSS v3 vector and returns the score
func ParseCVSSVector(vector string) (float64, error) {
	bm, err := cvss3.NewBase().Decode(vector)
	if err != nil {
		return 0, err
	}
	return bm.Score(), nil
}

// GetSeverity extracts severity from a Vulnerability
func GetSeverity(v *Vulnerability) Severity {
	if v == nil {
		return SeverityUnknown
	}

	for _, s := range v.Severity {
		if strings.HasPrefix(s.Type, "CVSS_V3") {
			score, err := ParseCVSSVector(s.Score)
			if err != nil {
				continue
			}
			return ParseSeverityFromScore(score)
		}
	}
	for _, s := range v.Severity {
		if strings.HasPrefix(s.Type, "CVSS_V2") {
			return estimateSeverityFromVector(s.Score)
		}
	}

	return SeverityUnknown
}

// estimateSeverityFromVector  which is a fallback for non-v3 CVSS
func estimateSeverityFromVector(vector string) Severity {
	if vector == "" {
		return SeverityUnknown
	}

	highCount := strings.Count(vector, ":H")
	lowCount := strings.Count(vector, ":L")
	isNetwork := strings.Contains(vector, "AV:N")
	noPrivRequired := strings.Contains(vector, "PR:N") || strings.Contains(vector, "Au:N")

	switch {
	case highCount >= 3 && isNetwork:
		return SeverityCritical
	case highCount >= 2 || (highCount >= 1 && isNetwork && noPrivRequired):
		return SeverityHigh
	case highCount >= 1 || lowCount >= 2:
		return SeverityMedium
	case lowCount >= 1:
		return SeverityLow
	default:
		return SeverityUnknown
	}
}

package osv

import "github.com/ElshadHu/vulnly/internal/lockfile"

func Scan(packages []lockfile.PackageDetails) (*ScanResult, error) {
	if len(packages) == 0 {
		return &ScanResult{}, nil
	}
	query := buildBatchQuery(packages)
	resp, err := MakeRequest(query)
	if err != nil {
		return nil, err
	}
	hydrated, err := Hydrate(resp)
	if err != nil {
		return nil, err
	}
	return buildResult(packages, hydrated), nil
}

func buildBatchQuery(packages []lockfile.PackageDetails) BatchedQuery {
	queries := make([]*Query, len(packages))
	for i, pkg := range packages {
		queries[i] = &Query{
			Package: Package{Name: pkg.Name, Ecosystem: pkg.Ecosystem},
			Version: pkg.Version,
		}
	}
	return BatchedQuery{Queries: queries}
}

func buildResult(packages []lockfile.PackageDetails, hydrated *HydratedBatchedResponse) *ScanResult {
	var vulns []PackageVuln
	summary := Summary{TotalDeps: len(packages)}

	for i, pkg := range packages {
		for _, v := range hydrated.Results[i].Vulns {
			if v == nil {
				continue
			}
			sev := GetSeverity(v)
			vulns = append(vulns, PackageVuln{
				Name:       pkg.Name,
				Version:    pkg.Version,
				Ecosystem:  pkg.Ecosystem,
				Vuln:       v,
				Severity:   sev,
				FixVersion: v.GetFixedVersion(pkg.Ecosystem, pkg.Name),
			})
			incrementSummary(&summary, sev)
		}
	}
	return &ScanResult{Packages: vulns, Summary: summary}
}

func incrementSummary(s *Summary, sev Severity) {
	switch sev {
	case SeverityCritical:
		s.Critical++
	case SeverityHigh:
		s.High++
	case SeverityMedium:
		s.Medium++
	case SeverityLow:
		s.Low++
	default:
		s.Unknown++
	}
}

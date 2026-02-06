package output

import (
	"fmt"
	"io"

	"github.com/ElshadHu/vulnly/internal/osv"
	"github.com/olekukonko/tablewriter"
)

func TableResult(w io.Writer, result *osv.ScanResult) error {
	total := result.Summary.Critical + result.Summary.High + result.Summary.Medium + result.Summary.Low + result.Summary.Unknown
	if _, err := fmt.Fprintf(w, "\nScan: %d deps, %d vulns\n", result.Summary.TotalDeps, total); err != nil {
		return fmt.Errorf("failed to write summary: %w", err)
	}
	if _, err := fmt.Fprintf(w, "Critical: %d | High: %d | Medium: %d | Low: %d | Unknown: %d\n\n",
		result.Summary.Critical, result.Summary.High, result.Summary.Medium, result.Summary.Low, result.Summary.Unknown); err != nil {
		return fmt.Errorf("failed to write severity breakdown: %w", err)
	}

	if len(result.Packages) == 0 {
		if _, err := fmt.Fprintln(w, "No vulnerabilities found"); err != nil {
			return fmt.Errorf("failed to write empty result %w", err)
		}
		return nil
	}

	table := tablewriter.NewWriter(w)
	table.Header([]string{"Package", "Version", "Vuln ID", "Severity", "Fix"})

	for _, pv := range result.Packages {
		fix := pv.FixVersion
		if fix == "" {
			fix = "-"
		}
		if err := table.Append([]string{pv.Name, pv.Version, pv.Vuln.ID, string(pv.Severity), fix}); err != nil {
			return fmt.Errorf("failed to append table: %w", err)
		}
	}
	if err := table.Render(); err != nil {
		return fmt.Errorf("failed to render table: %w", err)
	}
	return nil
}

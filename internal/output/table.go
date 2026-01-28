package output

import (
	"fmt"
	"io"

	"github.com/ElshadHu/vulnly/internal/osv"
	"github.com/olekukonko/tablewriter"
)

func TableResult(w io.Writer, result *osv.ScanResult) {
	total := result.Summary.Critical + result.Summary.High + result.Summary.Medium + result.Summary.Low + result.Summary.Unknown
	fmt.Fprintf(w, "\nScan: %d deps, %d vulns\n", result.Summary.TotalDeps, total)
	fmt.Fprintf(w, "Critical: %d | High: %d | Medium: %d | Low: %d | Unknown: %d\n\n",
		result.Summary.Critical, result.Summary.High, result.Summary.Medium, result.Summary.Low, result.Summary.Unknown)

	if len(result.Packages) == 0 {
		fmt.Fprintln(w, "No vulnerabilities found")
		return
	}

	table := tablewriter.NewWriter(w)
	table.Header([]string{"Package", "Version", "Vuln ID", "Severity", "Fix"})

	for _, pv := range result.Packages {
		fix := pv.FixVersion
		if fix == "" {
			fix = "-"
		}
		table.Append([]string{pv.Name, pv.Version, pv.Vuln.ID, string(pv.Severity), fix})
	}
	table.Render()
}

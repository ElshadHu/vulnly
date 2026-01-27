package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ElshadHu/vulnly/internal/lockfile"
	"github.com/ElshadHu/vulnly/internal/osv"
	"github.com/ElshadHu/vulnly/internal/output"
	"github.com/spf13/cobra"
)

func NewScanCmd() *cobra.Command {
	var failOnSeverity string

	cmd := &cobra.Command{
		Use:   "scan [directory]",
		Short: "Scan directory for dependency vulnerabilities",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) > 0 {
				dir = args[0]
			}

			packages, err := lockfile.ExtractAll(dir)
			if err != nil {
				return fmt.Errorf("failed to extract packages: %w", err)
			}

			if len(packages) == 0 {
				fmt.Println("No dependencies found")
				return nil
			}

			var vulns []*osv.Vulnerability // placeholder for now

			// Output results
			switch formatFlag {
			case "json":
				if err := output.JSON(os.Stdout, packages); err != nil {
					return err
				}
			default:
				output.Table(os.Stdout, packages)
			}

			// Check severity threshold after output
			if checkSeverityThreshold(vulns, failOnSeverity) {
				os.Exit(1)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&failOnSeverity, "fail-on-severity", "",
		"Exit with code 1 if vulnerabilities at or above this severity (CRITICAL, HIGH, MEDIUM, LOW)")

	return cmd
}

func checkSeverityThreshold(vulns []*osv.Vulnerability, threshold string) bool {
	if threshold == "" {
		return false
	}
	thresholdSev := osv.Severity(strings.ToUpper(threshold))
	for _, v := range vulns {
		severity := osv.GetSeverity(v)
		if severity.Priority() >= thresholdSev.Priority() {
			return true
		}
	}
	return false
}

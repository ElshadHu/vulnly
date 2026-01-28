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

			result, err := osv.Scan(packages)
			if err != nil {
				return fmt.Errorf("failed to query OSV: %w", err)
			}

			// Output results
			switch formatFlag {
			case "json":
				if err := output.JSONResult(os.Stdout, result); err != nil {
					return err
				}
			default:
				output.TableResult(os.Stdout, result)
			}

			// Check severity threshold after output
			if checkSeverityThreshold(result, failOnSeverity) {
				os.Exit(1)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&failOnSeverity, "fail-on-severity", "",
		"Exit with code 1 if vulnerabilities at or above this severity (CRITICAL, HIGH, MEDIUM, LOW)")

	return cmd
}

func checkSeverityThreshold(result *osv.ScanResult, threshold string) bool {
	if threshold == "" {
		return false
	}
	thresholdSev := osv.Severity(strings.ToUpper(threshold))
	for _, pv := range result.Packages {
		if pv.Severity.Priority() >= thresholdSev.Priority() {
			return true
		}
	}
	return false
}

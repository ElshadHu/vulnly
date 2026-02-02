package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ElshadHu/vulnly/internal/lockfile"
	"github.com/ElshadHu/vulnly/internal/osv"
	"github.com/ElshadHu/vulnly/internal/output"
	"github.com/ElshadHu/vulnly/internal/upload"
	"github.com/spf13/cobra"
)

var (
	tokenFlag   string
	apiURLFlag  string
	projectFlag string
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
				if err := output.JSONResult(cmd.OutOrStdout(), result); err != nil {
					return fmt.Errorf("failed to output JSON: %w", err)
				}
			default:
				if err := output.TableResult(cmd.OutOrStdout(), result); err != nil {
					return fmt.Errorf("failed to output table: %w", err)
				}
			}
			// Upload results if token provided
			if tokenFlag != "" {
				if apiURLFlag == "" {
					return fmt.Errorf("api-url is required when using --token")
				}
				project := projectFlag
				if project == "" {
					absDir, _ := filepath.Abs(dir)
					project = filepath.Base(absDir)
				}
				if err := uploadResults(tokenFlag, apiURLFlag, project, result); err != nil {
					return fmt.Errorf("failed to upload results: %w", err)
				}
				fmt.Printf("\nResults uploaded successfully %s\n", apiURLFlag)

			}

			// Check severity threshold after output
			if checkSeverityThreshold(result, failOnSeverity) {
				return fmt.Errorf("vulnerabilities found at or above %s severity", failOnSeverity)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&failOnSeverity, "fail-on-severity", "",
		"Exit with code 1 if vulnerabilities at or above this severity (CRITICAL, HIGH, MEDIUM, LOW)")
	cmd.Flags().StringVar(&tokenFlag, "token", os.Getenv("VULNLY_TOKEN"), "JWT token for API authentication")
	cmd.Flags().StringVar(&apiURLFlag, "api-url", os.Getenv("VULNLY_API_URL"), "API base URL")
	cmd.Flags().StringVar(&projectFlag, "project", "", "Project name (defaults to directory name)")

	return cmd
}

func uploadResults(token, apiURL, project string, result *osv.ScanResult) error {
	client := upload.NewClient(apiURL, token)

	req := &upload.IngestRequest{
		Project:   project,
		Ecosystem: detectEcosystem(result),
		Summary: upload.VulnSummary{
			Total:    result.Summary.TotalDeps,
			Critical: result.Summary.Critical,
			High:     result.Summary.High,
			Medium:   result.Summary.Medium,
			Low:      result.Summary.Low,
		},
	}
	for _, pv := range result.Packages {
		if pv.Vuln != nil {
			req.Vulnerabilities = append(req.Vulnerabilities, upload.Vulnerability{
				ID:           pv.Vuln.ID,
				Package:      pv.Name,
				Version:      pv.Version,
				Severity:     string(pv.Severity),
				FixedVersion: pv.FixVersion,
				Description:  pv.Vuln.Summary,
			})
		}
	}
	_, err := client.Ingest(req)
	return err
}

func detectEcosystem(result *osv.ScanResult) string {
	ecosystems := make(map[string]int)
	for _, pv := range result.Packages {
		ecosystems[pv.Ecosystem]++
	}
	maxCount := 0
	ecosystem := "unknown"
	for e, count := range ecosystems {
		if count > maxCount {
			maxCount = count
			ecosystem = e
		}
	}
	return ecosystem
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

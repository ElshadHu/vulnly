package cmd

import (
	"fmt"
	"os"

	"github.com/ElshadHu/vulnly/internal/lockfile"
	"github.com/ElshadHu/vulnly/internal/output"
	"github.com/spf13/cobra"
)

func NewScanCmd() *cobra.Command {
	return &cobra.Command{
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

			switch formatFlag {
			case "json":
				return output.JSON(os.Stdout, packages)
			default:
				output.Table(os.Stdout, packages)
			}

			return nil
		},
	}
}

package cmd

import (
	"github.com/spf13/cobra"
)

var (
	formatFlag string
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "vulnly",
		Short: "Dependency vulnerability scanner",
	}
	rootCmd.PersistentFlags().StringVarP(&formatFlag, "format", "f", "table", "Output format (json,table)")
	rootCmd.AddCommand(NewScanCmd())
	rootCmd.AddCommand(NewVersionCmd())

	return rootCmd
}

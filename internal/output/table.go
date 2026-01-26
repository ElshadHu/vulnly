package output

import (
	"io"

	"github.com/ElshadHu/vulnly/internal/lockfile"
	"github.com/olekukonko/tablewriter"
)

func Table(w io.Writer, packages []lockfile.PackageDetails) {
	table := tablewriter.NewWriter(w)
	table.Header([]string{"Name", "Version", "Ecosystem"})

	for _, pkg := range packages {
		table.Append([]string{pkg.Name, pkg.Version, pkg.Ecosystem})
	}
	table.Render()
}

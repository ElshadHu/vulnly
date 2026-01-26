package output

import (
	"encoding/json"
	"io"

	"github.com/ElshadHu/vulnly/internal/lockfile"
)

func JSON(w io.Writer, packages []lockfile.PackageDetails) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	return encoder.Encode(packages)
}

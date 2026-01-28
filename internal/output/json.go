package output

import (
	"encoding/json"
	"io"

	"github.com/ElshadHu/vulnly/internal/osv"
)

func JSONResult(w io.Writer, result *osv.ScanResult) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

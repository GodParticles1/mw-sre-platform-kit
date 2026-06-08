package report

import (
	"encoding/json"
	"io"

	"mw-sre-platform/internal/core"
)

func WriteJSON(w io.Writer, r core.CheckReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

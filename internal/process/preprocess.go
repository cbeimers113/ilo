package process

import (
	"cbeimers113/ilo/internal/constant"
	"strings"
)

// Preprocess cleans the raw input source code by converting
// Esperanto characters to their x-mode digraphs.
func Preprocess(data string) string {
	for eo, x := range constant.Orthography {
		data = strings.ReplaceAll(data, string(eo), x)
	}

	return data
}

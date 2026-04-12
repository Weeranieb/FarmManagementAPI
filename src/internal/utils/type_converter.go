package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// ConvertRepeatedFormInts converts repeated multipart / url-encoded form values to []int.
// raw is typically form.Value[fieldName] (e.g. several selectedPondIds entries).
func ConvertRepeatedFormInts(fieldLabel string, raw []string) ([]int, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("%s is required", fieldLabel)
	}
	out := make([]int, 0, len(raw))
	for _, s := range raw {
		s = strings.TrimSpace(s)
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("invalid %s: %s", fieldLabel, s)
		}
		out = append(out, n)
	}
	return out, nil
}

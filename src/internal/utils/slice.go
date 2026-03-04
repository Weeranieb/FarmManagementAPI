package utils

import "slices"

// AppendStringIfMissing returns a slice that contains elem: if slice is nil, returns []string{elem};
// otherwise if slice already contains elem returns slice unchanged; otherwise returns append(slice, elem).
func AppendStringIfMissing(slice []string, elem string) []string {
	if slice == nil {
		return []string{elem}
	}
	if slices.Contains(slice, elem) {
		return slice
	}
	return append(slice, elem)
}

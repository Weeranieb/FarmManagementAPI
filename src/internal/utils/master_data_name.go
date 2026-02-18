package utils

import "strings"

// Thai display prefixes (must match frontend). Stored names are without prefix.
const (
	THFarmDisplayPrefix = "ฟาร์ม "
	THPondDisplayPrefix = "บ่อ "
)

// NormalizeFarmNameForStore trims the Thai farm prefix and surrounding whitespace
// so only the name is stored. e.g. "  ฟาร์ม 1 " -> "1".
func NormalizeFarmNameForStore(s string) string {
	return trimPrefixAndSpace(s, THFarmDisplayPrefix)
}

// NormalizePondNameForStore trims the Thai pond prefix and surrounding whitespace
// so only the name is stored. e.g. " บ่อ 1 ซ้าย " -> "1 ซ้าย".
func NormalizePondNameForStore(s string) string {
	return trimPrefixAndSpace(s, THPondDisplayPrefix)
}

// trimPrefixAndSpace trims surrounding whitespace and the given prefix (with or without trailing space).
// e.g. "ฟาร์ม 1" or "  ฟาร์ม 1  " -> "1".
func trimPrefixAndSpace(s, prefix string) string {
	t := strings.TrimSpace(s)
	p := strings.TrimSpace(prefix)

	t = strings.TrimPrefix(t, p)
	return strings.TrimSpace(t)
}

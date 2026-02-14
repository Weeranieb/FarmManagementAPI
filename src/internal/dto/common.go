package dto

// DropdownItem is a common key-value pair for dropdown/select options.
type DropdownItem struct {
	Key   int    `json:"key"`
	Value string `json:"value"`
}

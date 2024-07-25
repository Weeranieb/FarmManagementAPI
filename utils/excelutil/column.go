package excelutil

func ColName(index int) string {
	column := make([]byte, 0)
	for index >= 0 {
		column = append(column, byte('A'+(index%26)))
		index = index/26 - 1
	}
	// Reverse the column name
	for i, j := 0, len(column)-1; i < j; i, j = i+1, j-1 {
		column[i], column[j] = column[j], column[i]
	}
	return string(column)
}

package utils

// Offset to get offset from page and limit, min value for page = 1
func Offset(page, limit int) int {
	offset := (page - 1) * limit
	if offset < 0 {
		return 0
	}
	return offset
}

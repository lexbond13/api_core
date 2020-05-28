package util

// ConvertCountBytesToCountMegabytes
func ConvertCountBytesToCountMegabytes(countBytes int64) int64 {
	if countBytes <= 0 {
		return 0
	}
	return countBytes / 1048576
}

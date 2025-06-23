package helper

// StringValue mengubah pointer string menjadi string biasa
func StringValue(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

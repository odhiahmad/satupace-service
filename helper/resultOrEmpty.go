package helper

func ResultOrEmpty[T any](data []T) []T {
	if data == nil {
		return []T{}
	}
	return data
}

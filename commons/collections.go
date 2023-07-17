package commons

func IsFull[T any](slice []T) bool {
	return len(slice) == cap(slice)
}

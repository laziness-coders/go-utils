package generic

// ToPointer returns a pointer to the given value.
// For example, ToPointer(42) returns *int pointing to 42.
func ToPointer[T any](v T) *T {
	return &v
}

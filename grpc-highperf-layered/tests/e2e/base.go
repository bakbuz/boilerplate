package e2e_test

func ptr[T any](v T) *T {
	return &v
}

package util

type transform[T1 any, T2 any] func(T1) T2

func Map[T1 any, T2 any](slice []T1, transform transform[T1, T2]) []T2 {
	result := make([]T2, len(slice))
	for i, v := range slice {
		result[i] = transform(v)
	}
	return result
}

package fp

// Map applies a transformation function to each element of the input slice
// and returns a new slice with the transformed elements.
func Map[T any, U any](input []T, transform func(T) U) []U {
	result := make([]U, len(input))
	for i, v := range input {
		result[i] = transform(v)
	}
	return result
}

// Filter returns a new slice containing only the elements of the input slice
// that satisfy the provided predicate function.
func Filter[T any](input []T, predicate func(T) bool) []T {
	result := []T{}
	for _, v := range input {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Reduce aggregates the elements of the input slice into a single value using
// the provided reduce function and an initial accumulator value.
func Reduce[T any, U any](input []T, initial U, reduceFunc func(U, T) U) U {
	accumulator := initial
	for _, v := range input {
		accumulator = reduceFunc(accumulator, v)
	}
	return accumulator
}

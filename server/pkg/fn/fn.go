package fn

type MapFunc[E any] func(e E) E

// Map applies the given function to each item of the slice and returns a slice of the results.
func Map[S ~[]E, E any](s S, f MapFunc[E]) S {
	result := make(S, len(s))
	for i := range s {
		result[i] = f(s[i])
	}
	return result
}

type keepFunc[E any] func(e E) bool

// Filter filters the iterable by applying the given function to each item, returning only the elements for which the function returns true.
func Filter[S ~[]E, E any](s S, f keepFunc[E]) S {
	result := S{}
	for _, v := range s {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

type reduceFunc[E any] func(cur, next E) E

// Reduce applies the given function cumulatively to the elements of the slice, from left to right, to reduce the slice to a single value.
func Reduce[E any](s []E, initialValue E, f reduceFunc[E]) E {
	result := initialValue
	for _, v := range s {
		result = f(result, v)
	}
	return result
}

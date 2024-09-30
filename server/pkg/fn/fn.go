package fn

// Map applies a function to each element of a slice and returns a new slice containing the results.
// The resulting slice must be of the same type as the input slice.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	squares := Map(numbers, func(x int) int { return x * x })
//	fmt.Println(squares) // Output: [1 4 9 16 25]
//
// Args:
//   - s: The input slice.
//   - f: The function to apply to each element.
//
// Returns:
//   - A new slice containing the results of applying f to each element of s.
func Map[S ~[]E, E any](s S, f func(E) E) S {
	result := make(S, len(s))
	for i := range s {
		result[i] = f(s[i])
	}
	return result
}

// Select applies a function to each element of a slice and returns a new slice containing the results.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	strings := fn.Select(numbers, func(x int) string { return fmt.Sprintf("%d", x) })
//	fmt.Println(strings) // Output: [1 2 3 4 5]
//
// Args:
//   - s: The input slice.
//   - f: The function to apply to each element.
//
// Returns:
//   - A new slice containing the results of applying f to each element of s.
func Select[E, R any](s []E, f func(E) R) []R {
	result := make([]R, len(s))
	for i := range s {
		result[i] = f(s[i])
	}
	return result
}

// Filter applies a predicate function to each element of a slice and returns a new slice containing only the elements for which the predicate returns true.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	evenNumbers := Filter(numbers, func(x int) bool { return x%2 == 0 })
//	fmt.Println(evenNumbers) // Output: [2 4]
//
// Args:
//   - s: The input slice.
//   - f: The predicate function to apply to each element.
//
// Returns:
//   - A new slice containing only the elements for which the predicate returns true.
func Filter[S ~[]E, E any](s S, f func(E) bool) S {
	result := S{}
	for _, v := range s {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

// Reduce applies a function to each element of a slice and an accumulator, returning the final accumulator value.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	sum := Reduce(numbers, 0, func(cur, next int) int { return acc + next })
//	fmt.Println(sum) // Output: 15
//
// Args:
//   - s: The input slice.
//   - initialValue: The initial value of the accumulator.
//   - f: The function to apply to the accumulator and each element of the slice.
//
// Returns:
//   - The final value of the accumulator.
func Reduce[E any](s []E, initialValue E, f func(acc, next E) E) E {
	accumulator := initialValue
	for _, v := range s {
		accumulator = f(accumulator, v)
	}
	return accumulator
}

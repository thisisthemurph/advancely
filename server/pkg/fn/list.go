package fn

import (
	"errors"
	"sort"
)

var (
	ErrEmptyList          = errors.New("empty list")
	ErrElementNotFound    = errors.New("element not found")
	ErrMoreThanOneElement = errors.New("more than one element")
	ErrIndexOutOfRange    = errors.New("index out of range")
)

type List[E comparable] struct {
	elements []E
}

// NewList instantiates a new List with the given elements.
//
//	list1 := fn.NewList(1, 2, 3)
//	list2 := fn.NewList(mySlice...)
func NewList[E comparable](elements ...E) *List[E] {
	return &List[E]{
		elements: elements,
	}
}

// NewEmptyList instantiates an empty List of type E.
//
//	list := fn.NewEmptyList[int]()
func NewEmptyList[E comparable]() *List[E] {
	return &List[E]{
		elements: []E{},
	}
}

// Add adds multiple elements to the end of the List.
//
//	numbers := int[]{4, 5, 6}
//	list := fn.NewList(1, 2, 3)
//	list.Add(numbers...)
//	list.Add(7)
//	list.Add(8, 9, 10)
func (l *List[E]) Add(elements ...E) {
	l.elements = append(l.elements, elements...)
}

// Insert inserts a new element at the given index.
//
// An error is returned if the index is out of range.
func (l *List[E]) Insert(index int, element E) error {
	if len(l.elements) < index {
		return errors.New("index out of range")
	}
	l.elements = append(l.elements[:index], append([]E{element}, l.elements[index:]...)...)
	return nil
}

// Remove removes the first occurrence of an element from the list.
func (l *List[E]) Remove(element E) {
	for i, el := range l.elements {
		if el == element {
			l.elements = append(l.elements[:i], l.elements[i+1:]...)
			return
		}
	}
}

// RemoveAt removes the element at the given index and returns the removed element.
//
// An error is returned if the index is out of range.
func (l *List[E]) RemoveAt(index int) (E, error) {
	var defaultValue E
	if index < 0 || index >= len(l.elements) {
		return defaultValue, ErrIndexOutOfRange
	}
	el := l.elements[index]
	l.elements = append(l.elements[:index], l.elements[index+1:]...)
	return el, nil
}

// Any returns a bool indicating if any of the elements match the given function.
func (l *List[E]) Any(f func(e E) bool) bool {
	for _, e := range l.elements {
		if f(e) {
			return true
		}
	}
	return false
}

// All returns a bool indicating if all elements match the given function.
func (l *List[E]) All(f func(e E) bool) bool {
	if len(l.elements) == 0 {
		return false
	}
	for _, e := range l.elements {
		if !f(e) {
			return false
		}
	}
	return true
}

// FirstElement returns the element at index 0.
// If there are no elements in the List, a default value and an error are returned.
func (l *List[E]) FirstElement() (E, error) {
	var defaultValue E
	if len(l.elements) == 0 {
		return defaultValue, ErrEmptyList
	}
	return l.elements[0], nil
}

// First returns the first element in the List matching the function.
// A default value and error is returned if there are no elements in the List
// or if no matching element is present in the list.
func (l *List[E]) First(f func(E) bool) (E, error) {
	var defaultValue E
	if len(l.elements) == 0 {
		return defaultValue, ErrEmptyList
	}
	for _, el := range l.elements {
		if f(el) {
			return el, nil
		}
	}
	return defaultValue, ErrElementNotFound
}

// FirstOrDefault returns the first element in the list matching the given function.
// If no matching element is present, a default value is returned.
// If you need to know if a default value is returned, use First instead.
func (l *List[E]) FirstOrDefault(f func(E) bool) E {
	el, _ := l.First(f)
	return el
}

// LastElement returns the last element in the List.
// A default value and an error is returned if the List has no elements.
func (l *List[E]) LastElement() (E, error) {
	var defaultValue E
	if len(l.elements) == 0 {
		return defaultValue, ErrEmptyList
	}
	return l.elements[len(l.elements)-1], nil
}

// Single returns the first element matching the given function if it is the only element matching the function.
// If multiple elements match the function, a default value and an error are returned.
// If there are no elements in the List a default value and an error are returned.
func (l *List[E]) Single(f func(E) bool) (E, error) {
	var found E
	var alreadyFound bool
	if len(l.elements) == 0 {
		return found, ErrEmptyList
	}
	for _, el := range l.elements {
		if f(el) {
			found = el
			if alreadyFound {
				var defaultValue E
				return defaultValue, ErrMoreThanOneElement
			}
			alreadyFound = true
		}
	}
	return found, nil
}

// SingleOrDefault returns the first element in the List matching the given function
// if it is the only element matching the function.
//
// A default value is returned if there are multiple matching elements in the List.
//
// A default value is returned if there are no elements in the List.
func (l *List[E]) SingleOrDefault(f func(E) bool) E {
	var defaultValue E
	el, err := l.Single(f)
	if err != nil {
		return defaultValue
	}
	return el
}

// Contains returns a bool indicating if the List contains the given element.
func (l *List[E]) Contains(element E) bool {
	for _, e := range l.elements {
		if e == element {
			return true
		}
	}
	return false
}

// IndexOf returns the index of the first instance of the given element.
// The index of -1 is returned if the element is not found.
func (l *List[E]) IndexOf(element E) int {
	for i, el := range l.elements {
		if el == element {
			return i
		}
	}
	return -1
}

// Len returns the number of elements in the List.
func (l *List[E]) Len() int {
	return len(l.elements)
}

// Empty returns a bool indicating if the List is empty.
func (l *List[E]) Empty() bool {
	return l.Len() == 0
}

// ToSlice returns the raw slice data backing the List.
func (l *List[E]) ToSlice() []E {
	return l.elements
}

// ToMap returns a map of the List, where the key can be any comparable type and the value can
// be any value.
//
// The below example creates a map[int]string of people by name. If there were duplicate ages
// earlier values would be overwritten:
//
//	people := []person{
//		{
//			Name: "Mike",
//			Age:  34,
//		},
//		{
//			Name: "Liam",
//			Age:  24,
//		},
//	}
//	list := fn.NewList(people...)
//	namesByAge := list.ToMap(func(p person) (int, string) {
//		return p.Age, p.Name
//	})
func ToMap[E comparable, K comparable, V any, M map[K]V](l *List[E], f func(E) (K, V)) M {
	result := make(map[K]V)
	for _, el := range l.elements {
		key, value := f(el)
		result[key] = value
	}
	return result
}

// ToMap returns a map[string][E] of the List, where the key is any string and value is the element.
// The given function should return the string key and the element to be associated with the key
//
//	people := []person{
//		{
//			Name: "Mike",
//			Age:  34,
//		},
//		{
//			Name: "Jack",
//			Age:  21,
//		},
//		{
//			Name: "Liam",
//			Age:  24,
//		},
//	}
//	list := fn.NewList(people...)
//	peopleByName := list.ToMap(func(p person) (string, person) { return p.Name, p })
func (l *List[E]) ToMap(f func(E) (string, E)) map[string]E {
	return ToMap(l, f)
}

// Sort orders the elements in the list using the given sorting function.
//
//	 people := []person{
//		    {
//			    Name: "Mike",
//			    Age:  34,
//		    },
//		    {
//			    Name: "Jack",
//			    Age:  21,
//		    },
//		    {
//			    Name: "Liam",
//			    Age:  24,
//		    },
//	 }
//	 list := fn.NewList(people...)
//	 list.Sort(func(cur, next person) bool {
//		    return cur.Age < next.Age
//	 })
func (l *List[E]) Sort(f func(cur, next E) bool) *List[E] {
	sort.Slice(l.elements, func(i, j int) bool {
		return f(l.elements[i], l.elements[j])
	})
	return l
}

// Deduplicate removes duplicate elements from the List.
func (l *List[E]) Deduplicate() *List[E] {
	seen := make(map[E]interface{})
	result := make([]E, 0, len(l.elements))

	for _, e := range l.elements {
		if _, ok := seen[e]; !ok {
			seen[e] = struct{}{}
			result = append(result, e)
		}
	}
	return &List[E]{elements: result}
}

// Map applies the given function to each element and returns a list of the results.
func (l *List[E]) Map(f func(E) E) *List[E] {
	result := Map(l.elements, f)
	return &List[E]{elements: result}
}

// Filter returns a list of all elements that return true from the keepFunc.
func (l *List[E]) Filter(f func(E) bool) *List[E] {
	result := Filter(l.elements, f)
	return &List[E]{elements: result}
}

// Reduce applies the given function cumulatively to the elements of the slice, from left to right, to reduce the slice to a single value.
func (l *List[E]) Reduce(initialValue E, f func(cur, next E) E) E {
	return Reduce(l.elements, initialValue, f)
}

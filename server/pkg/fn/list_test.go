package fn_test

import (
	"advancely/pkg/fn"
	"github.com/stretchr/testify/assert"
	"testing"
)

type person struct {
	Name string
	Age  int
}

func TestList_NewList(t *testing.T) {
	testCases := []struct {
		name     string
		elements []int
	}{
		{
			name:     "no elements results in empty list",
			elements: []int{},
		},
		{
			name:     "elements results in list with elements",
			elements: []int{1, 2, 3, 4},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			list := fn.NewList(tc.elements...)
			assert.NotNil(t, list)
			assert.Equal(t, tc.elements, list.ToSlice())
		})
	}
}

func TestList_NewEmptyList(t *testing.T) {
	list := fn.NewEmptyList[int]()
	assert.NotNil(t, list)
	assert.Equal(t, 0, list.Len())
}

func TestList_Add(t *testing.T) {
	list := fn.NewList(1, 2, 3)
	list.Add(4)
	list.Add(5)

	assert.Len(t, list.ToSlice(), 5)
}

func TestList_Insert(t *testing.T) {
	list := fn.NewList(1, 2, 3)

	err := list.Insert(0, 0)
	assert.NoError(t, err)
	assert.Equal(t, []int{0, 1, 2, 3}, list.ToSlice())

	err = list.Insert(1, 10)
	assert.NoError(t, err)
	assert.Equal(t, []int{0, 10, 1, 2, 3}, list.ToSlice())

	err = list.Insert(5, 20)
	assert.NoError(t, err)
	assert.Equal(t, []int{0, 10, 1, 2, 3, 20}, list.ToSlice())
}

func TestList_Remove(t *testing.T) {
	list := fn.NewList(1, 2, 3)
	list.Remove(2)

	assert.Len(t, list.ToSlice(), 2)
	for _, v := range list.ToSlice() {
		assert.NotEqual(t, v, 2)
	}
}

func TestList_RemoveAt(t *testing.T) {
	list := fn.NewList(1, 2, 3)
	el, err := list.RemoveAt(1)
	assert.NoError(t, err)
	assert.Equal(t, el, 2)
}

func TestList_Any(t *testing.T) {
	testCases := []struct {
		name     string
		elements []int
		expected bool
		f        func(int) bool
	}{
		{
			name:     "single true",
			elements: []int{1, 2, 3},
			expected: true,
			f: func(i int) bool {
				return i > 2
			},
		},
		{
			name:     "single true, first element",
			elements: []int{1, 2, 3},
			expected: true,
			f: func(i int) bool {
				return i == 1
			},
		},
		{
			name:     "single true, middle element",
			elements: []int{1, 2, 3},
			expected: true,
			f: func(i int) bool {
				return i == 2
			},
		},
		{
			name:     "single true, last element",
			elements: []int{1, 2, 3},
			expected: true,
			f: func(i int) bool {
				return i == 3
			},
		},
		{
			name:     "none false",
			elements: []int{1, 2, 3, 4},
			expected: false,
			f: func(i int) bool {
				return i == 100
			},
		},
		{
			name:     "empty should return false",
			elements: []int{},
			expected: false,
			f:        func(i int) bool { return i == 1 },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			list := fn.NewList(tc.elements...)
			assert.Equal(t, tc.expected, list.Any(tc.f))
		})
	}
}

func TestList_All(t *testing.T) {
	testCases := []struct {
		name     string
		elements []int
		expected bool
		f        func(int) bool
	}{
		{
			name:     "single match should return false",
			elements: []int{1, 2, 3},
			expected: false,
			f: func(i int) bool {
				return i > 2
			},
		},
		{
			name:     "no match should return false",
			elements: []int{1, 2, 3},
			expected: false,
			f: func(i int) bool {
				return i == 10
			},
		},
		{
			name:     "most match should return false",
			elements: []int{1, 2, 3},
			expected: false,
			f: func(i int) bool {
				return i <= 2
			},
		},
		{
			name:     "all match should return true",
			elements: []int{1, 2, 3},
			expected: true,
			f: func(i int) bool {
				return i < 10
			},
		},
		{
			name:     "empty should return false",
			elements: []int{},
			expected: false,
			f:        func(i int) bool { return i == 1 },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			list := fn.NewList(tc.elements...)
			assert.Equal(t, tc.expected, list.All(tc.f))
		})
	}
}

func TestList_FirstElement(t *testing.T) {
	list := fn.NewList(1, 2, 3)
	i, err := list.FirstElement()

	assert.NoError(t, err)
	assert.Equal(t, i, 1)
}

func TestList_FirstElement_ReturnsErrorWhenEmpty(t *testing.T) {
	var elements []int
	list := fn.NewList(elements...)
	i, err := list.FirstElement()

	assert.Error(t, err)
	assert.Equal(t, i, 0)
}

func TestList_First(t *testing.T) {
	testCases := []struct {
		name          string
		elements      []person
		f             func(person) bool
		expectedIndex int
		expectedError error
	}{
		{
			name:          "should return error when list is empty",
			elements:      []person{},
			f:             func(p person) bool { return true },
			expectedError: fn.ErrEmptyList,
		},
		{
			name: "should return error when no matching element in List",
			elements: []person{
				{
					Name: "Mike",
					Age:  34,
				},
			},
			f:             func(p person) bool { return false },
			expectedError: fn.ErrElementNotFound,
		},
		{
			name: "should return first element matching function",
			elements: []person{
				{
					Name: "Mike",
					Age:  34,
				},
				{
					Name: "Liam",
					Age:  24,
				},
			},
			f: func(p person) bool {
				return p.Name == "Liam"
			},
			expectedIndex: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			list := fn.NewList(tc.elements...)
			p, err := list.First(func(p person) bool { return tc.f(p) })
			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				assert.Equal(t, "", p.Name)
				assert.Equal(t, 0, p.Age)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.elements[tc.expectedIndex], p)
			}
		})
	}
}

func TestList_FirstOrDefault(t *testing.T) {
	testCases := []struct {
		name     string
		elements []person
		expected person
		f        func(person) bool
	}{
		{
			name:     "empty List should return default",
			elements: []person{},
			expected: person{},
			f:        func(person) bool { return true },
		},
		{
			name: "no matches should return default",
			elements: []person{
				{
					Name: "Mike",
					Age:  34,
				},
			},
			expected: person{},
			f:        func(person) bool { return false },
		},
		{
			name: "match at index 0 should return match",
			elements: []person{
				{
					Name: "Mike",
					Age:  34,
				},
				{
					Name: "Jack",
					Age:  34,
				},
			},
			expected: person{
				Name: "Mike",
				Age:  34,
			},
			f: func(p person) bool { return p.Age == 34 },
		},
		{
			name: "match at last index should return match",
			elements: []person{
				{
					Name: "Mike",
					Age:  34,
				},
				{
					Name: "Jack",
					Age:  21,
				},
			},
			expected: person{
				Name: "Jack",
				Age:  21,
			},
			f: func(p person) bool { return p.Age < 30 },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			list := fn.NewList(tc.elements...)
			assert.Equal(t, tc.expected, list.FirstOrDefault(tc.f))
		})
	}
}

func TestList_LastElement(t *testing.T) {
	testCases := []struct {
		name          string
		elements      []int
		expected      int
		expectedError error
	}{
		{
			name:          "empty List should return default",
			elements:      []int{},
			expectedError: fn.ErrEmptyList,
		},
		{
			name:     "should return last element when only one element in list",
			elements: []int{1},
			expected: 1,
		},
		{
			name:     "should return last element when multiple element in list",
			elements: []int{1, 2, 3, 4, 5},
			expected: 5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			list := fn.NewList(tc.elements...)
			v, err := list.LastElement()

			assert.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expected, v)
		})
	}
}

func TestList_Single(t *testing.T) {
	testCases := []struct {
		name        string
		elements    []person
		expected    person
		expectedErr error
		f           func(person) bool
	}{
		{
			name:        "empty list should return error",
			elements:    []person{},
			expected:    person{},
			expectedErr: fn.ErrEmptyList,
			f:           func(p person) bool { return true },
		},
		{
			name: "more than one match should return error",
			elements: []person{
				{
					Name: "Mike",
					Age:  34,
				},
				{
					Name: "Jack",
					Age:  21,
				},
			},
			expected:    person{},
			expectedErr: fn.ErrMoreThanOneElement,
			f:           func(p person) bool { return p.Age > 18 },
		},
		{
			name: "single match should return match",
			elements: []person{
				{
					Name: "Mike",
					Age:  34,
				},
				{
					Name: "Jack",
					Age:  21,
				},
			},
			expected: person{
				Name: "Jack",
				Age:  21,
			},
			f: func(p person) bool { return p.Age < 30 },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			list := fn.NewList(tc.elements...)
			p, err := list.Single(tc.f)

			assert.ErrorIs(t, err, tc.expectedErr)
			assert.Equal(t, tc.expected, p)
		})
	}
}

func TestList_SingleOrDefault(t *testing.T) {
	testCases := []struct {
		name     string
		elements []person
		expected person
		f        func(person) bool
	}{
		{
			name:     "empty list should return default",
			elements: []person{},
			expected: person{},
			f:        func(p person) bool { return true },
		},
		{
			name: "more than one match should return default",
			elements: []person{
				{
					Name: "Mike",
					Age:  34,
				},
				{
					Name: "Jack",
					Age:  21,
				},
			},
			expected: person{},
			f:        func(p person) bool { return p.Age > 18 },
		},
		{
			name: "single match should return match",
			elements: []person{
				{
					Name: "Mike",
					Age:  34,
				},
				{
					Name: "Jack",
					Age:  21,
				},
			},
			expected: person{
				Name: "Jack",
				Age:  21,
			},
			f: func(p person) bool { return p.Age < 30 },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			list := fn.NewList(tc.elements...)
			p := list.SingleOrDefault(tc.f)

			assert.Equal(t, tc.expected, p)
		})
	}
}

func TestList_Contains(t *testing.T) {
	testCases := []struct {
		name     string
		elements []person
		expected bool
	}{
		{
			name:     "should return false when no elements in List",
			elements: []person{},
			expected: false,
		},
		{
			name: "should return false when no matching element in List",
			elements: []person{
				{
					Name: "Jack",
					Age:  21,
				},
			},
			expected: false,
		},
		{
			name: "should return true when matching element in List",
			elements: []person{
				{
					Name: "Jack",
					Age:  21,
				},
				{
					Name: "Mike",
					Age:  34,
				},
			},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			list := fn.NewList(tc.elements...)
			assert.Equal(t, tc.expected, list.Contains(person{Name: "Mike", Age: 34}))
		})
	}
}

func TestList_IndexOf(t *testing.T) {
	testCases := []struct {
		name     string
		elements []person
		expected int
	}{
		{
			name:     "should return -1 when no elements in list",
			elements: []person{},
			expected: -1,
		},
		{
			name: "should return -1 when elements not in list",
			elements: []person{
				{
					Name: "Jack",
					Age:  21,
				},
			},
			expected: -1,
		},
		{
			name: "should return index when elements in list",
			elements: []person{
				{
					Name: "Jack",
					Age:  21,
				},
				{
					Name: "Mike",
					Age:  34,
				},
			},
			expected: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			list := fn.NewList(tc.elements...)
			assert.Equal(t, tc.expected, list.IndexOf(person{Name: "Mike", Age: 34}))
		})
	}
}

func TestList_Len(t *testing.T) {
	list := fn.NewList(1, 2, 3, 4, 5)
	assert.Equal(t, 5, list.Len())
}

func TestList_Empty(t *testing.T) {
	testCases := []struct {
		name     string
		elements []int
		expected bool
	}{
		{
			name:     "should return true when no elements in List",
			elements: []int{},
			expected: true,
		},
		{
			name:     "should return false when List has elements",
			elements: []int{1, 2, 3, 4, 5},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			list := fn.NewList(tc.elements...)
			assert.Equal(t, tc.expected, list.Empty())
		})
	}
}

func TestList_ToSlice(t *testing.T) {
	expected := []int{1, 2, 3, 4, 5}
	list := fn.NewList(expected...)

	actual := list.ToSlice()
	assert.Equal(t, expected, actual)
}

func TestList_ToMap(t *testing.T) {
	testCases := []struct {
		name     string
		elements []person
		expected map[string]person
	}{
		{
			name:     "empty list returns empty map",
			elements: []person{},
			expected: map[string]person{},
		},
		{
			name: "single element List should result in single element map",
			elements: []person{
				{
					Name: "Mike",
					Age:  34,
				},
			},
			expected: map[string]person{
				"Mike": {
					Name: "Mike",
					Age:  34,
				},
			},
		},
		{
			name: "multi element List should result in multi element map",
			elements: []person{
				{
					Name: "Jack",
					Age:  21,
				},
				{
					Name: "Mike",
					Age:  34,
				},
			},
			expected: map[string]person{
				"Jack": {
					Name: "Jack",
					Age:  21,
				},
				"Mike": {
					Name: "Mike",
					Age:  34,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			list := fn.NewList(tc.elements...)
			m := list.ToMap(func(p person) (string, person) {
				return p.Name, p
			})

			assert.Equal(t, tc.expected, m)
		})
	}
}

func Test_ToMap(t *testing.T) {
	testCases := []struct {
		name     string
		elements []person
		expected map[int]string
	}{
		{
			name:     "empty list returns empty map",
			elements: []person{},
			expected: map[int]string{},
		},
		{
			name: "single element List should result in single element map",
			elements: []person{
				{
					Name: "Mike",
					Age:  34,
				},
			},
			expected: map[int]string{
				34: "Mike",
			},
		},
		{
			name: "multi element List should result in multi element map",
			elements: []person{
				{
					Name: "Jack",
					Age:  21,
				},
				{
					Name: "Mike",
					Age:  34,
				},
			},
			expected: map[int]string{
				21: "Jack",
				34: "Mike",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			list := fn.NewList(tc.elements...)
			m := fn.ToMap(list, func(p person) (int, string) {
				return p.Age, p.Name
			})

			assert.Equal(t, tc.expected, m)
		})
	}
}

func TestList_Sort(t *testing.T) {
	people := []person{
		{
			Name: "Mike",
			Age:  34,
		},
		{
			Name: "Jack",
			Age:  21,
		},
		{
			Name: "Liam",
			Age:  24,
		},
	}

	expectedPeople := []person{
		{
			Name: "Jack",
			Age:  21,
		},
		{
			Name: "Liam",
			Age:  24,
		},
		{
			Name: "Mike",
			Age:  34,
		},
	}

	list := fn.NewList(people...)
	list.Sort(func(cur, next person) bool {
		return cur.Age < next.Age
	})

	expected := fn.NewList(expectedPeople...)
	assert.Equal(t, expected, list)
}

func TestList_Deduplicate(t *testing.T) {
	testCases := []struct {
		name     string
		elements []person
		expected []person
	}{
		{
			name:     "empty list should return empty list",
			elements: []person{},
			expected: []person{},
		},
		{
			name: "list without duplicates should return equal list",
			elements: []person{
				{
					Name: "Mike",
					Age:  34,
				},
				{
					Name: "Liam",
					Age:  24,
				},
			},
			expected: []person{
				{
					Name: "Mike",
					Age:  34,
				},
				{
					Name: "Liam",
					Age:  24,
				},
			},
		},
		{
			name: "list with duplicates should return deduplicated list",
			elements: []person{
				{
					Name: "Mike",
					Age:  34,
				},
				{
					Name: "Liam",
					Age:  24,
				},
				{
					Name: "Mike",
					Age:  34,
				},
			},
			expected: []person{
				{
					Name: "Mike",
					Age:  34,
				},
				{
					Name: "Liam",
					Age:  24,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			list := fn.NewList(tc.elements...)
			expected := fn.NewList(tc.expected...)
			assert.Equal(t, expected, list.Deduplicate())
		})
	}
}

func TestList_Map(t *testing.T) {
	testCases := []struct {
		name     string
		elements []person
		expected []person
		f        func(person) person
	}{
		{
			name:     "empty list returns empty List",
			elements: []person{},
			expected: []person{},
			f:        func(p person) person { return p },
		},
		{
			name: "returns mapped results",
			elements: []person{
				{
					Name: "Mike",
					Age:  34,
				},
				{
					Name: "Jack",
					Age:  21,
				},
			},
			expected: []person{
				{
					Name: "Mike!",
					Age:  35,
				},
				{
					Name: "Jack!",
					Age:  22,
				},
			},
			f: func(p person) person {
				return person{
					Name: p.Name + "!",
					Age:  p.Age + 1,
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			list := fn.NewList(tc.elements...)
			expected := fn.NewList(tc.expected...)
			assert.Equal(t, expected, list.Map(tc.f))
		})
	}
}

func TestList_Filter(t *testing.T) {
	testCases := []struct {
		name     string
		elements []person
		expected []person
		f        func(person) bool
	}{
		{
			name:     "empty list returns empty List",
			elements: []person{},
			expected: []person{},
			f:        func(p person) bool { return true },
		},
		{
			name: "no matches returns empty List",
			elements: []person{
				{
					Name: "Mike",
					Age:  34,
				},
			},
			expected: []person{},
			f:        func(p person) bool { return false },
		},
		{
			name: "return matching results",
			elements: []person{
				{
					Name: "Mike",
					Age:  34,
				},
				{
					Name: "Jack",
					Age:  21,
				},
			},
			expected: []person{
				{
					Name: "Jack",
					Age:  21,
				},
			},
			f: func(p person) bool { return p.Age < 30 },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			list := fn.NewList(tc.elements...)
			expected := fn.NewList(tc.expected...)
			assert.Equal(t, expected, list.Filter(tc.f))
		})
	}
}

func TestList_Reduce(t *testing.T) {
	testCases := []struct {
		name     string
		elements []person
		expected person
	}{
		{
			name:     "empty list returns initial value",
			elements: []person{},
			expected: person{},
		},
		{
			name: "single element List returns correct value",
			elements: []person{
				{
					Name: "Mike",
					Age:  34,
				},
			},
			expected: person{
				Name: "Mike",
				Age:  34,
			},
		},
		{
			name: "multi-element List returns correct value",
			elements: []person{
				{
					Name: "Mike",
					Age:  34,
				},
				{
					Name: "Jack",
					Age:  21,
				},
				{
					Name: "Liam",
					Age:  24,
				},
			},
			expected: person{
				Name: "Mike, Jack, Liam",
				Age:  79,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			list := fn.NewList(tc.elements...)
			p := list.Reduce(person{}, func(cur, next person) person {
				namePrefix := cur.Name
				if len(cur.Name) > 0 {
					namePrefix += ", "
				}
				return person{
					Name: namePrefix + next.Name,
					Age:  cur.Age + next.Age,
				}
			})

			assert.Equal(t, tc.expected, p)
		})
	}
}

func TestList_CanChainFunctions(t *testing.T) {
	items := fn.NewList(5, 1, 4, 2, 4, 4, 5, 3, 3, 3).
		Deduplicate().
		Filter(func(x int) bool { return x <= 3 }).
		Sort(func(cur, next int) bool { return cur < next }).
		ToSlice()

	assert.Equal(t, []int{1, 2, 3}, items)
}

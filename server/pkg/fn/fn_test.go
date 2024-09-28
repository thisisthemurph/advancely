package fn_test

import (
	"reflect"
	"testing"

	"advancely/pkg/fn"
)

func Test_Map_WithInt(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		want     []int
		function func(int) int
	}{
		{
			name:  "double",
			input: []int{1, 2, 3, 4, 5},
			want:  []int{2, 4, 6, 8, 10},
			function: func(x int) int {
				return x * 2
			},
		},
		{
			name:  "add",
			input: []int{1, 2, 3, 4, 5},
			want:  []int{11, 12, 13, 14, 15},
			function: func(x int) int {
				return x + 10
			},
		},
		{
			name:  "make 0",
			input: []int{-1, 2, 3, 4, 5},
			want:  []int{0, 0, 0, 0, 0},
			function: func(x int) int {
				return x - x
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := fn.Map(tc.input, tc.function)
			if !reflect.DeepEqual(result, tc.want) {
				t.Errorf("got %v, want %v", result, tc.want)
			}
		})
	}
}

func Test_Filter_WithInt(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		want     []int
		function func(int) bool
	}{
		{
			name:  "more than 10",
			input: []int{2, 4, 6, 8, 10, 12, 14},
			want:  []int{12, 14},
			function: func(x int) bool {
				return x > 10
			},
		},
		{
			name:  "is evan",
			input: []int{1, 2, 3, 4, 5},
			want:  []int{2, 4},
			function: func(x int) bool {
				return x%2 == 0
			},
		},
		{
			name:  "not negative",
			input: []int{-1, -2, -3, 4, 5},
			want:  []int{4, 5},
			function: func(x int) bool {
				return x >= 0
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := fn.Filter(tc.input, tc.function)
			if !reflect.DeepEqual(result, tc.want) {
				t.Errorf("got %v, want %v", result, tc.want)
			}
		})
	}
}

func Test_Reduce_WithInt(t *testing.T) {
	testCases := []struct {
		name         string
		input        []int
		want         int
		initialValue int
		function     func(int, int) int
	}{
		{
			name:     "add all",
			input:    []int{1, 2, 3, 4, 5},
			want:     15,
			function: func(cur, next int) int { return cur + next },
		},
		{
			name:         "multiply all",
			input:        []int{1, 2, 3, 4, 5},
			want:         120,
			initialValue: 1,
			function:     func(cur, next int) int { return cur * next },
		},
		{
			name:         "largest value",
			input:        []int{2, 1, 5, 4, 3},
			want:         5,
			initialValue: 0,
			function: func(cur, next int) int {
				if next < cur {
					return cur
				}
				return next
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := fn.Reduce(tc.input, tc.initialValue, tc.function)
			if result != tc.want {
				t.Errorf("got %v, want %v", result, tc.want)
			}
		})
	}
}

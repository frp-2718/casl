package main

import (
	"testing"

	"golang.org/x/exp/slices"
)

func TestFilter(t *testing.T) {
	elementsToFilter := []string{"a", "b", "c", "d", "e"}
	var tests = []struct {
		input []string
		want  []string
	}{
		{
			nil,
			nil,
		},
		{
			[]string{},
			[]string{},
		},
		{
			[]string{"d", "a"},
			[]string{},
		},
		{
			[]string{"a", "f", "b", "g"},
			[]string{"f", "g"},
		},
		{
			[]string{"f", "g"},
			[]string{"f", "g"},
		},
	}

	for _, test := range tests {
		got := filter(test.input, elementsToFilter)
		if equal := slices.Compare(got, test.want); equal != 0 {
			t.Errorf("filter(%v, %v) = %v, want %v", test.input, elementsToFilter, got, test.want)
		}
	}
}

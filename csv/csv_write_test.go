package csv

import (
	"reflect"
	"testing"
)

func Test_padRows(t *testing.T) {
	input := [][]string{
		{"col1", "col2", "col3", "col4", "col5"},
		{"a1", "a2", "a3", "a4", "a5"},
		{"b1", "b2"},
		{"c1", "c2", "c3", "c4", "c5", "c6"},
	}
	want := [][]string{
		{"col1", "col2", "col3", "col4", "col5"},
		{"a1", "a2", "a3", "a4", "a5"},
		{"b1", "b2", "", "", ""},
		{"c1", "c2", "c3", "c4", "c5", "c6"},
	}

	if got := padRows(input); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

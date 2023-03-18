package collection

import (
	"fmt"
	"reflect"
	"testing"
)

func TestPadSlice(t *testing.T) {
	tests := []struct {
		name   string
		slice  []int
		length int
		want   []int
	}{
		{
			name:  "slice is empty and length is 0",
			slice: []int{}, length: 0, want: []int{},
		},
		{
			name:  "slice is longer than pad length",
			slice: []int{1}, length: 0, want: []int{1},
		},
		{
			name:  "slice length is same as pad length",
			slice: []int{1}, length: 1, want: []int{1},
		},
		{
			name:  "slice is shorter than pad length",
			slice: []int{1}, length: 2, want: []int{1, 0},
		},
		{
			name:  "slice is much shorter than pad length",
			slice: []int{1}, length: 10, want: []int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PadSliceRight(tt.slice, tt.length)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewOrderedMap(t *testing.T) {
	m := NewOrderedMap[string, int]()
	if m.m == nil {
		t.Error("want map to be initialized, got nil")
	}
	if len(m.order) != 0 {
		t.Error("want order to be empty, got non-empty slice")
	}
}

func TestOrderedMap_Set(t *testing.T) {
	m := NewOrderedMap[string, int]()
	m.Set("one", 1)
	m.Set("two", 2)
	m.Set("three", 3)

	// Verify map:
	if _, exist := m.m["one"]; !exist {
		t.Error("want one to exist")
	}
	if _, exist := m.m["two"]; !exist {
		t.Error("want two to exist")
	}
	if _, exist := m.m["three"]; !exist {
		t.Error("want three to exist")
	}

	// Verify order:
	{
		got := []string{"one", "two", "three"}
		want := m.order
		if fmt.Sprint(got) != fmt.Sprint(want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestOrderedMap_Get(t *testing.T) {
	m := NewOrderedMap[string, int]()
	m.Set("one", 1)
	m.Set("two", 2)

	v1, exists1 := m.Get("one")
	if !exists1 || v1 != 1 {
		t.Error("want value 1 for key one, got", v1)
	}
	v2, exists2 := m.Get("two")
	if !exists2 || v2 != 2 {
		t.Error("want value 2 for key two, got", v2)
	}
	_, exists3 := m.Get("three")
	if exists3 {
		t.Error("want key three to not exist, got true")
	}
}

func TestOrderedMap_Slice(t *testing.T) {
	m := NewOrderedMap[string, int]()
	m.Set("one", 1)
	m.Set("two", 2)
	m.Set("three", 3)

	got := m.Slice()
	want := []int{1, 2, 3}

	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

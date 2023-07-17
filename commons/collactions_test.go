package commons

import "testing"

func TestIsFullFalse(t *testing.T) {
	slice := make([]int, 3, 5)
	slice[0], slice[1], slice[2] = 1, 2, 3
	got := IsFull[int](slice)
	want := false
	if got != want {
		t.Errorf("TestIsFull() failed, should be not full %t %t %d %d", got, want, len(slice), cap(slice))
	}
}

func TestIsFullTrue(t *testing.T) {
	slice := make([]int, 3, 5)
	slice[0], slice[1], slice[2] = 1, 2, 3
	slice = append(slice, 4, 5)
	got := IsFull[int](slice)
	want := true
	if got != want {
		t.Errorf("TestIsFull() failed, should be full %t %t %d %d", got, want, len(slice), cap(slice))
	}
}
